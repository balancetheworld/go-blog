package repo

import (
	"context"
	"errors"
	"time"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func GetLatestRoleApplicationByUserID(
	ctx context.Context,
	userID uint,
) (model.RoleApplication, error) {
	if db == nil {
		return model.RoleApplication{}, errors.New(
			"database is not initialized",
		)
	}

	var application model.RoleApplication
	err := db.WithContext(ctx).
		Preload("RequestedRole").
		Where("user_id = ?", userID).
		Order("id DESC").
		First(&application).
		Error

	return application, err
}

var ErrRoleApplicationNotPending = errors.New(
	"role application is not pending",
)

func ApproveRoleApplication(
	ctx context.Context,
	applicationID uint,
	reviewerID uint,
) error {
	if db == nil {
		return errors.New("database is not initialized")
	}
	//.Transaction(回调函数)开启数据库事务。  > 事务：这一块里面多条数据库操作，要么全部成功提交，要么全部回滚，不会出现一半成功一半失败
	//`tx`：事务内的数据库会话对象，事务里所有增删改查都必须用 `tx`，不能写 `db.xxx`。
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var application model.RoleApplication
		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&application, applicationID).
			Error
		if err != nil {
			return err
		}

		if application.Status != constant.RoleApplicationPending {
			return ErrRoleApplicationNotPending
		}

		var requestedRole model.Role
		err = tx.
			First(&requestedRole, application.RequestedRoleID).
			Error
		if err != nil {
			return err
		}

		if !requestedRole.Enabled {
			return ErrRoleNotRequestable
		}

		legacyRole := constant.RoleUser
		if requestedRole.Code == constant.RoleCodeEditor {
			legacyRole = constant.RoleEditor
		}
		if requestedRole.Code == constant.RoleCodeAdmin {
			return ErrRoleNotRequestable
		}

		result := tx.
			Model(&model.User{}).
			Where(
				"id = ? AND is_root = ?",
				application.UserID,
				false,
			).
			Updates(map[string]any{
				"role":    legacyRole,
				"role_id": requestedRole.ID,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		now := time.Now()
		return tx.
			Model(&application).
			Updates(map[string]any{
				"status":        constant.RoleApplicationApproved,
				"reviewer_id":   reviewerID,
				"reviewed_at":   &now,
				"reject_reason": "",
			}).
			Error
	})
}

 func RejectRoleApplication(
        ctx context.Context,
        applicationID uint,
        reviewerID uint,
        reason string,
  ) error {
        if db == nil {
                return errors.New("database is not initialized")
        }

        return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
                var application model.RoleApplication
                err := tx.
                        Clauses(clause.Locking{Strength: "UPDATE"}).
                        First(&application, applicationID).
                        Error
                if err != nil {
                        return err
                }

                if application.Status != constant.RoleApplicationPending {
                        return ErrRoleApplicationNotPending
                }

                now := time.Now()
                return tx.
                        Model(&application).
                        Updates(map[string]any{
                                "status":        constant.RoleApplicationRejected,
                                "reviewer_id":   reviewerID,
                                "reviewed_at":   &now,
                                "reject_reason": reason,
                        }).
                        Error
        })
  }
