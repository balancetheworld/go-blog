package repo

import (
	"context"
	"testing"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPostLabelsAndHeatCounters(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file:post_repo_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Label{},
		&model.Post{},
	); err != nil {
		t.Fatal(err)
	}

	previousDB := db
	db = database
	t.Cleanup(func() {
		db = previousDB
	})

	user := model.User{
		Username:     "editor",
		Email:        "editor@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleEditor,
	}
	if err := database.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	labels := []model.Label{
		{Name: "Go", Slug: "go"},
		{Name: "Hertz", Slug: "hertz"},
	}
	if err := database.Create(&labels).Error; err != nil {
		t.Fatal(err)
	}

	post := model.Post{
		PostBase: model.PostBase{
			Title:        "Post",
			Content:      "content",
			DraftContent: "content",
			Slug:         "post",
		},
		AuthorID: user.ID,
		Labels:   labels,
	}
	if err := CreatePost(context.Background(), &post); err != nil {
		t.Fatal(err)
	}

	created, err := GetPostByID(context.Background(), post.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(created.Labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(created.Labels))
	}

	created.Labels = labels[:1]
	if err := UpdatePost(context.Background(), &created); err != nil {
		t.Fatal(err)
	}

	updated, err := GetPostByID(context.Background(), post.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(updated.Labels) != 1 || updated.Labels[0].ID != labels[0].ID {
		t.Fatalf("unexpected labels: %#v", updated.Labels)
	}

	posts, total, err := ListPosts(context.Background(), PostListFilter{
		Limit:   10,
		LabelID: labels[0].ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(posts) != 1 || posts[0].ID != post.ID {
		t.Fatalf("unexpected filtered posts: total=%d posts=%#v", total, posts)
	}

	if err := IncrementPostViewCount(context.Background(), post.ID); err != nil {
		t.Fatal(err)
	}
	if err := UpdatePostCommentCount(context.Background(), post.ID, 2); err != nil {
		t.Fatal(err)
	}
	if err := UpdatePostCommentCount(context.Background(), post.ID, -3); err != nil {
		t.Fatal(err)
	}

	counted, err := GetPostByID(context.Background(), post.ID)
	if err != nil {
		t.Fatal(err)
	}
	if counted.ViewCount != 1 || counted.CommentCount != 0 || counted.Heat != 1 {
		t.Fatalf(
			"unexpected counters: views=%d comments=%d heat=%v",
			counted.ViewCount,
			counted.CommentCount,
			counted.Heat,
		)
	}

	rowsAffected, err := DeletePost(context.Background(), post.ID)
	if err != nil {
		t.Fatal(err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 deleted post, got %d", rowsAffected)
	}

	var associationCount int64
	if err := database.Table("post_labels").Count(&associationCount).Error; err != nil {
		t.Fatal(err)
	}
	if associationCount != 0 {
		t.Fatalf("expected no post label associations, got %d", associationCount)
	}
}

func TestMigrateBackfillsPostFields(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file:post_migrate_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(
		&model.User{},
		&model.Label{},
		&model.Post{},
	); err != nil {
		t.Fatal(err)
	}

	previousDB := db
	db = database
	t.Cleanup(func() {
		db = previousDB
	})

	user := model.User{
		Username:     "admin",
		Email:        "admin@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleAdmin,
	}
	if err := database.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	post := model.Post{
		PostBase: model.PostBase{
			Title:        "Published",
			Content:      "content",
			DraftContent: "content",
			Slug:         "published",
			LikeCount:    2,
			CommentCount: 3,
			ViewCount:    4,
		},
		AuthorID: user.ID,
	}
	if err := database.Create(&post).Error; err != nil {
		t.Fatal(err)
	}

	if err := migrate(); err != nil {
		t.Fatal(err)
	}

	var migrated model.Post
	if err := database.First(&migrated, post.ID).Error; err != nil {
		t.Fatal(err)
	}
	if migrated.PublishedAt == nil {
		t.Fatal("expected published_at to be backfilled")
	}
	if migrated.Heat != 25 {
		t.Fatalf("expected heat 25, got %v", migrated.Heat)
	}
}
