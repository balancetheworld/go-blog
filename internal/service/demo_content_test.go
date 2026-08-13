package service

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
)

func TestEnsureDemoContentIsIdempotent(t *testing.T) {
	t.Setenv(constant.EnvKeyDBDriver, "sqlite")
	t.Setenv(constant.EnvKeyDBSQLitePath, filepath.Join(t.TempDir(), "demo-content.db"))
	t.Setenv(constant.EnvKeyRootAdminUsername, "demo-root")
	t.Setenv(constant.EnvKeyRootAdminEmail, "demo-root@example.com")
	t.Setenv(constant.EnvKeyRootAdminPassword, "demo-password")
	if err := repo.InitDatabase(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := EnsureRootUser(ctx); err != nil {
		t.Fatal(err)
	}
	if err := EnsureDemoContent(ctx); err != nil {
		t.Fatal(err)
	}
	if err := EnsureDemoContent(ctx); err != nil {
		t.Fatal(err)
	}

	categories, err := repo.ListCategories(ctx)
	if err != nil {
		t.Fatal(err)
	}
	posts, total, err := repo.ListPosts(ctx, repo.PostListFilter{Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(categories) != len(demoCategories) {
		t.Fatalf("expected %d categories, got %d", len(demoCategories), len(categories))
	}
	if total != int64(len(demoPosts)) || len(posts) != len(demoPosts) {
		t.Fatalf("expected %d posts, got total=%d items=%d", len(demoPosts), total, len(posts))
	}
}
