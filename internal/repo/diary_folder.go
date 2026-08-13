package repo

import (
	"context"
	"errors"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
)

type fixedDiaryFolder struct {
	Name string
	Slug string
	Sort int
}

var fixedDiaryFolders = []fixedDiaryFolder{
	{Name: "旅行", Slug: "travel", Sort: 1},
	{Name: "日常", Slug: "daily", Sort: 2},
	{Name: "灵感", Slug: "inspiration", Sort: 3},
	{Name: "小记", Slug: "notes", Sort: 4},
	{Name: "收藏", Slug: "favorites", Sort: 5},
}

var fixedDiaryFolderSlugs = []string{
	"travel",
	"daily",
	"inspiration",
	"notes",
	"favorites",
}

func EnsureFixedDiaryFolders(ctx context.Context) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		folderIDs := make([]uint, 0, len(fixedDiaryFolders))
		for _, item := range fixedDiaryFolders {
			var folder model.DiaryFolder
			err := tx.Where("slug = ?", item.Slug).First(&folder).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				folder = model.DiaryFolder{
					Name:       item.Name,
					Slug:       item.Slug,
					Sort:       item.Sort,
					Visibility: constant.PostVisibilityPublic,
				}
				if err := tx.Create(&folder).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else if err := tx.Model(&folder).Updates(map[string]any{
				"name":       item.Name,
				"sort":       item.Sort,
				"visibility": constant.PostVisibilityPublic,
			}).Error; err != nil {
				return err
			}
			if err := tx.Model(&folder).Association("VisibleRoles").Clear(); err != nil {
				return err
			}
			folderIDs = append(folderIDs, folder.ID)
		}

		return tx.Model(&model.Diary{}).
			Where("folder_id IS NOT NULL AND folder_id NOT IN ?", folderIDs).
			UpdateColumn("folder_id", nil).
			Error
	})
}

func CreateDiaryFolder(ctx context.Context, folder *model.DiaryFolder) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("VisibleRoles").Create(folder).Error; err != nil {
			return err
		}

		return tx.Model(folder).Association("VisibleRoles").Replace(folder.VisibleRoles)
	})
}

func UpdateDiaryFolder(ctx context.Context, folder *model.DiaryFolder) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("VisibleRoles").Save(folder).Error; err != nil {
			return err
		}

		return tx.Model(folder).Association("VisibleRoles").Replace(folder.VisibleRoles)
	})
}

func DeleteDiaryFolder(ctx context.Context, id uint) (int64, error) {
	if db == nil {
		return 0, errors.New("database is not initialized")
	}

	var rowsAffected int64
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Diary{}).
			Where("folder_id = ?", id).
			UpdateColumn("folder_id", nil).
			Error; err != nil {
			return err
		}

		folder := model.DiaryFolder{}
		folder.ID = id
		result := tx.Select("VisibleRoles").Delete(&folder)
		rowsAffected = result.RowsAffected
		return result.Error
	})

	return rowsAffected, err
}

func GetDiaryFolderByID(ctx context.Context, id uint) (model.DiaryFolder, error) {
	if db == nil {
		return model.DiaryFolder{}, errors.New("database is not initialized")
	}

	var folder model.DiaryFolder
	err := db.WithContext(ctx).Preload("VisibleRoles").First(&folder, id).Error
	return folder, err
}

func CheckDiaryFolderSlugExists(ctx context.Context, slug string, excludeID uint) (bool, error) {
	if db == nil {
		return false, errors.New("database is not initialized")
	}

	query := db.WithContext(ctx).Model(&model.DiaryFolder{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

func ListDiaryFolders(ctx context.Context) ([]model.DiaryFolder, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	var folders []model.DiaryFolder
	err := db.WithContext(ctx).
		Preload("VisibleRoles").
		Where("slug IN ?", fixedDiaryFolderSlugs).
		Order("sort ASC, created_at ASC").
		Find(&folders).
		Error

	return folders, err
}
