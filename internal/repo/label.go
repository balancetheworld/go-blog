package repo

import (
	"context"
	"errors"

	"github.com/zyj/my-blog/internal/model"
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
