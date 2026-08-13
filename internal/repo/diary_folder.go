package repo

import (
	"context"
	"errors"

	"github.com/zyj/my-blog/internal/model"
	"gorm.io/gorm"
)

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
		Order("sort ASC, created_at ASC").
		Find(&folders).
		Error

	return folders, err
}
