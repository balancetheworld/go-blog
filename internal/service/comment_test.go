package service

import (
	"context"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
)

func TestCommentTargetsRepliesAndCounters(t *testing.T) {
	t.Setenv(constant.EnvKeyDBDriver, "sqlite")
	t.Setenv(
		constant.EnvKeyDBSQLitePath,
		filepath.Join(t.TempDir(), "comment-service.db"),
	)
	if err := repo.InitDatabase(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	memberRole, err := repo.GetRoleByCode(ctx, constant.RoleCodeMember)
	if err != nil {
		t.Fatal(err)
	}
	user := model.User{
		Username:     "comment-service-user",
		Email:        "comment-service-user@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleUser,
		RoleID:       &memberRole.ID,
	}
	if err := repo.CreateUser(ctx, &user); err != nil {
		t.Fatal(err)
	}

	post := model.Post{
		PostBase: model.PostBase{
			Title:        "Comment Service Post",
			Content:      "content",
			DraftContent: "content",
			Slug:         "comment-service-post",
		},
		AuthorID: user.ID,
	}
	if err := repo.CreatePost(ctx, &post); err != nil {
		t.Fatal(err)
	}

	top, err := CreateComment(
		ctx,
		user.ID,
		constant.RoleUser,
		memberRole.ID,
		dto.CreateCommentRequest{
			PostID:  post.ID,
			Content: "top",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if top.TargetType != constant.TargetPost || top.TargetID != uint64(post.ID) || top.Depth != 0 {
		t.Fatalf("unexpected top comment: %#v", top)
	}

	reply, err := CreateComment(
		ctx,
		user.ID,
		constant.RoleUser,
		memberRole.ID,
		dto.CreateCommentRequest{
			TargetType: constant.TargetComment,
			TargetID:   uint(top.ID),
			Content:    "reply",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if reply.ParentID == nil || *reply.ParentID != top.ID || reply.Depth != 1 {
		t.Fatalf("unexpected reply: %#v", reply)
	}

	nestedReply, err := CreateComment(
		ctx,
		user.ID,
		constant.RoleUser,
		memberRole.ID,
		dto.CreateCommentRequest{
			TargetType: constant.TargetComment,
			TargetID:   uint(reply.ID),
			Content:    "nested reply",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if nestedReply.RootID == nil || *nestedReply.RootID != top.ID || nestedReply.Depth != 2 {
		t.Fatalf("unexpected nested reply: %#v", nestedReply)
	}

	_, err = CreateComment(
		ctx,
		user.ID,
		constant.RoleUser,
		memberRole.ID,
		dto.CreateCommentRequest{
			TargetType: constant.TargetComment,
			TargetID:   uint(nestedReply.ID),
			Content:    "too deep",
		},
	)
	requireServiceStatus(t, err, http.StatusBadRequest)

	storedTop, err := repo.GetCommentByID(ctx, uint(top.ID))
	if err != nil {
		t.Fatal(err)
	}
	if storedTop.ReplyCount != 1 {
		t.Fatalf("expected one direct reply, got %d", storedTop.ReplyCount)
	}

	countedPost, err := repo.GetPostByID(ctx, post.ID)
	if err != nil {
		t.Fatal(err)
	}
	if countedPost.CommentCount != 3 || countedPost.Heat != 15 {
		t.Fatalf(
			"unexpected post counters: comments=%d heat=%v",
			countedPost.CommentCount,
			countedPost.Heat,
		)
	}

	if err := DeleteComment(ctx, uint(top.ID), user.ID, constant.RoleUser); err != nil {
		t.Fatal(err)
	}
	countedPost, err = repo.GetPostByID(ctx, post.ID)
	if err != nil {
		t.Fatal(err)
	}
	if countedPost.CommentCount != 0 || countedPost.Heat != 0 {
		t.Fatalf(
			"unexpected counters after delete: comments=%d heat=%v",
			countedPost.CommentCount,
			countedPost.Heat,
		)
	}
}

func TestDiaryCommentCounter(t *testing.T) {
	t.Setenv(constant.EnvKeyDBDriver, "sqlite")
	t.Setenv(
		constant.EnvKeyDBSQLitePath,
		filepath.Join(t.TempDir(), "diary-comment-service.db"),
	)
	if err := repo.InitDatabase(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	memberRole, err := repo.GetRoleByCode(ctx, constant.RoleCodeMember)
	if err != nil {
		t.Fatal(err)
	}
	user := model.User{
		Username:     "diary-comment-user",
		Email:        "diary-comment-user@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleUser,
		RoleID:       &memberRole.ID,
	}
	if err := repo.CreateUser(ctx, &user); err != nil {
		t.Fatal(err)
	}
	publishedAt := time.Now()
	diary := model.Diary{
		Content:      "content",
		DraftContent: "content",
		AuthorID:     user.ID,
		Visibility:   constant.PostVisibilityPublic,
		PublishedAt:  &publishedAt,
	}
	if err := repo.CreateDiary(ctx, &diary); err != nil {
		t.Fatal(err)
	}

	comment, err := CreateComment(
		ctx,
		user.ID,
		constant.RoleUser,
		memberRole.ID,
		dto.CreateCommentRequest{
			TargetType: constant.TargetDiary,
			TargetID:   diary.ID,
			Content:    "diary comment",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	countedDiary, err := repo.GetDiaryByID(ctx, diary.ID)
	if err != nil {
		t.Fatal(err)
	}
	if countedDiary.CommentCount != 1 {
		t.Fatalf("expected one diary comment, got %d", countedDiary.CommentCount)
	}

	if err := DeleteComment(ctx, uint(comment.ID), user.ID, constant.RoleUser); err != nil {
		t.Fatal(err)
	}
	countedDiary, err = repo.GetDiaryByID(ctx, diary.ID)
	if err != nil {
		t.Fatal(err)
	}
	if countedDiary.CommentCount != 0 {
		t.Fatalf("expected no diary comments, got %d", countedDiary.CommentCount)
	}
}
