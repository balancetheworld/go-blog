package task

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"

	"github.com/zyj/my-blog/internal/ai"
	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
)

type moderationProviderStub struct {
	result ai.ModerationResult
}

func (s moderationProviderStub) Moderate(
	context.Context,
	ai.ModerationInput,
) (ai.ModerationResult, error) {
	return s.result, nil
}

type moderationProviderErrorStub struct{}

func (moderationProviderErrorStub) Moderate(
	context.Context,
	ai.ModerationInput,
) (ai.ModerationResult, error) {
	return ai.ModerationResult{}, errors.New("provider unavailable")
}

type moderationProviderInvalidResultStub struct{}

func (moderationProviderInvalidResultStub) Moderate(
	context.Context,
	ai.ModerationInput,
) (ai.ModerationResult, error) {
	return ai.ModerationResult{
		Status:     constant.ModerationPending,
		Confidence: 0.8,
	}, nil
}

func TestCommentModerationHandlerProcessTask(t *testing.T) {
	t.Setenv(constant.EnvKeyDBDriver, "sqlite")
	t.Setenv(
		constant.EnvKeyDBSQLitePath,
		filepath.Join(t.TempDir(), "comment-moderation-handler.db"),
	)
	if err := repo.InitDatabase(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := repo.CloseDatabase(); err != nil {
			t.Fatal(err)
		}
	})

	database := repo.GetDB()
	user := model.User{
		Username:     "moderation-user",
		Email:        "moderation-user@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleUser,
	}
	if err := database.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	post := model.Post{
		PostBase: model.PostBase{
			Title:        "Moderation Post",
			Content:      "content",
			DraftContent: "content",
			Slug:         "moderation-post",
		},
		AuthorID: user.ID,
	}
	if err := database.Create(&post).Error; err != nil {
		t.Fatal(err)
	}
	comment := model.Comment{
		PostID:           &post.ID,
		TargetType:       constant.TargetPost,
		TargetID:         post.ID,
		AuthorID:         user.ID,
		Content:          "comment content",
		ModerationStatus: constant.ModerationPending,
	}
	if err := database.Create(&comment).Error; err != nil {
		t.Fatal(err)
	}
	aiTask := model.AITask{
		TaskType:   constant.AITaskCommentModeration,
		TargetType: constant.TargetComment,
		TargetID:   comment.ID,
		Status:     constant.AITaskQueued,
	}
	if err := database.Create(&aiTask).Error; err != nil {
		t.Fatal(err)
	}

	message, err := NewCommentModerationTask(aiTask.ID)
	if err != nil {
		t.Fatal(err)
	}
	handler := NewCommentModerationHandler(moderationProviderStub{
		result: ai.ModerationResult{
			Status:     constant.ModerationApproved,
			Categories: []string{"safe"},
			Confidence: 0.9,
			Reason:     "approved",
		},
	})
	if err := handler.ProcessTask(context.Background(), message); err != nil {
		t.Fatal(err)
	}

	updatedComment, err := repo.GetCommentByID(context.Background(), comment.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updatedComment.ModerationStatus != constant.ModerationApproved {
		t.Fatalf("unexpected moderation status: %s", updatedComment.ModerationStatus)
	}
	if updatedComment.ModerationReason != "approved" {
		t.Fatalf("unexpected moderation reason: %s", updatedComment.ModerationReason)
	}
	if updatedComment.ModerationConfidence != 0.9 {
		t.Fatalf("unexpected moderation confidence: %v", updatedComment.ModerationConfidence)
	}
	var categories []string
	if err := json.Unmarshal([]byte(updatedComment.ModerationCategories), &categories); err != nil {
		t.Fatal(err)
	}
	if len(categories) != 1 || categories[0] != "safe" {
		t.Fatalf("unexpected moderation categories: %#v", categories)
	}

	updatedTask, err := repo.GetAITaskByID(context.Background(), aiTask.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updatedTask.Status != constant.AITaskSucceeded || updatedTask.FinishedAt == nil {
		t.Fatalf("unexpected ai task: %#v", updatedTask)
	}
	var result ai.ModerationResult
	if err := json.Unmarshal([]byte(updatedTask.Result), &result); err != nil {
		t.Fatal(err)
	}
	if result.Status != constant.ModerationApproved || result.Confidence != 0.9 {
		t.Fatalf("unexpected ai task result: %#v", result)
	}
}

func TestCommentModerationHandlerProcessTaskProviderError(t *testing.T) {
	t.Setenv(constant.EnvKeyDBDriver, "sqlite")
	t.Setenv(
		constant.EnvKeyDBSQLitePath,
		filepath.Join(t.TempDir(), "comment-moderation-handler-error.db"),
	)
	if err := repo.InitDatabase(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := repo.CloseDatabase(); err != nil {
			t.Fatal(err)
		}
	})

	database := repo.GetDB()
	comment := model.Comment{
		TargetType:       constant.TargetGuestbook,
		AuthorID:         1,
		Content:          "comment content",
		ModerationStatus: constant.ModerationPending,
	}
	if err := database.Create(&comment).Error; err != nil {
		t.Fatal(err)
	}
	aiTask := model.AITask{
		TaskType:   constant.AITaskCommentModeration,
		TargetType: constant.TargetComment,
		TargetID:   comment.ID,
		Status:     constant.AITaskQueued,
	}
	if err := database.Create(&aiTask).Error; err != nil {
		t.Fatal(err)
	}

	message, err := NewCommentModerationTask(aiTask.ID)
	if err != nil {
		t.Fatal(err)
	}
	handler := NewCommentModerationHandler(moderationProviderErrorStub{})
	for range 3 {
		if err := handler.ProcessTask(context.Background(), message); err == nil {
			t.Fatal("expected provider error")
		}
	}

	updatedComment, err := repo.GetCommentByID(context.Background(), comment.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updatedComment.ModerationStatus != constant.ModerationPending {
		t.Fatalf("unexpected moderation status: %s", updatedComment.ModerationStatus)
	}

	updatedTask, err := repo.GetAITaskByID(context.Background(), aiTask.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updatedTask.Status != constant.AITaskDead ||
		updatedTask.Attempts != 3 ||
		updatedTask.FinishedAt == nil {
		t.Fatalf("unexpected ai task: %#v", updatedTask)
	}
	if updatedTask.Error != "moderate comment: provider unavailable" {
		t.Fatalf("unexpected ai task error: %s", updatedTask.Error)
	}
}

func TestNewCommentModerationMux(t *testing.T) {
	mux := NewCommentModerationMux(moderationProviderStub{})
	message, err := NewCommentModerationTask(1)
	if err != nil {
		t.Fatal(err)
	}

	handler, pattern := mux.Handler(message)
	if handler == nil || pattern != TypeCommentModeration {
		t.Fatalf("unexpected handler route: pattern=%s handler=%#v", pattern, handler)
	}
}

func TestNewCommentModerationMuxFromEnvMissingConfig(t *testing.T) {
	t.Setenv(constant.EnvKeyAIModerationBaseURL, "")
	t.Setenv(constant.EnvKeyAIModerationAPIKey, "")
	t.Setenv(constant.EnvKeyAIModerationModel, "")

	mux, err := NewCommentModerationMuxFromEnv()
	if err == nil {
		t.Fatal("expected configuration error")
	}
	if mux != nil {
		t.Fatal("expected nil mux")
	}
}

func TestNewCommentModerationMuxFromEnv(t *testing.T) {
	t.Setenv(constant.EnvKeyAIModerationBaseURL, "https://provider.test")
	t.Setenv(constant.EnvKeyAIModerationAPIKey, "test-key")
	t.Setenv(constant.EnvKeyAIModerationModel, "test-model")

	mux, err := NewCommentModerationMuxFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if mux == nil {
		t.Fatal("expected mux")
	}

	message, err := NewCommentModerationTask(1)
	if err != nil {
		t.Fatal(err)
	}
	handler, pattern := mux.Handler(message)
	if handler == nil || pattern != TypeCommentModeration {
		t.Fatalf("unexpected handler route: pattern=%s handler=%#v", pattern, handler)
	}
}

func TestCommentModerationHandlerProcessTaskInvalidResult(t *testing.T) {
	t.Setenv(constant.EnvKeyDBDriver, "sqlite")
	t.Setenv(
		constant.EnvKeyDBSQLitePath,
		filepath.Join(t.TempDir(), "comment-moderation-handler-invalid-result.db"),
	)
	if err := repo.InitDatabase(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := repo.CloseDatabase(); err != nil {
			t.Fatal(err)
		}
	})

	database := repo.GetDB()
	comment := model.Comment{
		TargetType:       constant.TargetGuestbook,
		AuthorID:         1,
		Content:          "comment content",
		ModerationStatus: constant.ModerationPending,
	}
	if err := database.Create(&comment).Error; err != nil {
		t.Fatal(err)
	}
	aiTask := model.AITask{
		TaskType:   constant.AITaskCommentModeration,
		TargetType: constant.TargetComment,
		TargetID:   comment.ID,
		Status:     constant.AITaskQueued,
	}
	if err := database.Create(&aiTask).Error; err != nil {
		t.Fatal(err)
	}

	message, err := NewCommentModerationTask(aiTask.ID)
	if err != nil {
		t.Fatal(err)
	}
	handler := NewCommentModerationHandler(moderationProviderInvalidResultStub{})
	for range 3 {
		if err := handler.ProcessTask(context.Background(), message); err == nil {
			t.Fatal("expected invalid moderation result error")
		}
	}

	updatedComment, err := repo.GetCommentByID(context.Background(), comment.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updatedComment.ModerationStatus != constant.ModerationPending {
		t.Fatalf("unexpected moderation status: %s", updatedComment.ModerationStatus)
	}

	updatedTask, err := repo.GetAITaskByID(context.Background(), aiTask.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updatedTask.Status != constant.AITaskDead ||
		updatedTask.Attempts != 3 ||
		updatedTask.FinishedAt == nil {
		t.Fatalf("unexpected ai task: %#v", updatedTask)
	}
	if updatedTask.Error != "validate moderation result: invalid moderation status" {
		t.Fatalf("unexpected ai task error: %s", updatedTask.Error)
	}
}
