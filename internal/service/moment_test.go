package service

import (
	"context"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
)

func TestMomentManagement(t *testing.T) {
	t.Setenv(constant.EnvKeyDBDriver, "sqlite")
	t.Setenv(
		constant.EnvKeyDBSQLitePath,
		filepath.Join(t.TempDir(), "moment-service.db"),
	)
	if err := repo.InitDatabase(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	memberRole, err := repo.GetRoleByCode(ctx, constant.RoleCodeMember)
	if err != nil {
		t.Fatal(err)
	}
	editor := model.User{
		Username:     "moment-editor",
		Email:        "moment-editor@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleEditor,
		RoleID:       &memberRole.ID,
	}
	if err := repo.CreateUser(ctx, &editor); err != nil {
		t.Fatal(err)
	}
	otherEditor := model.User{
		Username:     "other-moment-editor",
		Email:        "other-moment-editor@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleEditor,
		RoleID:       &memberRole.ID,
	}
	if err := repo.CreateUser(ctx, &otherEditor); err != nil {
		t.Fatal(err)
	}

	_, err = CreateMoment(ctx, editor.ID, constant.RoleUser, dto.CreateMomentRequest{Content: "无权限"})
	requireServiceStatus(t, err, http.StatusForbidden)
	_, err = CreateMoment(ctx, editor.ID, constant.RoleEditor, dto.CreateMomentRequest{Content: "  "})
	requireServiceStatus(t, err, http.StatusBadRequest)

	first, err := CreateMoment(ctx, editor.ID, constant.RoleEditor, dto.CreateMomentRequest{Content: "第一条"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := CreateMoment(ctx, editor.ID, constant.RoleEditor, dto.CreateMomentRequest{Content: " 第二条 "})
	if err != nil {
		t.Fatal(err)
	}
	if second.Content != "第二条" {
		t.Fatalf("unexpected content: %q", second.Content)
	}

	list, err := ListMoments(ctx, dto.ListMomentsRequest{Page: 1, PageSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if list.Total != 2 || len(list.Items) != 1 || list.Items[0].ID != second.ID {
		t.Fatalf("unexpected moment list: %#v", list)
	}

	err = DeleteMoment(ctx, uint(first.ID), otherEditor.ID, constant.RoleEditor)
	requireServiceStatus(t, err, http.StatusForbidden)
	if err := DeleteMoment(ctx, uint(first.ID), otherEditor.ID, constant.RoleAdmin); err != nil {
		t.Fatal(err)
	}
	if err := DeleteMoment(ctx, uint(second.ID), editor.ID, constant.RoleEditor); err != nil {
		t.Fatal(err)
	}

	list, err = ListMoments(ctx, dto.ListMomentsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if list.Total != 0 || len(list.Items) != 0 {
		t.Fatalf("expected no moments, got %#v", list)
	}
}
