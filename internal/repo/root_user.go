package repo

import (
	"context"
	"errors"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
)

func GetRootUser(ctx context.Context) (model.User, error) {
	if db == nil {
		return model.User{}, errors.New("database is not initialized")
	}

	var user model.User
	err := db.WithContext(ctx).
		Unscoped().
		Where("is_root = ?", true).
		First(&user).
		Error
	return user, err
}

func RestoreRootUser(ctx context.Context, userID uint, adminRoleID uint) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).
		Unscoped().
		Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"role":       constant.RoleAdmin,
			"role_id":    adminRoleID,
			"is_root":    true,
			"deleted_at": nil,
		}).
		Error
}
