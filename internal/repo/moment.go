package repo

import (
	"context"
	"errors"

	"github.com/zyj/my-blog/internal/model"
)

func CreateMoment(ctx context.Context, moment *model.Moment) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Omit("Author").Create(moment).Error
}

func GetMomentByID(ctx context.Context, id uint) (model.Moment, error) {
	if db == nil {
		return model.Moment{}, errors.New("database is not initialized")
	}

	var moment model.Moment
	err := db.WithContext(ctx).Preload("Author").First(&moment, id).Error

	return moment, err
}

func ListMoments(ctx context.Context, offset int, limit int) ([]model.Moment, int64, error) {
	if db == nil {
		return nil, 0, errors.New("database is not initialized")
	}

	query := db.WithContext(ctx).Model(&model.Moment{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var moments []model.Moment
	err := query.
		Preload("Author").
		Order("created_at DESC").
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&moments).
		Error

	return moments, total, err
}

func DeleteMoment(ctx context.Context, id uint) (int64, error) {
	if db == nil {
		return 0, errors.New("database is not initialized")
	}

	result := db.WithContext(ctx).Delete(&model.Moment{}, id)
	return result.RowsAffected, result.Error
}
