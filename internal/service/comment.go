package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/errs"
	"gorm.io/gorm"
)

const maxCommentDepth uint8 = 2

type commentTarget struct {
	targetType constant.TargetType
	targetID   uint
	post       *model.Post
	diary      *model.Diary
}

func ListComments(
	ctx context.Context,
	req dto.CommentListRequest,
	viewerID uint,
	role constant.Role,
	viewerRoleID uint,
) (dto.CommentListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PostID == 0 && req.TargetType == "" && req.TargetID == 0 {
		return ListAdminComments(ctx, dto.AdminCommentListRequest{
			Page:     req.Page,
			PageSize: req.PageSize,
		}, role)
	}

	targetType, targetID, err := normalizeCommentListTarget(req)
	if err != nil {
		return dto.CommentListResponse{}, err
	}
	if _, err := getCommentTarget(
		ctx,
		targetType,
		targetID,
		viewerID,
		role,
		viewerRoleID,
	); err != nil {
		return dto.CommentListResponse{}, err
	}

	comments, total, err := repo.ListComments(ctx, repo.CommentListFilter{
		TargetType:   targetType,
		TargetID:     targetID,
		TopLevelOnly: true,
		Offset:       (req.Page - 1) * req.PageSize,
		Limit:        req.PageSize,
	})
	if err != nil {
		return dto.CommentListResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"list comments failed",
		)
	}

	return commentListResponse(comments, total, req.Page, req.PageSize), nil
}

func ListCommentReplies(
	ctx context.Context,
	parentID uint,
	viewerID uint,
	role constant.Role,
	viewerRoleID uint,
) ([]dto.CommentResponse, error) {
	parent, err := repo.GetCommentByID(ctx, parentID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NewNotFound(http.StatusNotFound, "comment not found")
	}
	if err != nil {
		return nil, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get comment failed",
		)
	}
	if _, err := getCommentTarget(
		ctx,
		parent.TargetType,
		parent.TargetID,
		viewerID,
		role,
		viewerRoleID,
	); err != nil {
		return nil, err
	}

	replies, err := repo.ListCommentReplies(ctx, parentID)
	if err != nil {
		return nil, errs.NewInternalServer(
			http.StatusInternalServerError,
			"list comment replies failed",
		)
	}

	items := make([]dto.CommentResponse, 0, len(replies))
	for _, reply := range replies {
		items = append(items, reply.ToDto())
	}

	return items, nil
}

func ListAdminComments(
	ctx context.Context,
	req dto.AdminCommentListRequest,
	role constant.Role,
) (dto.CommentListResponse, error) {
	if role != constant.RoleAdmin {
		return dto.CommentListResponse{}, errs.NewForbidden(
			http.StatusForbidden,
			"admin comment access denied",
		)
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if (req.TargetType == "") != (req.TargetID == 0) {
		return dto.CommentListResponse{}, errs.NewBadRequest(
			http.StatusBadRequest,
			"comment target type and target id are required together",
		)
	}

	comments, total, err := repo.ListComments(ctx, repo.CommentListFilter{
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		AuthorID:   req.AuthorID,
		Keyword:    req.Keyword,
		Offset:     (req.Page - 1) * req.PageSize,
		Limit:      req.PageSize,
	})
	if err != nil {
		return dto.CommentListResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"list comments failed",
		)
	}

	return commentListResponse(comments, total, req.Page, req.PageSize), nil
}

func CreateComment(
	ctx context.Context,
	authorID uint,
	role constant.Role,
	viewerRoleID uint,
	req dto.CreateCommentRequest,
) (dto.CommentResponse, error) {
	if authorID == 0 || role == constant.RoleGuest {
		return dto.CommentResponse{}, errs.NewUnauthorized(
			http.StatusUnauthorized,
			"login required",
		)
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return dto.CommentResponse{}, errs.NewBadRequest(
			http.StatusBadRequest,
			"comment content is required",
		)
	}

	targetType, targetID, parent, err := resolveCreateCommentTarget(ctx, req)
	if err != nil {
		return dto.CommentResponse{}, err
	}
	if _, err := getCommentTarget(
		ctx,
		targetType,
		targetID,
		authorID,
		role,
		viewerRoleID,
	); err != nil {
		return dto.CommentResponse{}, err
	}

	comment := model.Comment{
		TargetType: targetType,
		TargetID:   targetID,
		AuthorID:   authorID,
		Content:    content,
	}
	if targetType == constant.TargetPost || targetType == constant.TargetPage {
		comment.PostID = &targetID
	}
	if parent != nil {
		if parent.Depth >= maxCommentDepth {
			return dto.CommentResponse{}, errs.NewBadRequest(
				http.StatusBadRequest,
				"maximum comment reply depth exceeded",
			)
		}
		comment.ParentID = &parent.ID
		comment.ReplyToUserID = &parent.AuthorID
		comment.Depth = parent.Depth + 1
		if parent.RootID != nil {
			comment.RootID = parent.RootID
		} else {
			comment.RootID = &parent.ID
		}
	}

	if err := repo.CreateComment(ctx, &comment); err != nil {
		return dto.CommentResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"create comment failed",
		)
	}

	createdComment, err := repo.GetCommentByID(ctx, comment.ID)
	if err != nil {
		return dto.CommentResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get created comment failed",
		)
	}

	return createdComment.ToDto(), nil
}

