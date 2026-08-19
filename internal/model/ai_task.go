package model

import (
	"time"

	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
)


//   TaskType       任务类型
//   TargetType     目标类型，例如 comment
//   TargetID       评论 ID
//   Status         任务状态
//   Attempts       重试次数
//   Provider       AI 供应商
//   Model          模型名称
//   PromptVersion  提示词版本
//   InputHash      输入内容哈希
//   Result         AI 返回的 JSON
//   Error          失败原因
//   StartedAt      开始时间
//   FinishedAt     完成时间
type AITask struct {
	gorm.Model
	TaskType      constant.AITaskType   `gorm:"size:32;not null;index"`
	TargetType    constant.TargetType   `gorm:"size:16;not null;index"`
	TargetID      uint                  `gorm:"not null;index"`
	Status        constant.AITaskStatus `gorm:"size:20;not null;default:queued;index"`
	Attempts      uint                  `gorm:"not null;default:0"`
	Provider      string                `gorm:"size:64"`
	ModelName         string                `gorm:"size:128"`
	PromptVersion string                `gorm:"size:64"`
	InputHash     string                `gorm:"size:64;index"`
	Result        string                `gorm:"type:text"`
	Error         string                `gorm:"type:text"`
	StartedAt     *time.Time
	FinishedAt    *time.Time
}
