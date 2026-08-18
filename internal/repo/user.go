package repo

import (
	"context"
	"errors"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
)

func RegisterUserWithSession(
	ctx context.Context,
	user *model.User,
	session *model.Session,
	requestedRoleID *uint,
	beforeCommit func() error,
) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		roleCode := constant.RoleCodeMember
		legacyRole := constant.RoleUser

		var requestedRole *model.Role

		if requestedRoleID != nil {
			var role model.Role
			err := tx.
				Where(
					"id = ? AND enabled = ? AND is_requestable = ?",
					*requestedRoleID,
					true,
					true,
				).
				First(&role).
				Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRoleNotRequestable
			}
			if err != nil {
				return err
			}

			requestedRole = &role
		}

		var currentRole model.Role
		if err := tx.
			Where("code = ? AND enabled = ?", roleCode, true).
			First(&currentRole).
			Error; err != nil {
			return err
		}

		user.Role = legacyRole
		user.RoleID = &currentRole.ID

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		session.UserID = user.ID

		if err := tx.Create(session).Error; err != nil {
			return err
		}

		if requestedRole != nil {
			application := model.RoleApplication{
				UserID:          user.ID,
				RequestedRoleID: requestedRole.ID,
				Status:          constant.RoleApplicationPending,
			}

			if err := tx.Create(&application).Error; err != nil {
				return err
			}
		}

		if beforeCommit != nil {
			if err := beforeCommit(); err != nil {
				return err
			}
		}

		return nil
	})
}

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
	err := db.WithContext(ctx).
		Preload("CurrentRole").
		First(&user, id).
		Error
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

func GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	if db == nil {
		return model.User{}, errors.New("database is not initialized")
	}

	var user model.User         // 1. 声明变量，准备存放查到的数据
	err := db.WithContext(ctx). //GORM 的方法：**复制一份全新的 DB 会话，绑定传入的 ctx**。
					Where("username = ?", username).
					First(&user). //链式查询，`First()`：查询第一条匹配的数据，自动加 `ORDER BY id LIMIT 1`，`&user` 传变量地址：GORM 查到数据后，直接把数据库字段赋值给 user 对象。
					Error

	return user, err
}

func GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	if db == nil {
		return model.User{}, errors.New("database is not initialized")
	}

	var user model.User
	err := db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error
	return user, err
}

func GetUserByGithubID(ctx context.Context, githubID uint64) (model.User, error) {
	if db == nil {
		return model.User{}, errors.New("database is not initialized")
	}

	var user model.User
	err := db.WithContext(ctx).
		Where("github_id = ?", githubID).
		First(&user).
		Error
	return user, err
}

func BindGithubID(ctx context.Context, userID uint, githubID uint64) (bool, error) {
	if db == nil {
		return false, errors.New("database is not initialized")
	}

	result := db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ? AND github_id IS NULL", userID).
		Update("github_id", githubID)
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected == 1, nil
}

func GetUserByUsernameOrEmail(ctx context.Context, value string) (model.User, error) {
	if db == nil {
		return model.User{}, errors.New("database is not initialized")
	}

	var user model.User
	err := db.WithContext(ctx).
		Where("username = ? OR email = ?", value, value).
		First(&user).
		Error

	return user, err
}

func CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	if db == nil {
		return false, errors.New("database is not initialized")
	}
	var count int64
	err := db.WithContext(ctx).
		Model(&model.User{}). //指定本次操作对应的数据库表为 `users`
		Where("username = ?", username).
		Count(&count). //统计满足条件的记录总条数
		Error

	return count > 0, err
}

func CheckEmailExists(ctx context.Context, email string) (bool, error) {
	if db == nil {
		return false, errors.New("database is not initialized")
	}

	var count int64
	err := db.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		Count(&count).
		Error

	return count > 0, err
}

func CountUsers(ctx context.Context) (int64, error) {
	if db == nil {
		return 0, errors.New("database is not initialized")
	}

	var count int64
	err := db.WithContext(ctx).
		Model(&model.User{}).
		Count(&count).
		Error
	return count, err
}
