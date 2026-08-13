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
		Where(
			"id IN ? AND enabled = ? AND code NOT IN ?",
			ids,
			true,
			[]string{
				constant.RoleCodeGuest,
				constant.RoleCodeAdmin,
				constant.RoleCodeEditor,
			},
		).
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
		Where(
			"enabled = ? AND code NOT IN ?",
			true,
			[]string{
				constant.RoleCodeGuest,
				constant.RoleCodeAdmin,
				constant.RoleCodeEditor,
			},
		).
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

	query := db.WithContext(ctx).
		Model(&model.Role{}).
		Where("code <> ?", constant.RoleCodeEditor)
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
			Code:          constant.RoleCodeGuest,
			Name:          "游客",
			Description:   "未登录访问站点时使用的身份",
			IsSystem:      true,
			IsDefault:     false,
			IsRequestable: false,
			Enabled:       true,
		},
		{
			Code:          constant.RoleCodeMember,
			Name:          "普通访客",
			Description:   "用户登录后的默认身份",
			IsSystem:      true,
			IsDefault:     true,
			IsRequestable: false,
			Enabled:       true,
		},
		{
			Code:          constant.RoleCodeAdmin,
			Name:          "管理员",
			Description:   "博客主人专用的最高权限身份",
			IsSystem:      true,
			IsDefault:     false,
			IsRequestable: false,
			Enabled:       true,
		},
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range roles {
			var role model.Role
			err := tx.Unscoped().Where("code = ?", roles[i].Code).First(&role).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := tx.Create(&roles[i]).Error; err != nil {
					return err
				}
				continue
			}
			if err != nil {
				return err
			}

			if err := tx.Unscoped().Model(&role).Updates(map[string]any{
				"name":           roles[i].Name,
				"description":    roles[i].Description,
				"is_system":      roles[i].IsSystem,
				"is_default":     roles[i].IsDefault,
				"is_requestable": roles[i].IsRequestable,
				"enabled":        roles[i].Enabled,
				"deleted_at":     nil,
			}).Error; err != nil {
				return err
			}
		}

		var editorRole model.Role
		err := tx.Unscoped().Where("code = ?", constant.RoleCodeEditor).First(&editorRole).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			editorRole = model.Role{
				Code:          constant.RoleCodeEditor,
				Name:          "编辑者",
				Description:   "已停用的旧版内容编辑身份",
				IsRequestable: false,
				Enabled:       false,
			}
			return tx.Create(&editorRole).Error
		}
		if err != nil {
			return err
		}

		return tx.Unscoped().Model(&editorRole).Updates(map[string]any{
			"description":    "已停用的旧版内容编辑身份",
			"is_system":      false,
			"is_default":     false,
			"is_requestable": false,
			"enabled":        false,
			"deleted_at":     nil,
		}).Error
	})
}

func RetireEditorRole(ctx context.Context) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	var editorRole model.Role
	if err := db.WithContext(ctx).
		Unscoped().
		Where("code = ?", constant.RoleCodeEditor).
		First(&editorRole).
		Error; err != nil {
		return err
	}

	var memberRole model.Role
	if err := db.WithContext(ctx).
		Where("code = ?", constant.RoleCodeMember).
		First(&memberRole).
		Error; err != nil {
		return err
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.User{}).
			Where("role = ? OR role_id = ?", constant.RoleEditor, editorRole.ID).
			Updates(map[string]any{
				"role":    constant.RoleUser,
				"role_id": memberRole.ID,
			}).
			Error; err != nil {
			return err
		}

		var excludedRoleIDs []uint
		if err := tx.Model(&model.Role{}).
			Unscoped().
			Where(
				"code IN ?",
				[]string{
					constant.RoleCodeGuest,
					constant.RoleCodeAdmin,
					constant.RoleCodeEditor,
				},
			).
			Pluck("id", &excludedRoleIDs).
			Error; err != nil {
			return err
		}

		if err := tx.Table("post_visible_roles").
			Where("role_id IN ?", excludedRoleIDs).
			Delete(nil).
			Error; err != nil {
			return err
		}

		return nil
	})
}

func EnforceRootAdminRole(ctx context.Context) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	memberRole, err := GetRoleByCode(ctx, constant.RoleCodeMember)
	if err != nil {
		return err
	}
	adminRole, err := GetRoleByCode(ctx, constant.RoleCodeAdmin)
	if err != nil {
		return err
	}

	return db.WithContext(ctx).
		Model(&model.User{}).
		Where(
			"is_root = ? AND (role = ? OR role_id = ?)",
			false,
			constant.RoleAdmin,
			adminRole.ID,
		).
		Updates(map[string]any{
			"role":    constant.RoleUser,
			"role_id": memberRole.ID,
		}).
		Error
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
