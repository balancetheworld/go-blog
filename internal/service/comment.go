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

	if req.PostID > 0 {
		post, err := repo.GetPostByID(ctx, req.PostID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.CommentListResponse{}, errs.NewNotFound(
				http.StatusNotFound,
				"post not found",
			)
		}
		if err != nil {
			return dto.CommentListResponse{}, errs.NewInternalServer(
				http.StatusInternalServerError,
				"get post failed",
			)
		}
		if !canViewPost(post, viewerID, role, viewerRoleID) {
			return dto.CommentListResponse{}, errs.NewNotFound(
				http.StatusNotFound,
				"post not found",
			)
		}
	} else if viewerID == 0 || role != constant.RoleAdmin {
		return dto.CommentListResponse{}, errs.NewForbidden(
			http.StatusForbidden,
			"comment list access denied",
		)
	}

	comments, total, err := repo.ListComments(
		ctx,
		req.PostID,
		(req.Page-1)*req.PageSize,
		req.PageSize,
	)
	if err != nil {
		return dto.CommentListResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"list comments failed",
		)
	}

	items := make([]dto.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		items = append(items, comment.ToDto())
	}

	return dto.CommentListResponse{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
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

	post, err := repo.GetPostByID(ctx, req.PostID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.CommentResponse{}, errs.NewNotFound(
			http.StatusNotFound,
			"post not found",
		)
	}
	if err != nil {
		return dto.CommentResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get post failed",
		)
	}
	if post.Content == "" {
		return dto.CommentResponse{}, errs.NewBadRequest(
			http.StatusBadRequest,
			"cannot comment on draft post",
		)
	}
	if !canViewPost(post, authorID, role, viewerRoleID) {
		return dto.CommentResponse{}, errs.NewNotFound(
			http.StatusNotFound,
			"post not found",
		)
	}

	comment := model.Comment{
		PostID:   post.ID,
		AuthorID: authorID,
		Content:  content,
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
		return errs.NewNotFound(
			http.StatusNotFound,
			"comment not found",
		)
	}
	if err != nil {
		return errs.NewInternalServer(
			http.StatusInternalServerError,
			"get comment failed",
		)
	}

	canDelete := viewerID == comment.AuthorID ||
		role == constant.RoleAdmin ||
		canManagePost(comment.Post, viewerID, role)
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
		return errs.NewNotFound(
			http.StatusNotFound,
			"comment not found",
		)
	}

	return nil
}
