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

func momentToResponse(moment model.Moment) dto.MomentResponse {
	return dto.MomentResponse{
		ID:        uint64(moment.ID),
		Content:   moment.Content,
		Author:    toUserResponse(moment.Author),
		CreatedAt: moment.CreatedAt,
		UpdatedAt: moment.UpdatedAt,
	}
}

func CreateMoment(
	ctx context.Context,
	authorID uint,
	role constant.Role,
	req dto.CreateMomentRequest,
) (dto.MomentResponse, error) {
	if authorID == 0 || (role != constant.RoleEditor && role != constant.RoleAdmin) {
		return dto.MomentResponse{}, errs.NewForbidden(http.StatusForbidden, "create moment access denied")
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return dto.MomentResponse{}, errs.NewBadRequest(http.StatusBadRequest, "moment content is required")
	}

	moment := model.Moment{
		Content:  content,
		AuthorID: authorID,
	}
	if err := repo.CreateMoment(ctx, &moment); err != nil {
		return dto.MomentResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "create moment failed")
	}

	createdMoment, err := repo.GetMomentByID(ctx, moment.ID)
	if err != nil {
		return dto.MomentResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "get created moment failed")
	}

	return momentToResponse(createdMoment), nil
}

func ListMoments(ctx context.Context, req dto.ListMomentsRequest) (dto.ListMomentsResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	moments, total, err := repo.ListMoments(
		ctx,
		(req.Page-1)*req.PageSize,
		req.PageSize,
	)
	if err != nil {
		return dto.ListMomentsResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "list moments failed")
	}

	items := make([]dto.MomentResponse, 0, len(moments))
	for _, moment := range moments {
		items = append(items, momentToResponse(moment))
	}

	return dto.ListMomentsResponse{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func DeleteMoment(ctx context.Context, id uint, viewerID uint, role constant.Role) error {
	moment, err := repo.GetMomentByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.NewNotFound(http.StatusNotFound, "moment not found")
	}
	if err != nil {
		return errs.NewInternalServer(http.StatusInternalServerError, "get moment failed")
	}
	if role != constant.RoleAdmin && (role != constant.RoleEditor || viewerID == 0 || moment.AuthorID != viewerID) {
		return errs.NewForbidden(http.StatusForbidden, "delete moment access denied")
	}

	rowsAffected, err := repo.DeleteMoment(ctx, id)
	if err != nil {
		return errs.NewInternalServer(http.StatusInternalServerError, "delete moment failed")
	}
	if rowsAffected == 0 {
		return errs.NewNotFound(http.StatusNotFound, "moment not found")
	}

	return nil
}