func DeleteComment(
	ctx context.Context,
	id uint,
	viewerID uint,
	role constant.Role,
) error {
	comment, err := repo.GetCommentByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.NewNotFound(http.StatusNotFound, "comment not found")
	}
	if err != nil {
		return errs.NewInternalServer(
			http.StatusInternalServerError,
			"get comment failed",
		)
	}

	canDelete := viewerID == comment.AuthorID || role == constant.RoleAdmin
	if !canDelete {
		canDelete, err = canManageCommentTarget(ctx, comment, viewerID, role)
		if err != nil {
			return err
		}
	}
	if !canDelete {
		return errs.NewForbidden(
			http.StatusForbidden,
			"delete comment access denied",
		)
	}

	rowsAffected, err := repo.DeleteComment(ctx, comment)
	if err != nil {
		return errs.NewInternalServer(
			http.StatusInternalServerError,
			"delete comment failed",
		)
	}
	if rowsAffected == 0 {
		return errs.NewNotFound(http.StatusNotFound, "comment not found")
	}

	return nil
}

func normalizeCommentListTarget(
	req dto.CommentListRequest,
) (constant.TargetType, uint, error) {
	if req.TargetID == 0 && req.PostID > 0 {
		return constant.TargetPost, req.PostID, nil
	}
	if req.TargetType == "" {
		req.TargetType = constant.TargetPost
	}
	if req.TargetID == 0 {
		return "", 0, errs.NewBadRequest(
			http.StatusBadRequest,
			"comment target is required",
		)
	}

	return req.TargetType, req.TargetID, nil
}

func resolveCreateCommentTarget(
	ctx context.Context,
	req dto.CreateCommentRequest,
) (constant.TargetType, uint, *model.Comment, error) {
	var parent *model.Comment
	parentID := req.ParentID
	if req.TargetType == constant.TargetComment {
		if req.TargetID == 0 {
			return "", 0, nil, errs.NewBadRequest(
				http.StatusBadRequest,
				"reply target comment is required",
			)
		}
		value := req.TargetID
		parentID = &value
	}

	if parentID != nil {
		value, err := repo.GetCommentByID(ctx, *parentID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", 0, nil, errs.NewNotFound(
				http.StatusNotFound,
				"parent comment not found",
			)
		}
		if err != nil {
			return "", 0, nil, errs.NewInternalServer(
				http.StatusInternalServerError,
				"get parent comment failed",
			)
		}
		parent = &value
	}

	if parent != nil {
		if req.TargetType != "" && req.TargetType != constant.TargetComment {
			if req.TargetType != parent.TargetType || req.TargetID != parent.TargetID {
				return "", 0, nil, errs.NewBadRequest(
					http.StatusBadRequest,
					"parent comment target mismatch",
				)
			}
		}

		return parent.TargetType, parent.TargetID, parent, nil
	}

	if req.TargetID == 0 && req.PostID > 0 {
		return constant.TargetPost, req.PostID, nil, nil
	}
	if req.TargetType == "" {
		req.TargetType = constant.TargetPost
	}
	if req.TargetType == constant.TargetComment || req.TargetID == 0 {
		return "", 0, nil, errs.NewBadRequest(
			http.StatusBadRequest,
			"comment target is required",
		)
	}

	return req.TargetType, req.TargetID, nil, nil
}

func getCommentTarget(
	ctx context.Context,
	targetType constant.TargetType,
	targetID uint,
	viewerID uint,
	role constant.Role,
	viewerRoleID uint,
) (commentTarget, error) {
	switch targetType {
	case constant.TargetPost, constant.TargetPage:
		post, err := repo.GetPostByID(ctx, targetID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return commentTarget{}, errs.NewNotFound(http.StatusNotFound, "comment target not found")
		}
		if err != nil {
			return commentTarget{}, errs.NewInternalServer(
				http.StatusInternalServerError,
				"get comment target failed",
			)
		}
		if (targetType == constant.TargetPage) != (post.Type == "page") {
			return commentTarget{}, errs.NewNotFound(http.StatusNotFound, "comment target not found")
		}
		if post.Content == "" || !canViewPost(post, viewerID, role, viewerRoleID) {
			return commentTarget{}, errs.NewNotFound(http.StatusNotFound, "comment target not found")
		}

		return commentTarget{targetType: targetType, targetID: targetID, post: &post}, nil
	case constant.TargetDiary:
		diary, err := repo.GetDiaryByID(ctx, targetID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return commentTarget{}, errs.NewNotFound(http.StatusNotFound, "comment target not found")
		}
		if err != nil {
			return commentTarget{}, errs.NewInternalServer(
				http.StatusInternalServerError,
				"get comment target failed",
			)
		}
		if diary.PublishedAt == nil || !canViewDiary(diary, viewerID, role, viewerRoleID) {
			return commentTarget{}, errs.NewNotFound(http.StatusNotFound, "comment target not found")
		}

		return commentTarget{targetType: targetType, targetID: targetID, diary: &diary}, nil
	default:
		return commentTarget{}, errs.NewBadRequest(
			http.StatusBadRequest,
			"unsupported comment target type",
		)
	}
}

func canManageCommentTarget(
	ctx context.Context,
	comment model.Comment,
	viewerID uint,
	role constant.Role,
) (bool, error) {
	switch comment.TargetType {
	case constant.TargetPost, constant.TargetPage:
		post, err := repo.GetPostByID(ctx, comment.TargetID)
		if err != nil {
			return false, errs.NewInternalServer(
				http.StatusInternalServerError,
				"get comment target failed",
			)
		}
		return canManagePost(post, viewerID, role), nil
	case constant.TargetDiary:
		diary, err := repo.GetDiaryByID(ctx, comment.TargetID)
		if err != nil {
			return false, errs.NewInternalServer(
				http.StatusInternalServerError,
				"get comment target failed",
			)
		}
		return canManageDiary(diary, viewerID, role), nil
	default:
		return false, nil
	}
}

func commentListResponse(
	comments []model.Comment,
	total int64,
	page int,
	pageSize int,
) dto.CommentListResponse {
	items := make([]dto.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		items = append(items, comment.ToDto())
	}

	return dto.CommentListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
