package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/errs"
	"gorm.io/gorm"
)

func canManageDiary(diary model.Diary, viewerID uint, role constant.Role) bool {
	if role == constant.RoleAdmin {
		return true
	}

	return role == constant.RoleEditor && viewerID > 0 && diary.AuthorID == viewerID
}

func canViewDiaryFolder(
	folder *model.DiaryFolder,
	role constant.Role,
	viewerRoleID uint,
) bool {
	if folder == nil || role == constant.RoleAdmin || role == constant.RoleEditor {
		return true
	}

	switch folder.Visibility {
	case constant.PostVisibilityPublic, "":
		return true
	case constant.PostVisibilityRoles:
		for _, visibleRole := range folder.VisibleRoles {
			if visibleRole.ID == viewerRoleID {
				return true
			}
		}
	}

	return false
}

func canViewDiary(
	diary model.Diary,
	viewerID uint,
	role constant.Role,
	viewerRoleID uint,
) bool {
	if canManageDiary(diary, viewerID, role) {
		return true
	}
	if diary.PublishedAt == nil || !canViewDiaryFolder(diary.Folder, role, viewerRoleID) {
		return false
	}

	switch diary.Visibility {
	case constant.PostVisibilityPublic, "":
		return true
	case constant.PostVisibilityRoles:
		if viewerRoleID == 0 {
			return false
		}

		for _, visibleRole := range diary.VisibleRoles {
			if visibleRole.ID == viewerRoleID {
				return true
			}
		}
	}

	return false
}

func diaryFolderToResponse(folder model.DiaryFolder) dto.DiaryFolderResponse {
	visibleRoles := make([]dto.RoleOptionResponse, 0, len(folder.VisibleRoles))
	for _, role := range folder.VisibleRoles {
		visibleRoles = append(visibleRoles, dto.RoleOptionResponse{
			ID:          role.ID,
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
		})
	}

	return dto.DiaryFolderResponse{
		ID:           uint64(folder.ID),
		Name:         folder.Name,
		Slug:         folder.Slug,
		Description:  folder.Description,
		Cover:        folder.Cover,
		Sort:         folder.Sort,
		Visibility:   folder.Visibility,
		VisibleRoles: visibleRoles,
		CreatedAt:    folder.CreatedAt,
		UpdatedAt:    folder.UpdatedAt,
	}
}

func diaryToResponse(diary model.Diary, includeDraft bool) dto.DiaryResponse {
	visibleRoles := make([]dto.RoleOptionResponse, 0, len(diary.VisibleRoles))
	for _, role := range diary.VisibleRoles {
		visibleRoles = append(visibleRoles, dto.RoleOptionResponse{
			ID:          role.ID,
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
		})
	}

	status := "published"
	if diary.PublishedAt == nil {
		status = "draft"
	}

	var folder *dto.DiaryFolderResponse
	if diary.Folder != nil {
		value := diaryFolderToResponse(*diary.Folder)
		folder = &value
	}

	result := dto.DiaryResponse{
		ID:           uint64(diary.ID),
		Title:        diary.Title,
		Slug:         diary.Slug,
		Description:  diary.Description,
		Cover:        diary.Cover,
		Content:      diary.Content,
		Author:       diary.Author.ToDto(),
		Folder:       folder,
		Visibility:   diary.Visibility,
		VisibleRoles: visibleRoles,
		ViewCount:    diary.ViewCount,
		LikeCount:    diary.LikeCount,
		CommentCount: diary.CommentCount,
		Status:       status,
		PublishedAt:  diary.PublishedAt,
		CreatedAt:    diary.CreatedAt,
		UpdatedAt:    diary.UpdatedAt,
	}
	if includeDraft {
		draftContent := diary.DraftContent
		result.DraftContent = &draftContent
	}

	return result
}

func GetDiary(
	ctx context.Context,
	id uint,
	viewerID uint,
	role constant.Role,
	viewerRoleID uint,
) (dto.DiaryResponse, error) {
	return getDiaryResponse(ctx, func() (model.Diary, error) {
		return repo.GetDiaryByID(ctx, id)
	}, viewerID, role, viewerRoleID)
}

func GetDiaryByIdentifier(
	ctx context.Context,
	identifier string,
	viewerID uint,
	role constant.Role,
	viewerRoleID uint,
) (dto.DiaryResponse, error) {
	return getDiaryResponse(ctx, func() (model.Diary, error) {
		if id, err := strconv.ParseUint(identifier, 10, 64); err == nil && id > 0 {
			return repo.GetDiaryByID(ctx, uint(id))
		}

		return repo.GetDiaryBySlug(ctx, identifier)
	}, viewerID, role, viewerRoleID)
}

