package service

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/internal/task"
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

func toCommentResponse(comment model.Comment) dto.CommentResponse {
	var parentID *uint64
	if comment.ParentID != nil {
		value := uint64(*comment.ParentID)
		parentID = &value
	}

	var rootID *uint64
	if comment.RootID != nil {
		value := uint64(*comment.RootID)
		rootID = &value
	}

	var replyToUser *dto.UserDto
	if comment.ReplyToUser != nil {
		value := toUserResponse(*comment.ReplyToUser)
		replyToUser = &value
	}

	postID := uint64(0)
	if comment.PostID != nil {
		postID = uint64(*comment.PostID)
	}
	if comment.TargetType == constant.TargetPost || comment.TargetType == constant.TargetPage {
		postID = uint64(comment.TargetID)
	}

	return dto.CommentResponse{
		ID:                   uint64(comment.ID),
		PostID:               postID,
		TargetType:           comment.TargetType,
		TargetID:             uint64(comment.TargetID),
		ParentID:             parentID,
		RootID:               rootID,
		ReplyToUser:          replyToUser,
		Content:              comment.Content,
		ModerationStatus:     comment.ModerationStatus,
		ModerationReason:     comment.ModerationReason,
		ModerationCategories: comment.ModerationCategories,
		ModerationConfidence: comment.ModerationConfidence,
		ModeratedAt:          comment.ModeratedAt,
		Author:               toUserResponse(comment.Author),
		Depth:                comment.Depth,
		ReplyCount:           comment.ReplyCount,
		LikeCount:            comment.LikeCount,
		CreatedAt:            comment.CreatedAt,
		UpdatedAt:            comment.UpdatedAt,
	}
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
		TargetType:       targetType,
		TargetID:         targetID,
		ModerationStatus: constant.ModerationApproved,
		TopLevelOnly:     true,
		Offset:           (req.Page - 1) * req.PageSize,
		Limit:            req.PageSize,
		NewestFirst:      targetType == constant.TargetGuestbook,
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
	if parent.ModerationStatus != constant.ModerationApproved {
		return []dto.CommentResponse{}, nil
	}

	replies, err := repo.ListCommentReplies(
		ctx,
		parentID,
		constant.ModerationApproved,
	)
	if err != nil {
		return nil, errs.NewInternalServer(
			http.StatusInternalServerError,
			"list comment replies failed",
		)
	}

	items := make([]dto.CommentResponse, 0, len(replies))
	for _, reply := range replies {
		items = append(items, toCommentResponse(reply))
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

	moderationTask, err := repo.GetCommentModerationTask(ctx, comment.ID)
	if err != nil {
		log.Printf(
			"get comment moderation task failed: comment_id=%d err=%v",
			comment.ID,
			err,
		)
	} else if err := task.EnqueueCommentModeration(ctx, moderationTask.ID); err != nil {
		log.Printf(
			"enqueue comment moderation task failed: task_id=%d err=%v",
			moderationTask.ID,
			err,
		)
	}

	createdComment, err := repo.GetCommentByID(ctx, comment.ID)
	if err != nil {
		return dto.CommentResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get created comment failed",
		)
	}

	return toCommentResponse(createdComment), nil
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

func ModerateComment(
	ctx context.Context,
	id uint,
	role constant.Role,
	req dto.UpdateCommentModerationRequest,
) (dto.CommentResponse, error) {
	if role != constant.RoleAdmin {
		return dto.CommentResponse{}, errs.NewForbidden(
			http.StatusForbidden,
			"admin comment moderation denied",
		)
	}

	comment, err := repo.GetCommentByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.CommentResponse{}, errs.NewNotFound(http.StatusNotFound, "comment not found")
	}
	if err != nil {
		return dto.CommentResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get comment failed",
		)
	}

	if err := repo.UpdateCommentModeration(
		ctx,
		id,
		req.Status,
		strings.TrimSpace(req.Reason),
	); err != nil {
		return dto.CommentResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"update comment moderation failed",
		)
	}

	updatedComment, err := repo.GetCommentByID(ctx, comment.ID)
	if err != nil {
		return dto.CommentResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get updated comment failed",
		)
	}

	return toCommentResponse(updatedComment), nil
}

func GetCommentModeration(
	ctx context.Context,
	id uint,
	viewerID uint,
	role constant.Role,
) (dto.CommentModerationResponse, error) {
	comment, err := repo.GetCommentByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.CommentModerationResponse{}, errs.NewNotFound(http.StatusNotFound, "comment not found")
	}
	if err != nil {
		return dto.CommentModerationResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get comment failed",
		)
	}

	if role != constant.RoleAdmin && (viewerID == 0 || comment.AuthorID != viewerID) {
		return dto.CommentModerationResponse{}, errs.NewForbidden(
			http.StatusForbidden,
			"comment moderation access denied",
		)
	}

	return dto.CommentModerationResponse{
		ID:               uint64(comment.ID),
		ModerationStatus: comment.ModerationStatus,
		ModerationReason: comment.ModerationReason,
		ModeratedAt:      comment.ModeratedAt,
	}, nil
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
	case constant.TargetGuestbook:
		if targetID != 1 {
			return commentTarget{}, errs.NewNotFound(http.StatusNotFound, "comment target not found")
		}

		return commentTarget{targetType: targetType, targetID: targetID}, nil
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
		items = append(items, toCommentResponse(comment))
	}

	return dto.CommentListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
