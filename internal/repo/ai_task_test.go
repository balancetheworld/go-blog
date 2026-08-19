package repo

import (
	"context"
	"testing"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAITaskLifecycle(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file:ai_task_repo_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(&model.AITask{}); err != nil {
		t.Fatal(err)
	}

	previousDB := db
	db = database
	t.Cleanup(func() {
		db = previousDB
	})

	ctx := context.Background()
	task := model.AITask{
		TaskType:   constant.AITaskCommentModeration,
		TargetType: constant.TargetComment,
		TargetID:   1,
		Status:     constant.AITaskQueued,
	}
	if err := database.Create(&task).Error; err != nil {
		t.Fatal(err)
	}

	claimed, ok, err := ClaimAITask(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || claimed.Status != constant.AITaskProcessing || claimed.Attempts != 1 || claimed.StartedAt == nil {
		t.Fatalf("unexpected claimed task: %#v", claimed)
	}

	_, ok, err = ClaimAITask(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected second claim to fail")
	}

	completed, err := CompleteAITask(ctx, task.ID, `{"decision":"approved"}`)
	if err != nil {
		t.Fatal(err)
	}
	if !completed {
		t.Fatal("expected task completion to succeed")
	}

	completedTask, err := GetAITaskByID(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if completedTask.Status != constant.AITaskSucceeded || completedTask.FinishedAt == nil {
		t.Fatalf("unexpected completed task: %#v", completedTask)
	}

	retryTask := model.AITask{
		TaskType:   constant.AITaskCommentModeration,
		TargetType: constant.TargetComment,
		TargetID:   2,
		Status:     constant.AITaskQueued,
	}
	if err := database.Create(&retryTask).Error; err != nil {
		t.Fatal(err)
	}

	if _, ok, err := ClaimAITask(ctx, retryTask.ID); err != nil || !ok {
		t.Fatalf("unexpected first claim: ok=%v err=%v", ok, err)
	}
	if updated, err := FailAITask(ctx, retryTask.ID, "timeout", 2); err != nil || !updated {
		t.Fatalf("unexpected first failure: updated=%v err=%v", updated, err)
	}

	queuedTask, err := GetAITaskByID(ctx, retryTask.ID)
	if err != nil {
		t.Fatal(err)
	}
	if queuedTask.Status != constant.AITaskQueued || queuedTask.Attempts != 1 {
		t.Fatalf("unexpected queued retry task: %#v", queuedTask)
	}

	if _, ok, err := ClaimAITask(ctx, retryTask.ID); err != nil || !ok {
		t.Fatalf("unexpected second claim: ok=%v err=%v", ok, err)
	}
	if updated, err := FailAITask(ctx, retryTask.ID, "timeout", 2); err != nil || !updated {
		t.Fatalf("unexpected second failure: updated=%v err=%v", updated, err)
	}

	deadTask, err := GetAITaskByID(ctx, retryTask.ID)
	if err != nil {
		t.Fatal(err)
	}
	if deadTask.Status != constant.AITaskDead || deadTask.Attempts != 2 || deadTask.FinishedAt == nil {
		t.Fatalf("unexpected dead task: %#v", deadTask)
	}
}
