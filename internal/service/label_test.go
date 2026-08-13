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

func TestLabelManagementService(t *testing.T) {
	t.Setenv(constant.EnvKeyDBDriver, "sqlite")
	t.Setenv(
		constant.EnvKeyDBSQLitePath,
		filepath.Join(t.TempDir(), "label-service.db"),
	)
	if err := repo.InitDatabase(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	_, err := CreateLabel(ctx, constant.RoleUser, dto.CreateLabelRequest{
		Name: "Go",
		Slug: "go",
	})
	requireServiceStatus(t, err, http.StatusForbidden)

	created, err := CreateLabel(ctx, constant.RoleEditor, dto.CreateLabelRequest{
		Name: " Go ",
		Slug: " Go Lang ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.Name != "Go" || created.Slug != "go-lang" {
		t.Fatalf("unexpected created label: %#v", created)
	}

	_, err = CreateLabel(ctx, constant.RoleEditor, dto.CreateLabelRequest{
		Name: "Go",
		Slug: "other",
	})
	requireServiceStatus(t, err, http.StatusConflict)

	name := "Golang"
	slug := "golang"
	updated, err := UpdateLabel(
		ctx,
		uint(created.ID),
		constant.RoleAdmin,
		dto.UpdateLabelRequest{
			Name: &name,
			Slug: &slug,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != name || updated.Slug != slug {
		t.Fatalf("unexpected updated label: %#v", updated)
	}

	if err := DeleteLabel(ctx, uint(created.ID), constant.RoleEditor); err != nil {
		t.Fatal(err)
	}
	if err := DeleteLabel(ctx, uint(created.ID), constant.RoleEditor); err == nil {
		t.Fatal("expected missing label deletion to fail")
	} else {
		requireServiceStatus(t, err, http.StatusNotFound)
	}
}

func TestDeleteLabelRemovesPostAssociation(t *testing.T) {
	t.Setenv(constant.EnvKeyDBDriver, "sqlite")
	t.Setenv(
		constant.EnvKeyDBSQLitePath,
		filepath.Join(t.TempDir(), "label-association.db"),
	)
	if err := repo.InitDatabase(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	label := model.Label{Name: "Go", Slug: "go"}
	if err := repo.CreateLabel(ctx, &label); err != nil {
		t.Fatal(err)
	}

	user := model.User{
		Username:     "label-editor",
		Email:        "label-editor@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleEditor,
	}
	if err := repo.CreateUser(ctx, &user); err != nil {
		t.Fatal(err)
	}

	post := model.Post{
		PostBase: model.PostBase{
			Title:        "Label post",
			Content:      "content",
			DraftContent: "content",
			Slug:         "label-post",
		},
		AuthorID: user.ID,
		Labels:   []model.Label{label},
	}
	if err := repo.CreatePost(ctx, &post); err != nil {
		t.Fatal(err)
	}

	if err := DeleteLabel(ctx, label.ID, constant.RoleAdmin); err != nil {
		t.Fatal(err)
	}

	var count int64
	if err := repo.GetDB().Table("post_labels").Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected no post label associations, got %d", count)
	}
}
