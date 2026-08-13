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

func ListDiaryFolders(
	ctx context.Context,
	role constant.Role,
	viewerRoleID uint,
	all bool,
) (dto.ListDiaryFoldersResponse, error) {
	if all && role != constant.RoleEditor && role != constant.RoleAdmin {
		return dto.ListDiaryFoldersResponse{}, errs.NewForbidden(http.StatusForbidden, "diary folder list access denied")
	}

	folders, err := repo.ListDiaryFolders(ctx)
	if err != nil {
		return dto.ListDiaryFoldersResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "list diary folders failed")
	}

	items := make([]dto.DiaryFolderResponse, 0, len(folders))
	for _, folder := range folders {
		if all || canViewDiaryFolder(&folder, role, viewerRoleID) {
			items = append(items, diaryFolderToResponse(folder))
		}
	}

	return dto.ListDiaryFoldersResponse{Items: items}, nil
}

func CreateDiaryFolder(
	ctx context.Context,
	role constant.Role,
	req dto.CreateDiaryFolderRequest,
) (dto.DiaryFolderResponse, error) {
	if role != constant.RoleEditor && role != constant.RoleAdmin {
		return dto.DiaryFolderResponse{}, errs.NewForbidden(http.StatusForbidden, "create diary folder access denied")
	}

	visibility := req.Visibility
	if visibility == "" {
		visibility = constant.PostVisibilityPublic
	}
	visibleRoles, err := resolvePostVisibleRoles(ctx, visibility, req.VisibleRoleIDs)
	if err != nil {
		return dto.DiaryFolderResponse{}, err
	}
	slug := normalizeDiarySlug(req.Slug)
	if slug == "" {
		slug = normalizeDiarySlug(req.Name)
	}
	if slug == "" {
		return dto.DiaryFolderResponse{}, errs.NewBadRequest(http.StatusBadRequest, "diary folder slug is required")
	}
	exists, err := repo.CheckDiaryFolderSlugExists(ctx, slug, 0)
	if err != nil {
		return dto.DiaryFolderResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "check diary folder slug failed")
	}
	if exists {
		return dto.DiaryFolderResponse{}, errs.NewConflict(http.StatusConflict, "diary folder slug already exists")
	}

	folder := model.DiaryFolder{
		Name:         strings.TrimSpace(req.Name),
		Slug:         slug,
		Description:  strings.TrimSpace(req.Description),
		Cover:        strings.TrimSpace(req.Cover),
		Sort:         req.Sort,
		Visibility:   visibility,
		VisibleRoles: visibleRoles,
	}
	if err := repo.CreateDiaryFolder(ctx, &folder); err != nil {
		return dto.DiaryFolderResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "create diary folder failed")
	}

	created, err := repo.GetDiaryFolderByID(ctx, folder.ID)
	if err != nil {
		return dto.DiaryFolderResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "get diary folder failed")
	}

	return diaryFolderToResponse(created), nil
}

func UpdateDiaryFolder(
	ctx context.Context,
	id uint,
	role constant.Role,
	req dto.UpdateDiaryFolderRequest,
) (dto.DiaryFolderResponse, error) {
	if role != constant.RoleEditor && role != constant.RoleAdmin {
		return dto.DiaryFolderResponse{}, errs.NewForbidden(http.StatusForbidden, "update diary folder access denied")
	}

	folder, err := repo.GetDiaryFolderByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.DiaryFolderResponse{}, errs.NewNotFound(http.StatusNotFound, "diary folder not found")
	}
	if err != nil {
		return dto.DiaryFolderResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "get diary folder failed")
	}

	if req.Name != nil {
		folder.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		folder.Description = strings.TrimSpace(*req.Description)
	}
	if req.Cover != nil {
		folder.Cover = strings.TrimSpace(*req.Cover)
	}
	if req.Sort != nil {
		folder.Sort = *req.Sort
	}
	if req.Slug != nil {
		slug := normalizeDiarySlug(*req.Slug)
		if slug == "" {
			return dto.DiaryFolderResponse{}, errs.NewBadRequest(http.StatusBadRequest, "diary folder slug is required")
		}
		exists, checkErr := repo.CheckDiaryFolderSlugExists(ctx, slug, id)
		if checkErr != nil {
			return dto.DiaryFolderResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "check diary folder slug failed")
		}
		if exists {
			return dto.DiaryFolderResponse{}, errs.NewConflict(http.StatusConflict, "diary folder slug already exists")
		}
		folder.Slug = slug
	}

	visibility := folder.Visibility
	if req.Visibility != nil {
		visibility = *req.Visibility
	}
	visibleRoleIDs := make([]uint, 0, len(folder.VisibleRoles))
	for _, visibleRole := range folder.VisibleRoles {
		visibleRoleIDs = append(visibleRoleIDs, visibleRole.ID)
	}
	if req.VisibleRoleIDs != nil {
		visibleRoleIDs = *req.VisibleRoleIDs
	}
	visibleRoles, err := resolvePostVisibleRoles(ctx, visibility, visibleRoleIDs)
	if err != nil {
		return dto.DiaryFolderResponse{}, err
	}
	folder.Visibility = visibility
	folder.VisibleRoles = visibleRoles

	if err := repo.UpdateDiaryFolder(ctx, &folder); err != nil {
		return dto.DiaryFolderResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "update diary folder failed")
	}

	updated, err := repo.GetDiaryFolderByID(ctx, folder.ID)
	if err != nil {
		return dto.DiaryFolderResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "get diary folder failed")
	}

	return diaryFolderToResponse(updated), nil
}

func DeleteDiaryFolder(ctx context.Context, id uint, role constant.Role) error {
	if role != constant.RoleEditor && role != constant.RoleAdmin {
		return errs.NewForbidden(http.StatusForbidden, "delete diary folder access denied")
	}

	rowsAffected, err := repo.DeleteDiaryFolder(ctx, id)
	if err != nil {
		return errs.NewInternalServer(http.StatusInternalServerError, "delete diary folder failed")
	}
	if rowsAffected == 0 {
		return errs.NewNotFound(http.StatusNotFound, "diary folder not found")
	}

	return nil
}
