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
