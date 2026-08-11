package repo

import (
	"context"
	"errors"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
)

func GetRoleByCode(ctx context.Context, code string) (model.Role, error) {
	if db == nil {
		return model.Role{}, errors.New("database is not initialized")
	}
	var role model.Role
	err := db.WithContext(ctx).
		Where("code = ? AND enabled = ?", code, true).
		First(&role).
		Error
	return role, err
}

func EnsureSystemRoles(ctx context.Context) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	roles := []model.Role{
		{
			Code:          constant.RoleCodeMember,
			Name:          "普通访客",
			IsSystem:      true,
			IsDefault:     true,
			IsRequestable: false,
			Enabled:       true,
		},
		{
			Code:          constant.RoleCodeEditor,
			Name:          "编辑者",
			IsSystem:      true,
			IsDefault:     false,
			IsRequestable: true,
			Enabled:       true,
		},
		{
			Code:          constant.RoleCodeAdmin,
			Name:          "管理员",
			IsSystem:      true,
			IsDefault:     false,
			IsRequestable: false,
			Enabled:       true,
		},
	}
	for i := range roles {
		if err := db.WithContext(ctx).
			Where("code = ?", roles[i].Code).
			FirstOrCreate(&roles[i]).
			Error; err != nil {
			return err
		}
	}

	return nil
}

func BackfillUserRoleIDs(ctx context.Context) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	mappings := []struct {
		LegacyRole constant.Role
		RoleCode   string
	}{
		{constant.RoleUser, constant.RoleCodeMember},
		{constant.RoleEditor, constant.RoleCodeEditor},
		{constant.RoleAdmin, constant.RoleCodeAdmin},
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, mapping := range mappings {
			var role model.Role
			if err := tx.Where("code = ?", mapping.RoleCode).
				First(&role).
				Error; err != nil {
				return err
			}

			if err := tx.Model(&model.User{}).
				Where("role = ? AND role_id IS NULL", mapping.LegacyRole).
				Update("role_id", role.ID).
				Error; err != nil {
				return err
			}
		}

		return nil
	})
}

var ErrRoleNotRequestable = errors.New("role is not requestable")

func GetRequestableRoleByID(
	ctx context.Context,
	id uint,
) (model.Role, error) {
	if db == nil {
                return model.Role{}, errors.New("database is not initialized")
        } 
	var role model.Role
	 err := db.WithContext(ctx).
                Where(
                        "id = ? AND enabled = ? AND is_requestable = ?",
                        id,
                        true,
                        true,
                ).
                First(&role).
                Error
        return role, err
}

func ListRequestableRoles(
        ctx context.Context,
  ) ([]model.Role, error) {
        if db == nil {
                return nil, errors.New("database is not initialized")
        }

        var roles []model.Role
        err := db.WithContext(ctx).
                Where("enabled = ? AND is_requestable = ?", true, true).
                Order("id ASC").
                Find(&roles).
                Error
        return roles, err
  }