func getDiaryResponse(
	ctx context.Context,
	find func() (model.Diary, error),
	viewerID uint,
	role constant.Role,
	viewerRoleID uint,
) (dto.DiaryResponse, error) {
	diary, err := find()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.DiaryResponse{}, errs.NewNotFound(http.StatusNotFound, "diary not found")
	}
	if err != nil {
		return dto.DiaryResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "get diary failed")
	}
	if !canViewDiary(diary, viewerID, role, viewerRoleID) {
		return dto.DiaryResponse{}, errs.NewNotFound(http.StatusNotFound, "diary not found")
	}

	if diary.PublishedAt != nil && !canManageDiary(diary, viewerID, role) {
		if err := repo.IncreaseDiaryViewCount(ctx, diary.ID); err == nil {
			diary.ViewCount++
		}
	}

	return diaryToResponse(diary, canManageDiary(diary, viewerID, role)), nil
}

func ListDiaries(
	ctx context.Context,
	req dto.ListDiariesRequest,
	viewerID uint,
	role constant.Role,
	viewerRoleID uint,
) (dto.ListDiariesResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	filter := repo.DiaryListFilter{
		Offset:       (req.Page - 1) * req.PageSize,
		Limit:        req.PageSize,
		Status:       "published",
		Keyword:      req.Keyword,
		FolderID:     req.FolderID,
		PublicOnly:   true,
		ViewerID:     viewerID,
		ViewerRoleID: viewerRoleID,
	}
	if req.Status == "draft" || req.Status == "all" {
		switch role {
		case constant.RoleAdmin:
			filter.PublicOnly = false
			filter.Status = req.Status
		case constant.RoleEditor:
			if viewerID == 0 {
				return dto.ListDiariesResponse{}, errs.NewForbidden(http.StatusForbidden, "diary list access denied")
			}
			filter.PublicOnly = false
			filter.Status = req.Status
			filter.AuthorID = viewerID
		default:
			return dto.ListDiariesResponse{}, errs.NewForbidden(http.StatusForbidden, "diary list access denied")
		}
	}

	diaries, total, err := repo.ListDiaries(ctx, filter)
	if err != nil {
		return dto.ListDiariesResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "list diaries failed")
	}

	items := make([]dto.DiaryResponse, 0, len(diaries))
	for _, diary := range diaries {
		items = append(items, diaryToResponse(diary, canManageDiary(diary, viewerID, role)))
	}

	return dto.ListDiariesResponse{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

var invalidDiarySlugCharacters = regexp.MustCompile(`[^\p{L}\p{N}]+`)

func normalizeDiarySlug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = invalidDiarySlugCharacters.ReplaceAllString(value, "-")
	return strings.Trim(value, "-")
}

func resolveDiarySlug(ctx context.Context, value string, title string, excludeID uint) (string, error) {
	base := normalizeDiarySlug(value)
	if base == "" {
		base = normalizeDiarySlug(title)
	}
	if base == "" || regexp.MustCompile(`^\d+$`).MatchString(base) {
		base = "diary-" + time.Now().Format("20060102-150405")
	}

	for sequence := 1; sequence <= 1000; sequence++ {
		slug := base
		if sequence > 1 {
			slug = fmt.Sprintf("%s-%d", base, sequence)
		}
		exists, err := repo.CheckDiarySlugExists(ctx, slug, excludeID)
		if err != nil {
			return "", errs.NewInternalServer(http.StatusInternalServerError, "check diary slug failed")
		}
		if !exists {
			return slug, nil
		}
	}

	return "", errs.NewConflict(http.StatusConflict, "diary slug already exists")
}

func resolveDiaryFolder(ctx context.Context, folderID *uint) (*model.DiaryFolder, error) {
	if folderID == nil || *folderID == 0 {
		return nil, nil
	}

	folder, err := repo.GetDiaryFolderByID(ctx, *folderID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.NewBadRequest(http.StatusBadRequest, "diary folder not found")
	}
	if err != nil {
		return nil, errs.NewInternalServer(http.StatusInternalServerError, "get diary folder failed")
	}

	return &folder, nil
}

func CreateDiary(
	ctx context.Context,
	authorID uint,
	role constant.Role,
	req dto.CreateDiaryRequest,
) (dto.DiaryResponse, error) {
	if authorID == 0 || (role != constant.RoleEditor && role != constant.RoleAdmin) {
		return dto.DiaryResponse{}, errs.NewForbidden(http.StatusForbidden, "create diary access denied")
	}

	visibility := req.Visibility
	if visibility == "" {
		visibility = constant.PostVisibilityPublic
	}
	visibleRoles, err := resolvePostVisibleRoles(ctx, visibility, req.VisibleRoleIDs)
	if err != nil {
		return dto.DiaryResponse{}, err
	}
	folder, err := resolveDiaryFolder(ctx, req.FolderID)
	if err != nil {
		return dto.DiaryResponse{}, err
	}
	slug, err := resolveDiarySlug(ctx, req.Slug, req.Title, 0)
	if err != nil {
		return dto.DiaryResponse{}, err
	}

	draftContent := req.DraftContent
	content := ""
	var publishedAt *time.Time
	if req.Publish {
		if strings.TrimSpace(draftContent) == "" {
			return dto.DiaryResponse{}, errs.NewBadRequest(http.StatusBadRequest, "published diary content is required")
		}
		content = draftContent
		now := time.Now()
		publishedAt = &now
	}

	diary := model.Diary{
		Title:        strings.TrimSpace(req.Title),
		Slug:         slug,
		Description:  strings.TrimSpace(req.Description),
		Cover:        strings.TrimSpace(req.Cover),
		Content:      content,
		DraftContent: draftContent,
		AuthorID:     authorID,
		FolderID:     req.FolderID,
		Folder:       folder,
		Visibility:   visibility,
		VisibleRoles: visibleRoles,
		PublishedAt:  publishedAt,
	}
	if err := repo.CreateDiary(ctx, &diary); err != nil {
		return dto.DiaryResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "create diary failed")
	}

	createdDiary, err := repo.GetDiaryByID(ctx, diary.ID)
	if err != nil {
		return dto.DiaryResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "get created diary failed")
	}

	return diaryToResponse(createdDiary, true), nil
}

