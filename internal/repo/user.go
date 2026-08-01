package repo

import (
	"context"
	"errors"

	"github.com/zyj/my-blog/internal/model"
)

func UserExists(ctx context.Context, username, email string) (bool, error) {
	if db == nil {
		return false, errors.New("database is not initialized")
	}

	var count int64
	err := db.WithContext(ctx).
                Model(&model.User{}).
                Where("username = ? OR email = ?", username, email).
                Count(&count).
                Error
        if err != nil {
                return false, err
        }

        return count > 0, nil
}

func CreateUser(ctx context.Context, user *model.User) error {
	if db == nil {
		return errors.New("database is not initialized")
	}
	// 向user表插入一条新记录，携带上下文ctx，返回数据库错误
	return db.WithContext(ctx).Create(user).Error
}

func ListUsers(ctx context.Context, offset, limit int) ([]model.User, int64, error) {
	if db == nil {
		return nil, 0, errors.New("database is not initialized")
	}

	var users []model.User
	var total int64
	database := db.WithContext(ctx).Model(&model.User{})
	if err := database.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := database.Order("id DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func GetUserByID(ctx context.Context, id uint64) (model.User, error) {
	if db == nil {
		return model.User{}, errors.New("database is not initialized")
	}

	var user model.User
	err := db.WithContext(ctx).First(&user, id).Error
	return user, err
}

func UserExistsExcept(ctx context.Context, id uint64, username, email string) (bool, error) {
	if db == nil {
		return false, errors.New("database is not initialized")
	}

	var count int64
	err := db.WithContext(ctx).
		Model(&model.User{}).
		Where("id <> ? AND (username = ? OR email = ?)", id, username, email).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func UpdateUser(ctx context.Context, user *model.User) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Save(user).Error
}

func DeleteUser(ctx context.Context, id uint64) (int64, error) {
	if db == nil {
		return 0, errors.New("database is not initialized")
	}

	result := db.WithContext(ctx).Delete(&model.User{}, id)
	return result.RowsAffected, result.Error
}
