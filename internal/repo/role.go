package repo

import (
	"context"
	"errors"
	"strings"

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

func GetRoleByID(
	ctx context.Context,
	id uint,
) (model.Role, error) {
	if db == nil {
		return model.Role{}, errors.New("database is not initialized")
	}

	var role model.Role
	err := db.WithContext(ctx).
		First(&role, id).
		Error

	return role, err
}

func GetEnabledRolesByIDs(
	ctx context.Context,
	ids []uint,
) ([]model.Role, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	if len(ids) == 0 {
		return []model.Role{}, nil
	}

	var roles []model.Role
	err := db.WithContext(ctx).
		Where("id IN ? AND enabled = ?", ids, true).
		Find(&roles).
		Error

	return roles, err
}

func ListEnabledRoles(ctx context.Context) ([]model.Role, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	var roles []model.Role
	err := db.WithContext(ctx).
		Where("enabled = ?", true).
		Order("is_system DESC, id ASC").
		Find(&roles).
		Error

	return roles, err
}

func CheckRoleCodeExists(
	ctx context.Context,
	code string,
) (bool, error) {
	if db == nil {
		return false, errors.New("database is not initialized")
	}

	var count int64
	err := db.WithContext(ctx).
		Unscoped().
		Model(&model.Role{}).
		Where("code = ?", code).
		Count(&count).
		Error

	return count > 0, err
}

func ListRoles(
	ctx context.Context,
	keyword string,
	offset int,
	limit int,
) ([]model.Role, int64, error) {
	if db == nil {
		return nil, 0, errors.New("database is not initialized")
	}

	query := db.WithContext(ctx).Model(&model.Role{})
	if keyword = strings.TrimSpace(keyword); keyword != "" {
		pattern := "%" + strings.ToLower(keyword) + "%"
		query = query.Where(
			"LOWER(code) LIKE ? OR LOWER(name) LIKE ?",
			pattern,
			pattern,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var roles []model.Role
	err := query.
		Order("is_system DESC, id ASC").
		Offset(offset).
		Limit(limit).
		Find(&roles).
		Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func CreateRole(ctx context.Context, role *model.Role) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Create(role).Error
}

func UpdateRole(ctx context.Context, role model.Role) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).
		Model(&model.Role{}).
		Where("id = ?", role.ID).
		Updates(map[string]any{
			"name":           role.Name,
			"description":    role.Description,
			"is_requestable": role.IsRequestable,
			"enabled":        role.Enabled,
		}).
		Error
}

func CountRoleReferences(
	ctx context.Context,
	roleID uint,
) (int64, int64, error) {
	if db == nil {
		return 0, 0, errors.New("database is not initialized")
	}

	var userCount int64
	if err := db.WithContext(ctx).
		Unscoped().
		Model(&model.User{}).
		Where("role_id = ?", roleID).
		Count(&userCount).
		Error; err != nil {
		return 0, 0, err
	}

	var applicationCount int64
	if err := db.WithContext(ctx).
		Unscoped().
		Model(&model.RoleApplication{}).
		Where("requested_role_id = ?", roleID).
		Count(&applicationCount).
		Error; err != nil {
		return 0, 0, err
	}

	return userCount, applicationCount, nil
}

func DeleteRole(ctx context.Context, id uint) (int64, error) {
	if db == nil {
		return 0, errors.New("database is not initialized")
	}

	result := db.WithContext(ctx).Delete(&model.Role{}, id)
	return result.RowsAffected, result.Error
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
