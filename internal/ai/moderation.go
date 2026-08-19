package ai

import (
	"context"
	"errors"
	"math"

	"github.com/zyj/my-blog/pkg/constant"
)

type ModerationInput struct {
	CommentID  uint
	Content    string
	TargetType constant.TargetType
	TargetID   uint
}

type ModerationResult struct {
	Status     constant.ModerationStatus
	Categories []string
	Confidence float64
	Reason     string
}

type ModerationProvider interface {
	Moderate(context.Context, ModerationInput) (ModerationResult, error)
}

// 对大模型返回的 ModareationResult 做合法性校验的函数
func ValidateModerationResult(result ModerationResult) error {
	switch result.Status {
	case constant.ModerationApproved, constant.ModerationRejected, constant.ModerationManualReview:
	default:
		return errors.New("invalid moderation status")
	}
	if math.IsNaN(result.Confidence) || result.Confidence < 0 || result.Confidence > 1 {
		return errors.New("invalid moderation confidence")
	}
	if len(result.Categories) > 10 {
		return errors.New("too many moderation categories")
	}
	if len(result.Reason) > 500 {
		return errors.New("moderation reason is too long")
	}

	return nil
}
