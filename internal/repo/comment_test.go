package repo

import (
	"context"
	"testing"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCommentCreateListAndDelete(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file:comment_repo_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Label{},
		&model.Post{},
		&model.Comment{},
		&model.AITask{},
	); err != nil {
		t.Fatal(err)
	}

	previousDB := db
	db = database
	t.Cleanup(func() {
		db = previousDB
	})

	user := model.User{
		Username:     "commenter",
		Email:        "commenter@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleUser,
	}
	if err := database.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	post := model.Post{
		PostBase: model.PostBase{
			Title:        "Commented Post",
			Content:      "content",
			DraftContent: "content",
			Slug:         "commented-post",
		},
		AuthorID: user.ID,
	}
	if err := CreatePost(context.Background(), &post); err != nil {
		t.Fatal(err)
	}

	comment := model.Comment{
		PostID:   &post.ID,
		AuthorID: user.ID,
		Content:  "content",
	}
	if err := CreateComment(context.Background(), &comment); err != nil {
		t.Fatal(err)
	}
	var task model.AITask
	if err := database.
		Where(
			"task_type = ? AND target_type = ? AND target_id = ?",
			constant.AITaskCommentModeration,
			constant.TargetComment,
			comment.ID,
		).
		First(&task).
		Error; err != nil {
		t.Fatal(err)
	}
	if task.Status != constant.AITaskQueued {
		t.Fatalf("unexpected task status: %s", task.Status)
	}
	if err := UpdateCommentModeration(
		context.Background(),
		comment.ID,
		constant.ModerationApproved,
		"",
	); err != nil {
		t.Fatal(err)
	}

	comments, total, err := ListComments(
		context.Background(),
		CommentListFilter{
			TargetType: constant.TargetPost,
			TargetID:   post.ID,
			Offset:     0,
			Limit:      20,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(comments) != 1 {
		t.Fatalf("unexpected comments: total=%d comments=%#v", total, comments)
	}
	if comments[0].Author.ID != user.ID {
		t.Fatalf("expected author %d, got %d", user.ID, comments[0].Author.ID)
	}

	createdPost, err := GetPostByID(context.Background(), post.ID)
	if err != nil {
		t.Fatal(err)
	}
	if createdPost.CommentCount != 1 || createdPost.Heat != 5 {
		t.Fatalf(
			"unexpected counters after create: comments=%d heat=%v",
			createdPost.CommentCount,
			createdPost.Heat,
		)
	}

	rowsAffected, err := DeleteComment(context.Background(), comment)
	if err != nil {
		t.Fatal(err)
	}
	if rowsAffected != 1 {
		t.Fatalf("expected 1 deleted comment, got %d", rowsAffected)
	}

	deletedPost, err := GetPostByID(context.Background(), post.ID)
	if err != nil {
		t.Fatal(err)
	}
	if deletedPost.CommentCount != 0 || deletedPost.Heat != 0 {
		t.Fatalf(
			"unexpected counters after delete: comments=%d heat=%v",
			deletedPost.CommentCount,
			deletedPost.Heat,
		)
	}
}

func TestCommentMigrationBackfillsLegacyPostTarget(t *testing.T) {
	type legacyComment struct {
		gorm.Model
		PostID   uint   `gorm:"not null;index"`
		AuthorID uint   `gorm:"not null;index"`
		Content  string `gorm:"type:text;not null"`
	}

	database, err := gorm.Open(
		sqlite.Open("file:comment_migration_test?mode=memory&cache=shared"),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(&model.Post{}); err != nil {
		t.Fatal(err)
	}
	if err := database.Table("comments").AutoMigrate(&legacyComment{}); err != nil {
		t.Fatal(err)
	}
	page := model.Post{
		PostBase: model.PostBase{
			Title: "Legacy Page",
			Type:  "page",
			Slug:  "legacy-page",
		},
		AuthorID: 8,
	}
	if err := database.Create(&page).Error; err != nil {
		t.Fatal(err)
	}
	legacy := legacyComment{
		PostID:   page.ID,
		AuthorID: 8,
		Content:  "legacy comment",
	}
	if err := database.Table("comments").Create(&legacy).Error; err != nil {
		t.Fatal(err)
	}

	if err := database.AutoMigrate(&model.Comment{}); err != nil {
		t.Fatal(err)
	}
	if err := database.Model(&model.Comment{}).
		Where("target_id = ? AND post_id > ?", 0, 0).
		Updates(map[string]any{
			"target_type": constant.TargetPost,
			"target_id":   gorm.Expr("post_id"),
		}).
		Error; err != nil {
		t.Fatal(err)
	}
	if err := database.Model(&model.Comment{}).
		Where(
			"target_type = ? AND target_id IN (?)",
			constant.TargetPost,
			database.Model(&model.Post{}).Select("id").Where("type = ?", "page"),
		).
		UpdateColumn("target_type", constant.TargetPage).
		Error; err != nil {
		t.Fatal(err)
	}

	var migrated model.Comment
	if err := database.First(&migrated, legacy.ID).Error; err != nil {
		t.Fatal(err)
	}
	if migrated.PostID == nil || *migrated.PostID != legacy.PostID {
		t.Fatalf("unexpected legacy post id: %#v", migrated.PostID)
	}
	if migrated.TargetType != constant.TargetPage || migrated.TargetID != legacy.PostID {
		t.Fatalf("unexpected migrated target: %#v", migrated)
	}
}
