package task

import (
	"context"
	"errors"

	"github.com/hibiken/asynq"
)

// 包级全局变量，保存asynq客户端单例实例，整个程序复用这一个client对象
var client *asynq.Client

// InitClient 初始化全局的asynq客户端，懒加载模式
func InitClient() {
	// 如果client还没有初始化（nil），才去创建；避免重复多次NewClient
	if client == nil {
		// NewClient()就是上面写的，内部调用NewRedisClientOpt组装redis配置，生成*asynq.Client
		client = NewClient()
	}
}

// CloseClient 关闭asynq客户端，释放redis连接，程序退出的时候调用
func CloseClient() error {
	// 没有初始化，直接返回nil，不用处理
	if client == nil {
		return nil
	}

	// 局部变量current拷贝引用，防止并发情况下client被修改
	current := client
	// 把全局client置nil，标记已经关闭
	client = nil
	// 调用asynq库Close，关闭底层redis连接
	return current.Close()
}

// EnqueueCommentModeration 【对外业务方法】：把评论审核任务投递进入asynq redis队列
// ctx：上下文，传递超时、取消信号；aiTaskID：MySQL AITask表的任务主键ID
func EnqueueCommentModeration(
	ctx context.Context,
	aiTaskID uint,
) error {
	// 保护：如果client没有调用InitClient初始化，直接返回错误，不能投递任务
	if client == nil {
		return errors.New("asynq client is not initialized")
	}

	// 调用工厂函数 NewCommentModerationTask，构造 *asynq.Task
	// 内部：构建CommentModerationPayload结构体 → json.Marshal得到payload字节 → asynq.NewTask(任务类型,payload)
	message, err := NewCommentModerationTask(aiTaskID)
	if err != nil {
		return err
	}

	// client.EnqueueContext：带context的入队方法，真正写入Redis队列
	// 参数1 ctx：上下文，如果ctx取消，入队操作可以中断
	// 参数2 message：*asynq.Task 准备投递的任务对象
	// asynq.Queue("default") 指定投递到default队列
	// asynq.MaxRetry(3) asynq层面最大重试3次，worker执行失败会自动重试
	_, err = client.EnqueueContext(
		ctx,
		message,
		asynq.Queue("default"),
		asynq.MaxRetry(3),
	)
	// 返回错误；nil代表成功写入Redis队列
	return err
}
