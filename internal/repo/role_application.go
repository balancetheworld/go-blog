package repo

import (
	"context"
	"errors"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
)

func ListRoleApplications(
	ctx context.Context,
	status constant.RoleApplicationStatus,
	offset int,
	limit int,
) ([]model.RoleApplication, int64, error) {
	if db == nil {
		return nil, 0, errors.New("database is not initialized")
	}
	database := db.WithContext(ctx).
		Model(&model.RoleApplication{})

	if status != "" {
		database = database.Where("status = ?", status)
	}

	var total int64
	if err := database.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var applications []model.RoleApplication
	err := database.
		Preload("User").
		Preload("RequestedRole").
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&applications).
		Error
	if err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}
