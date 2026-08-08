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
		PostID:   post.ID,
		AuthorID: user.ID,
		Content:  "content",
	}
	if err := CreateComment(context.Background(), &comment); err != nil {
		t.Fatal(err)
	}

	comments, total, err := ListCommentsByPostID(
		context.Background(),
		post.ID,
		0,
		20,
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
