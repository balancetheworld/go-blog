package service

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
	serviceerrs "github.com/zyj/my-blog/pkg/errs"
)

func requireServiceStatus(t *testing.T, err error, status int) {
	t.Helper()

	var serviceErr *serviceerrs.ServiceError
	if !errors.As(err, &serviceErr) {
		t.Fatalf("expected service error, got %v", err)
	}
	if serviceErr.HTTPStatus != status {
		t.Fatalf("expected status %d, got %d", status, serviceErr.HTTPStatus)
	}
}

func TestRoleManagementService(t *testing.T) {
	t.Setenv(constant.EnvKeyDBDriver, "sqlite")
	t.Setenv(
		constant.EnvKeyDBSQLitePath,
		filepath.Join(t.TempDir(), "role-service.db"),
	)
	if err := repo.InitDatabase(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	_, err := CreateRole(ctx, dto.CreateRoleRequest{
		Code: "INVALID CODE",
		Name: "非法身份",
	})
	requireServiceStatus(t, err, http.StatusBadRequest)

	created, err := CreateRole(ctx, dto.CreateRoleRequest{
		Code:          " Photographer ",
		Name:          " 摄影师 ",
		IsRequestable: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.Code != "photographer" || created.Name != "摄影师" || !created.Enabled {
		t.Fatalf("unexpected created role: %#v", created)
	}

	_, err = CreateRole(ctx, dto.CreateRoleRequest{
		Code: "photographer",
		Name: "重复身份",
	})
	requireServiceStatus(t, err, http.StatusConflict)

	roles, err := ListRoles(ctx, dto.ListRolesRequest{
		Keyword: "PHOTO",
	})
	if err != nil {
		t.Fatal(err)
	}
	if roles.Total != 1 || len(roles.Items) != 1 || roles.Items[0].ID != created.ID {
		t.Fatalf("unexpected roles: %#v", roles)
	}

	name := "摄影作者"
	enabled := false
	updated, err := UpdateRole(ctx, created.ID, dto.UpdateRoleRequest{
		Name:    &name,
		Enabled: &enabled,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != name || updated.Enabled {
		t.Fatalf("unexpected updated role: %#v", updated)
	}

	memberRole, err := repo.GetRoleByCode(ctx, constant.RoleCodeMember)
	if err != nil {
		t.Fatal(err)
	}
	_, err = UpdateRole(ctx, memberRole.ID, dto.UpdateRoleRequest{
		Name: &name,
	})
	requireServiceStatus(t, err, http.StatusForbidden)

	referenced, err := CreateRole(ctx, dto.CreateRoleRequest{
		Code: "writer",
		Name: "作者",
	})
	if err != nil {
		t.Fatal(err)
	}
	user := model.User{
		Username:     "role-service-user",
		Email:        "role-service-user@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleUser,
		RoleID:       &referenced.ID,
	}
	if err := repo.CreateUser(ctx, &user); err != nil {
		t.Fatal(err)
	}
	if err := DeleteRole(ctx, referenced.ID); err == nil {
		t.Fatal("expected referenced role deletion to fail")
	} else {
		requireServiceStatus(t, err, http.StatusConflict)
	}

	if err := DeleteRole(ctx, memberRole.ID); err == nil {
		t.Fatal("expected system role deletion to fail")
	} else {
		requireServiceStatus(t, err, http.StatusForbidden)
	}

	if err := DeleteRole(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
}
