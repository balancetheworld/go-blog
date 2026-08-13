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

func canManageLabel(role constant.Role) bool {
	return role == constant.RoleEditor || role == constant.RoleAdmin
}

func CreateLabel(
	ctx context.Context,
	role constant.Role,
	req dto.CreateLabelRequest,
) (dto.LabelResponse, error) {
	if !canManageLabel(role) {
		return dto.LabelResponse{}, errs.NewForbidden(
			http.StatusForbidden,
			"create label access denied",
		)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return dto.LabelResponse{}, errs.NewBadRequest(
			http.StatusBadRequest,
			"label name is required",
		)
	}

	slugSource := strings.TrimSpace(req.Slug)
	if slugSource == "" {
		slugSource = name
	}
	slug := normalizePostSlug(slugSource)
	if slug == "" {
		return dto.LabelResponse{}, errs.NewBadRequest(
			http.StatusBadRequest,
			"label slug is required",
		)
	}

	exists, err := repo.CheckLabelExists(ctx, name, slug, 0)
	if err != nil {
		return dto.LabelResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"check label failed",
		)
	}
	if exists {
		return dto.LabelResponse{}, errs.NewConflict(
			http.StatusConflict,
			"label name or slug already exists",
		)
	}

	label := model.Label{
		Name: name,
		Slug: slug,
	}
	if err := repo.CreateLabel(ctx, &label); err != nil {
		return dto.LabelResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"create label failed",
		)
	}

	return toLabelResponse(label), nil
}

func UpdateLabel(
	ctx context.Context,
	id uint,
	role constant.Role,
	req dto.UpdateLabelRequest,
) (dto.LabelResponse, error) {
	if !canManageLabel(role) {
		return dto.LabelResponse{}, errs.NewForbidden(
			http.StatusForbidden,
			"update label access denied",
		)
	}

	label, err := repo.GetLabelByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.LabelResponse{}, errs.NewNotFound(
			http.StatusNotFound,
			"label not found",
		)
	}
	if err != nil {
		return dto.LabelResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get label failed",
		)
	}

	name := label.Name
	if req.Name != nil {
		name = strings.TrimSpace(*req.Name)
		if name == "" {
			return dto.LabelResponse{}, errs.NewBadRequest(
				http.StatusBadRequest,
				"label name is required",
			)
		}
	}

	slug := label.Slug
	if req.Slug != nil {
		slugSource := strings.TrimSpace(*req.Slug)
		if slugSource == "" {
			slugSource = name
		}
		slug = normalizePostSlug(slugSource)
		if slug == "" {
			return dto.LabelResponse{}, errs.NewBadRequest(
				http.StatusBadRequest,
				"label slug is required",
			)
		}
	}

	exists, err := repo.CheckLabelExists(ctx, name, slug, label.ID)
	if err != nil {
		return dto.LabelResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"check label failed",
		)
	}
	if exists {
		return dto.LabelResponse{}, errs.NewConflict(
			http.StatusConflict,
			"label name or slug already exists",
		)
	}

	label.Name = name
	label.Slug = slug
	if err := repo.UpdateLabel(ctx, &label); err != nil {
		return dto.LabelResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"update label failed",
		)
	}

	return toLabelResponse(label), nil
}

func DeleteLabel(ctx context.Context, id uint, role constant.Role) error {
	if !canManageLabel(role) {
		return errs.NewForbidden(
			http.StatusForbidden,
			"delete label access denied",
		)
	}

	rowsAffected, err := repo.DeleteLabel(ctx, id)
	if err != nil {
		return errs.NewInternalServer(
			http.StatusInternalServerError,
			"delete label failed",
		)
	}
	if rowsAffected == 0 {
		return errs.NewNotFound(
			http.StatusNotFound,
			"label not found",
		)
	}

	return nil
}
