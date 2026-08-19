package repo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
)

func GetAITaskByID(ctx context.Context, id uint) (model.AITask, error) {
	if db == nil {
		return model.AITask{}, errors.New("database is not initialized")
	}

	var task model.AITask
	err := db.WithContext(ctx).First(&task, id).Error
	return task, err
}

// 认领 AI 任务
// 把一条待处理 AI 任务抢占过来，标记为我这个实例负责执行，防止多个 worker 重复处理同一个任务
func ClaimAITask(ctx context.Context, id uint) (model.AITask, bool, error) {
	if db == nil {
		return model.AITask{}, false, errors.New("database is not initialized")
	}
	result := db.WithContext(ctx).
		Model(&model.AITask{}).
		Where("id = ? AND status = ?", id, constant.AITaskQueued).
		Updates(map[string]any{ //map [string] any  零值会生效。比如你想把某个字段更新为 `""`、`0`，map 会正常更新到数据库。
			"status":     constant.AITaskProcessing,
			"attempts":   gorm.Expr("attempts + ?", 1),
			"started_at": time.Now(),
		})
	if result.Error != nil {
		return model.AITask{}, false, result.Error
	}
	if result.RowsAffected == 0 {
		return model.AITask{}, false, nil
	}

	task, err := GetAITaskByID(ctx, id)
	if err != nil {
		return model.AITask{}, false, err
	}
	return task, true, nil
}

// 将 AI 任务标记为完成，是任务生命周期里的收尾操作。
func CompleteAITask(ctx context.Context, id uint, result string) (bool, error) {
	if db == nil {
		return false, errors.New("database is not initialized")
	}
	update := db.WithContext(ctx).
		Model(&model.AITask{}).
		Where("id = ? AND status = ?", id, constant.AITaskProcessing).
		Updates(map[string]any{
			"status":      constant.AITaskSucceeded,
			"result":      result,
			"error":       "",
			"finished_at": time.Now(),
		})
	if update.Error != nil {
		return false, update.Error
	}
	return update.RowsAffected > 0, nil
}

// 把 AI 任务标记为失败，任务生命周期的失败分支处理函数。
func FailAITask(
	ctx context.Context,
	id uint,
	taskError string,
	maxAttempts uint,
) (bool, error) {
	if db == nil {
		return false, errors.New("database is not initialized")
	}
	taskError = strings.TrimSpace(taskError)
	if taskError == "" {
		taskError = "ai task failed"
	}
	updated := false
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var task model.AITask
		if err := tx.First(&task, id).Error; err != nil {
			return err
		}
		if task.Status != constant.AITaskProcessing {
			return nil
		}

		values := map[string]any{
			"status":      constant.AITaskQueued,
			"error":       taskError,
			"started_at":  nil,
			"finished_at": nil,
		}
		if maxAttempts == 0 || task.Attempts >= maxAttempts {
			values["status"] = constant.AITaskDead
			values["finished_at"] = time.Now()
		}

		result := tx.Model(&model.AITask{}).
			Where("id = ? AND status = ?", id, constant.AITaskProcessing).
			Updates(values)
		if result.Error != nil {
			return result.Error
		}
		updated = result.RowsAffected > 0
		return nil
	})
	if err != nil {
		return false, err
	}
	return updated, nil
}

func GetCommentModerationTask(
	ctx context.Context,
	commentID uint,
) (model.AITask, error) {
	if db == nil {
		return model.AITask{}, errors.New("database is not initialized")
	}
	var task model.AITask
	err := db.WithContext(ctx).
		Where(
			"task_type = ? AND target_type = ? AND target_id = ?",
			constant.AITaskCommentModeration,
			constant.TargetComment,
			commentID,
		).
		First(&task).
		Error
	return task, err
}
