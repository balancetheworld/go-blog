package repo

import (
	"context"
	"errors"

	"github.com/zyj/my-blog/internal/model"
	"gorm.io/gorm"
)

func GetLabelsByIDs(ctx context.Context, ids []uint) ([]model.Label, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	if len(ids) == 0 {
		return []model.Label{}, nil
	}

	var labels []model.Label
	err := db.WithContext(ctx).
		Where("id IN ?", ids).
		Order("id ASC").
		Find(&labels).
		Error

	return labels, err
}

func ListLabels(ctx context.Context) ([]model.Label, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	var labels []model.Label
	err := db.WithContext(ctx).
		Order("name ASC").
		Find(&labels).
		Error

	return labels, err
}

func CreateLabel(ctx context.Context, label *model.Label) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Create(label).Error
}

func GetLabelByID(ctx context.Context, id uint) (model.Label, error) {
	if db == nil {
		return model.Label{}, errors.New("database is not initialized")
	}

	var label model.Label
	err := db.WithContext(ctx).First(&label, id).Error

	return label, err
}

func UpdateLabel(ctx context.Context, label *model.Label) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Save(label).Error
}

func DeleteLabel(ctx context.Context, id uint) (int64, error) {
	if db == nil {
		return 0, errors.New("database is not initialized")
	}

	var rowsAffected int64
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM post_labels WHERE label_id = ?", id).Error; err != nil {
			return err
		}

		result := tx.Unscoped().Delete(&model.Label{}, id)
		rowsAffected = result.RowsAffected
		return result.Error
	})

	return rowsAffected, err
}

func CheckLabelExists(
	ctx context.Context,
	name string,
	slug string,
	excludeID uint,
) (bool, error) {
	if db == nil {
		return false, errors.New("database is not initialized")
	}

	query := db.WithContext(ctx).
		Model(&model.Label{}).
		Where("name = ? OR slug = ?", name, slug)

	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