func UpdateDiary(
	ctx context.Context,
	id uint,
	viewerID uint,
	role constant.Role,
	req dto.UpdateDiaryRequest,
) (dto.DiaryResponse, error) {
	diary, err := repo.GetDiaryByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.DiaryResponse{}, errs.NewNotFound(http.StatusNotFound, "diary not found")
	}
	if err != nil {
		return dto.DiaryResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "get diary failed")
	}
	if !canManageDiary(diary, viewerID, role) {
		return dto.DiaryResponse{}, errs.NewForbidden(http.StatusForbidden, "update diary access denied")
	}

	if req.Title != nil {
		diary.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		diary.Description = strings.TrimSpace(*req.Description)
	}
	if req.Cover != nil {
		diary.Cover = strings.TrimSpace(*req.Cover)
	}
	if req.DraftContent != nil {
		diary.DraftContent = *req.DraftContent
	}
	if req.ClearFolder {
		diary.FolderID = nil
		diary.Folder = nil
	} else if req.FolderID != nil {
		folder, folderErr := resolveDiaryFolder(ctx, req.FolderID)
		if folderErr != nil {
			return dto.DiaryResponse{}, folderErr
		}
		diary.FolderID = req.FolderID
		diary.Folder = folder
	}
	if req.Slug != nil || diary.Slug == "" {
		slugValue := diary.Slug
		if req.Slug != nil {
			slugValue = *req.Slug
		}
		diary.Slug, err = resolveDiarySlug(ctx, slugValue, diary.Title, diary.ID)
		if err != nil {
			return dto.DiaryResponse{}, err
		}
	}

	visibility := diary.Visibility
	if req.Visibility != nil {
		visibility = *req.Visibility
	}
	visibleRoleIDs := make([]uint, 0, len(diary.VisibleRoles))
	for _, visibleRole := range diary.VisibleRoles {
		visibleRoleIDs = append(visibleRoleIDs, visibleRole.ID)
	}
	if req.VisibleRoleIDs != nil {
		visibleRoleIDs = *req.VisibleRoleIDs
	}
	visibleRoles, err := resolvePostVisibleRoles(ctx, visibility, visibleRoleIDs)
	if err != nil {
		return dto.DiaryResponse{}, err
	}
	diary.Visibility = visibility
	diary.VisibleRoles = visibleRoles

	if req.Publish != nil {
		if *req.Publish {
			if strings.TrimSpace(diary.DraftContent) == "" {
				return dto.DiaryResponse{}, errs.NewBadRequest(http.StatusBadRequest, "published diary content is required")
			}
			diary.Content = diary.DraftContent
			if diary.PublishedAt == nil {
				now := time.Now()
				diary.PublishedAt = &now
			}
		} else {
			diary.Content = ""
			diary.PublishedAt = nil
		}
	}

	if err := repo.UpdateDiary(ctx, &diary); err != nil {
		return dto.DiaryResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "update diary failed")
	}

	updatedDiary, err := repo.GetDiaryByID(ctx, diary.ID)
	if err != nil {
		return dto.DiaryResponse{}, errs.NewInternalServer(http.StatusInternalServerError, "get updated diary failed")
	}

	return diaryToResponse(updatedDiary, true), nil
}

func DeleteDiary(
	ctx context.Context,
	id uint,
	viewerID uint,
	role constant.Role,
) error {
	diary, err := repo.GetDiaryByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.NewNotFound(http.StatusNotFound, "diary not found")
	}
	if err != nil {
		return errs.NewInternalServer(http.StatusInternalServerError, "get diary failed")
	}
	if !canManageDiary(diary, viewerID, role) {
		return errs.NewForbidden(http.StatusForbidden, "delete diary access denied")
	}

	rowsAffected, err := repo.DeleteDiary(ctx, id)
	if err != nil {
		return errs.NewInternalServer(http.StatusInternalServerError, "delete diary failed")
	}
	if rowsAffected == 0 {
		return errs.NewNotFound(http.StatusNotFound, "diary not found")
	}

	return nil
}
