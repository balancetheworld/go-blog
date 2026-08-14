package repo

import (
	"context"
	"errors"
	"testing"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRoleManagementRepository(t *testing.T) {
	database, err := gorm.Open(
		sqlite.Open("file:role_repo_test?mode=memory&cache=shared"),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.RoleApplication{},
	); err != nil {
		t.Fatal(err)
	}

	previousDB := db
	db = database
	t.Cleanup(func() {
		db = previousDB
	})

	ctx := context.Background()
	role := model.Role{
		Code:          "photographer",
		Name:          "摄影师",
		IsRequestable: true,
		Enabled:       true,
	}
	if err := CreateRole(ctx, &role); err != nil {
		t.Fatal(err)
	}

	exists, err := CheckRoleCodeExists(ctx, role.Code)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected role code to exist")
	}

	roles, total, err := ListRoles(ctx, "PHOTO", 0, 20)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(roles) != 1 || roles[0].ID != role.ID {
		t.Fatalf("unexpected roles: total=%d roles=%#v", total, roles)
	}

	role.Name = "摄影作者"
	role.IsRequestable = false
	role.Enabled = false
	if err := UpdateRole(ctx, role); err != nil {
		t.Fatal(err)
	}

	updated, err := GetRoleByID(ctx, role.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != role.Name || updated.IsRequestable || updated.Enabled {
		t.Fatalf("unexpected updated role: %#v", updated)
	}

	user := model.User{
		Username:     "role-user",
		Email:        "role-user@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleUser,
		RoleID:       &role.ID,
	}
	if err := database.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	application := model.RoleApplication{
		UserID:          user.ID,
		RequestedRoleID: role.ID,
		Status:          constant.RoleApplicationPending,
	}
	if err := database.Create(&application).Error; err != nil {
		t.Fatal(err)
	}

	userCount, applicationCount, err := CountRoleReferences(ctx, role.ID)
	if err != nil {
		t.Fatal(err)
	}
	if userCount != 1 || applicationCount != 1 {
		t.Fatalf(
			"unexpected references: users=%d applications=%d",
			userCount,
			applicationCount,
		)
	}

	unusedRole := model.Role{
		Code:    "unused",
		Name:    "未使用身份",
		Enabled: true,
	}
	if err := CreateRole(ctx, &unusedRole); err != nil {
		t.Fatal(err)
	}

	rowsAffected, err := DeleteRole(ctx, unusedRole.ID)
	if err != nil {
		t.Fatal(err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 deleted role, got %d", rowsAffected)
	}
	if _, err := GetRoleByID(ctx, unusedRole.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected deleted role to be hidden, got %v", err)
	}

	exists, err = CheckRoleCodeExists(ctx, unusedRole.Code)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected deleted role code to remain reserved")
	}
}

func TestSystemRolePolicy(t *testing.T) {
	database, err := gorm.Open(
		sqlite.Open("file:system_role_policy_test?mode=memory&cache=shared"),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.RoleApplication{},
	); err != nil {
		t.Fatal(err)
	}

	previousDB := db
	db = database
	t.Cleanup(func() {
		db = previousDB
	})

	ctx := context.Background()
	if err := EnsureSystemRoles(ctx); err != nil {
		t.Fatal(err)
	}
	if err := EnsureSystemRoles(ctx); err != nil {
		t.Fatal(err)
	}

	guestRole, err := GetRoleByCode(ctx, constant.RoleCodeGuest)
	if err != nil {
		t.Fatal(err)
	}
	memberRole, err := GetRoleByCode(ctx, constant.RoleCodeMember)
	if err != nil {
		t.Fatal(err)
	}
	var editorRole model.Role
	if err := database.Unscoped().
		Where("code = ?", constant.RoleCodeEditor).
		First(&editorRole).
		Error; err != nil {
		t.Fatal(err)
	}
	adminRole, err := GetRoleByCode(ctx, constant.RoleCodeAdmin)
	if err != nil {
		t.Fatal(err)
	}

	if guestRole.Name != "游客" || !guestRole.IsSystem || guestRole.IsDefault {
		t.Fatalf("unexpected guest role: %#v", guestRole)
	}
	if memberRole.Name != "普通访客" || !memberRole.IsSystem || !memberRole.IsDefault {
		t.Fatalf("unexpected member role: %#v", memberRole)
	}
	if editorRole.IsSystem || editorRole.Enabled || editorRole.IsRequestable {
		t.Fatalf("editor role must be retired: %#v", editorRole)
	}
	if adminRole.Name != "管理员" || !adminRole.IsSystem || adminRole.IsRequestable {
		t.Fatalf("unexpected admin role: %#v", adminRole)
	}

	customRole := model.Role{
		Code:    "friends",
		Name:    "朋友",
		Enabled: true,
	}
	if err := CreateRole(ctx, &customRole); err != nil {
		t.Fatal(err)
	}
	enabledRoles, err := ListEnabledRoles(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(enabledRoles) != 2 || enabledRoles[0].ID != memberRole.ID || enabledRoles[1].ID != customRole.ID {
		t.Fatalf("unexpected visible role options: %#v", enabledRoles)
	}

	rootUser := model.User{
		Username:     "root-user",
		Email:        "root-user@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleAdmin,
		RoleID:       &adminRole.ID,
		IsRoot:       true,
	}
	legacyAdmin := model.User{
		Username:     "legacy-admin",
		Email:        "legacy-admin@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleAdmin,
		RoleID:       &adminRole.ID,
	}
	if err := database.Create(&rootUser).Error; err != nil {
		t.Fatal(err)
	}
	if err := database.Create(&legacyAdmin).Error; err != nil {
		t.Fatal(err)
	}

	if err := EnforceRootAdminRole(ctx); err != nil {
		t.Fatal(err)
	}
	if err := database.First(&rootUser, rootUser.ID).Error; err != nil {
		t.Fatal(err)
	}
	if err := database.First(&legacyAdmin, legacyAdmin.ID).Error; err != nil {
		t.Fatal(err)
	}
	if rootUser.Role != constant.RoleAdmin || rootUser.RoleID == nil || *rootUser.RoleID != adminRole.ID {
		t.Fatalf("root role changed unexpectedly: %#v", rootUser)
	}
	if legacyAdmin.Role != constant.RoleUser || legacyAdmin.RoleID == nil || *legacyAdmin.RoleID != memberRole.ID {
		t.Fatalf("legacy admin was not demoted: %#v", legacyAdmin)
	}
}

func TestDisablingRoleFallsBackUsersAndRevokesSessions(t *testing.T) {
	database, err := gorm.Open(
		sqlite.Open("file:disable_role_policy_test?mode=memory&cache=shared"),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.Session{},
	); err != nil {
		t.Fatal(err)
	}

	previousDB := db
	db = database
	t.Cleanup(func() {
		db = previousDB
	})

	ctx := context.Background()
	if err := EnsureSystemRoles(ctx); err != nil {
		t.Fatal(err)
	}
	memberRole, err := GetRoleByCode(ctx, constant.RoleCodeMember)
	if err != nil {
		t.Fatal(err)
	}
	customRole := model.Role{
		Code:    "verified",
		Name:    "认证用户",
		Enabled: true,
	}
	if err := CreateRole(ctx, &customRole); err != nil {
		t.Fatal(err)
	}

	user := model.User{
		Username:     "verified-user",
		Email:        "verified-user@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleUser,
		RoleID:       &customRole.ID,
	}
	if err := database.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	session := model.Session{
		UserID:    user.ID,
		SessionID: "verified-session",
	}
	if err := database.Create(&session).Error; err != nil {
		t.Fatal(err)
	}

	customRole.Enabled = false
	if err := UpdateRole(ctx, customRole); err != nil {
		t.Fatal(err)
	}

	if err := database.First(&user, user.ID).Error; err != nil {
		t.Fatal(err)
	}
	if user.Role != constant.RoleUser || user.RoleID == nil || *user.RoleID != memberRole.ID {
		t.Fatalf("disabled role user was not reset: %#v", user)
	}
	valid, err := IsSessionValidForUser(ctx, session.SessionID, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Fatal("disabled role user session must be revoked")
	}
}
