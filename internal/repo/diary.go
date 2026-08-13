package repo

import (
	"context"
	"errors"
	"strings"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
)

type DiaryListFilter struct {
	Offset       int
	Limit        int
	AuthorID     uint
	Status       string
	Keyword      string
	FolderID     uint
	PublicOnly   bool
	ViewerID     uint
	ViewerRoleID uint
}

func CreateDiary(ctx context.Context, diary *model.Diary) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Author", "Folder", "VisibleRoles").Create(diary).Error; err != nil {
			return err
		}

		return tx.Model(diary).Association("VisibleRoles").Replace(diary.VisibleRoles)
	})
}

func UpdateDiary(ctx context.Context, diary *model.Diary) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Author", "Folder", "VisibleRoles").Save(diary).Error; err != nil {
			return err
		}

		return tx.Model(diary).Association("VisibleRoles").Replace(diary.VisibleRoles)
	})
}

func DeleteDiary(ctx context.Context, id uint) (int64, error) {
	if db == nil {
		return 0, errors.New("database is not initialized")
	}

	var rowsAffected int64
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		diary := model.Diary{}
		diary.ID = id
		if err := tx.Where(
			"target_type = ? AND target_id = ?",
			constant.TargetDiary,
			id,
		).Delete(&model.Comment{}).Error; err != nil {
			return err
		}

		result := tx.Select("VisibleRoles").Delete(&diary)
		rowsAffected = result.RowsAffected
		return result.Error
	})

	return rowsAffected, err
}

func diaryQuery(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx).
		Preload("Author").
		Preload("Folder").
		Preload("Folder.VisibleRoles").
		Preload("VisibleRoles")
}

func GetDiaryByID(ctx context.Context, id uint) (model.Diary, error) {
	if db == nil {
		return model.Diary{}, errors.New("database is not initialized")
	}

	var diary model.Diary
	err := diaryQuery(ctx).First(&diary, id).Error

	return diary, err
}

func GetDiaryBySlug(ctx context.Context, slug string) (model.Diary, error) {
	if db == nil {
		return model.Diary{}, errors.New("database is not initialized")
	}

	var diary model.Diary
	err := diaryQuery(ctx).Where("slug = ?", slug).First(&diary).Error

	return diary, err
}

func CheckDiarySlugExists(ctx context.Context, slug string, excludeID uint) (bool, error) {
	if db == nil {
		return false, errors.New("database is not initialized")
	}

	query := db.WithContext(ctx).Model(&model.Diary{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

func ListDiaries(
	ctx context.Context,
	filter DiaryListFilter,
) ([]model.Diary, int64, error) {
	if db == nil {
		return nil, 0, errors.New("database is not initialized")
	}

	query := db.WithContext(ctx).Model(&model.Diary{})
	if filter.PublicOnly {
		query = query.
			Where("published_at IS NOT NULL").
			Where(
				"visibility = ? OR author_id = ? OR (visibility = ? AND EXISTS (SELECT 1 FROM diary_visible_roles WHERE diary_visible_roles.diary_id = diaries.id AND diary_visible_roles.role_id = ?))",
				constant.PostVisibilityPublic,
				filter.ViewerID,
				constant.PostVisibilityRoles,
				filter.ViewerRoleID,
			).
			Where(
				"folder_id IS NULL OR author_id = ? OR EXISTS (SELECT 1 FROM diary_folders WHERE diary_folders.id = diaries.folder_id AND diary_folders.deleted_at IS NULL AND (diary_folders.visibility = ? OR (diary_folders.visibility = ? AND EXISTS (SELECT 1 FROM diary_folder_visible_roles WHERE diary_folder_visible_roles.diary_folder_id = diary_folders.id AND diary_folder_visible_roles.role_id = ?))))",
				filter.ViewerID,
				constant.PostVisibilityPublic,
				constant.PostVisibilityRoles,
				filter.ViewerRoleID,
			)
	}
	if filter.AuthorID > 0 {
		query = query.Where("author_id = ?", filter.AuthorID)
	}
	if filter.FolderID > 0 {
		query = query.Where("folder_id = ?", filter.FolderID)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		value := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", value, value)
	}
	switch filter.Status {
	case "published":
		query = query.Where("published_at IS NOT NULL")
	case "draft":
		query = query.Where("published_at IS NULL")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var diaries []model.Diary
	err := query.
		Preload("Author").
		Preload("Folder").
		Preload("Folder.VisibleRoles").
		Preload("VisibleRoles").
		Order("created_at DESC").
		Offset(filter.Offset).
		Limit(filter.Limit).
		Find(&diaries).
		Error

	return diaries, total, err
}

func IncreaseDiaryViewCount(ctx context.Context, id uint) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Model(&model.Diary{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).
		Error
}
