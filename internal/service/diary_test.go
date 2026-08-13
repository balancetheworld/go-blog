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

func TestDiaryManagementAndVisibility(t *testing.T) {
	t.Setenv(constant.EnvKeyDBDriver, "sqlite")
	t.Setenv(
		constant.EnvKeyDBSQLitePath,
		filepath.Join(t.TempDir(), "diary-service.db"),
	)
	if err := repo.InitDatabase(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	memberRole, err := repo.GetRoleByCode(ctx, constant.RoleCodeMember)
	if err != nil {
		t.Fatal(err)
	}
	customRole := model.Role{
		Code:    "verified",
		Name:    "已认证",
		Enabled: true,
	}
	if err := repo.CreateRole(ctx, &customRole); err != nil {
		t.Fatal(err)
	}

	editor := model.User{
		Username:     "diary-editor",
		Email:        "diary-editor@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleEditor,
		RoleID:       &memberRole.ID,
	}
	if err := repo.CreateUser(ctx, &editor); err != nil {
		t.Fatal(err)
	}

	publicDiary, err := CreateDiary(
		ctx,
		editor.ID,
		constant.RoleEditor,
		dto.CreateDiaryRequest{
			DraftContent: "公开日记",
			Publish:      true,
			Visibility:   constant.PostVisibilityPublic,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	roleDiary, err := CreateDiary(
		ctx,
		editor.ID,
		constant.RoleEditor,
		dto.CreateDiaryRequest{
			DraftContent:   "身份日记",
			Publish:        true,
			Visibility:     constant.PostVisibilityRoles,
			VisibleRoleIDs: []uint{customRole.ID},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = CreateDiary(
		ctx,
		editor.ID,
		constant.RoleEditor,
		dto.CreateDiaryRequest{
			DraftContent: "草稿",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	guestList, err := ListDiaries(
		ctx,
		dto.ListDiariesRequest{},
		0,
		constant.RoleGuest,
		0,
	)
	if err != nil {
		t.Fatal(err)
	}
	if guestList.Total != 1 || guestList.Items[0].ID != publicDiary.ID {
		t.Fatalf("unexpected guest list: %#v", guestList)
	}

	roleList, err := ListDiaries(
		ctx,
		dto.ListDiariesRequest{},
		999,
		constant.RoleUser,
		customRole.ID,
	)
	if err != nil {
		t.Fatal(err)
	}
	if roleList.Total != 2 {
		t.Fatalf("unexpected role list: %#v", roleList)
	}

	_, err = GetDiary(
		ctx,
		uint(roleDiary.ID),
		0,
		constant.RoleGuest,
		0,
	)
	requireServiceStatus(t, err, http.StatusNotFound)

	adminList, err := ListDiaries(
		ctx,
		dto.ListDiariesRequest{Status: "all"},
		1,
		constant.RoleAdmin,
		0,
	)
	if err != nil {
		t.Fatal(err)
	}
	if adminList.Total != 3 {
		t.Fatalf("unexpected admin list: %#v", adminList)
	}

	otherEditor := model.User{
		Username:     "other-diary-editor",
		Email:        "other-diary-editor@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleEditor,
		RoleID:       &memberRole.ID,
	}
	if err := repo.CreateUser(ctx, &otherEditor); err != nil {
		t.Fatal(err)
	}
	content := "修改"
	_, err = UpdateDiary(
		ctx,
		uint(publicDiary.ID),
		otherEditor.ID,
		constant.RoleEditor,
		dto.UpdateDiaryRequest{DraftContent: &content},
	)
	requireServiceStatus(t, err, http.StatusForbidden)

	if err := DeleteDiary(
		ctx,
		uint(roleDiary.ID),
		editor.ID,
		constant.RoleEditor,
	); err != nil {
		t.Fatal(err)
	}

	var associationCount int64
	if err := repo.GetDB().Table("diary_visible_roles").Count(&associationCount).Error; err != nil {
		t.Fatal(err)
	}
	if associationCount != 0 {
		t.Fatalf("expected no diary role associations, got %d", associationCount)
	}
}

func TestDiaryFolderSlugAndVisibility(t *testing.T) {
	t.Setenv(constant.EnvKeyDBDriver, "sqlite")
	t.Setenv(
		constant.EnvKeyDBSQLitePath,
		filepath.Join(t.TempDir(), "diary-folder-service.db"),
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
		Username:     "folder-editor",
		Email:        "folder-editor@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleEditor,
		RoleID:       &memberRole.ID,
	}
	if err := repo.CreateUser(ctx, &editor); err != nil {
		t.Fatal(err)
	}

	folder, err := CreateDiaryFolder(
		ctx,
		constant.RoleEditor,
		dto.CreateDiaryFolderRequest{
			Name:       "私密文件夹",
			Visibility: constant.PostVisibilityPrivate,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	folderID := uint(folder.ID)
	diary, err := CreateDiary(
		ctx,
		editor.ID,
		constant.RoleEditor,
		dto.CreateDiaryRequest{
			Title:        "同名日记",
			FolderID:     &folderID,
			DraftContent: "日记内容",
			Publish:      true,
			Visibility:   constant.PostVisibilityPublic,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	secondDiary, err := CreateDiary(
		ctx,
		editor.ID,
		constant.RoleEditor,
		dto.CreateDiaryRequest{
			Title:        "同名日记",
			DraftContent: "第二篇内容",
			Publish:      true,
			Visibility:   constant.PostVisibilityPublic,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if diary.Slug == secondDiary.Slug {
		t.Fatalf("expected unique slugs, got %q", diary.Slug)
	}

	guestList, err := ListDiaries(ctx, dto.ListDiariesRequest{}, 0, constant.RoleGuest, 0)
	if err != nil {
		t.Fatal(err)
	}
	if guestList.Total != 1 || guestList.Items[0].ID != secondDiary.ID {
		t.Fatalf("unexpected guest list: %#v", guestList)
	}
	_, err = GetDiaryByIdentifier(ctx, diary.Slug, 0, constant.RoleGuest, 0)
	requireServiceStatus(t, err, http.StatusNotFound)

	if err := DeleteDiaryFolder(ctx, uint(folder.ID), constant.RoleEditor); err != nil {
		t.Fatal(err)
	}
	storedDiary, err := repo.GetDiaryByID(ctx, uint(diary.ID))
	if err != nil {
		t.Fatal(err)
	}
	if storedDiary.FolderID != nil {
		t.Fatalf("expected diary folder to be cleared, got %d", *storedDiary.FolderID)
	}
}
