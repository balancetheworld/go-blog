package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/zyj/my-blog/internal/ai"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
)

type CommentModerationHandler struct {
	Provider    ai.ModerationProvider
	MaxAttempts uint
}

// `var _`
// 下划线 `_` 代表丢弃变量，**不会创建实际变量，不会分配内存**。
// 只是告诉编译器：做一次类型检查，这个变量我不用。
var _ asynq.Handler = (*CommentModerationHandler)(nil)

func NewCommentModerationHandler(provider ai.ModerationProvider) *CommentModerationHandler {
	return &CommentModerationHandler{
		Provider:    provider,
		MaxAttempts: 3,
	}
}

func (h *CommentModerationHandler) ProcessTask(
	ctx context.Context,
	task *asynq.Task,
) error {
	if h == nil {
		return errors.New("comment moderation handler is nil")
	}
	if task == nil {
		return errors.New("task is nil")
	}
	var payload CommentModerationPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("decode comment moderation payload: %w", err)
	}
	if payload.AITaskID == 0 {
		return errors.New("invalid ai task id")
	}

	aiTask, claimed, err := repo.ClaimAITask(ctx, payload.AITaskID)
	if err != nil {
		return fmt.Errorf("claim ai task: %w", err)
	}
	if !claimed {
		return nil
	}

	fail := func(taskErr error) error {
		_, err := repo.FailAITask(
			ctx,
			aiTask.ID,
			taskErr.Error(),
			h.MaxAttempts,
		)
		if err != nil {
			return fmt.Errorf("fail ai task: %v: %w", taskErr, err)
		}
		return taskErr
	}

	if aiTask.TaskType != constant.AITaskCommentModeration ||
		aiTask.TargetType != constant.TargetComment {
		return fail(errors.New("invalid comment moderation task target"))
	}

	comment, err := repo.GetCommentByID(ctx, aiTask.TargetID)
	if err != nil {
		return fail(fmt.Errorf("get comment: %w", err))
	}

	if h.Provider == nil {
		return fail(errors.New("moderation provider is not configured"))
	}

	result, err := h.Provider.Moderate(ctx, ai.ModerationInput{
		CommentID:  comment.ID,
		Content:    comment.Content,
		TargetType: comment.TargetType,
		TargetID:   comment.TargetID,
	})
	if err != nil {
		return fail(fmt.Errorf("moderate comment: %w", err))
	}
	if err := ai.ValidateModerationResult(result); err != nil {
		return fail(fmt.Errorf("validate moderation result: %w", err))
	}

	categories, err := json.Marshal(result.Categories)
	if err != nil {
		return fail(fmt.Errorf("encode moderation categories: %w", err))
	}
	if err := repo.UpdateCommentModerationResult(
		ctx,
		comment.ID,
		result.Status,
		result.Reason,
		string(categories),
		result.Confidence,
	); err != nil {
		return fail(fmt.Errorf("update comment moderation: %w", err))
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fail(fmt.Errorf("encode moderation result: %w", err))
	}

	completed, err := repo.CompleteAITask(
		ctx,
		aiTask.ID,
		string(resultJSON),
	)
	if err != nil {
		return fail(fmt.Errorf("complete ai task: %w", err))
	}
	if !completed {
		return fail(errors.New("ai task was not completed"))
	}

	return nil
}

//  `asynq.ServeMux`
//  和 Go 标准库 http 的`http.ServeMux`思想一模一样：**路由分发器**。
//  HTTP mux：收到 http 请求路径，匹配路径，执行对应的 http handler。
//  Asynq mux：worker 从 Redis 拿到任务，读取任务的类型，匹配类型，执行对应的`ProcessTask`。
// `mux.Handle(任务类型, handler实例)`
// 含义：当收到这个类型的任务，就交给这个 handler 去处理（调用它的 ProcessTask）。


// 2. NewCommentModerationMux 做的事
// 输入：传入 `ModerationProvider`（比如 Gemini 的实现）
// 1. 创建空的 mux 路由器
// 2. 调用`NewCommentModerationHandler(provider)`，把 provider 塞给 handler
// 3. `mux.Handle`：把任务类型常量和 handler 绑定在一起
// 4. 返回配置完毕的 mux
//  此时只是内存里面的对象，Redis 没有连接，worker 没有跑。
func NewCommentModerationMux(
	provider ai.ModerationProvider,
) *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.Handle(
		TypeCommentModeration,
		NewCommentModerationHandler(provider),
	)
	return mux
}

// NewCommentModerationMuxFromEnv 从环境变量读取配置，构造出asynq mux 返回配置好的异步任务路由器（asynq.ServeMux）
// 把「组装路由」和「启动 Worker」两件事拆开。
// 函数只负责：把任务类型和 handler 绑定好，造出路由器。
// 它不负责启动 worker，不连接 Redis，不阻塞。把造好的路由器交出去，由上层 main 决定什么时候跑。
func NewCommentModerationMuxFromEnv() (*asynq.ServeMux, error) {
	// 从环境变量(BaseURL、APIKey、Model等)构建OpenAIProvider实例
	provider := ai.NewOpenAIProviderFromEnv()
	// 校验provider配置是否完整（BaseURL、APIKey、Model不为空）
	if err := provider.ValidateConfig(); err != nil {
		return nil, err
	}
	// 把provider传入，调用之前写好的mux构造函数，返回组装完成的asynq路由mux
	return NewCommentModerationMux(provider), nil
}