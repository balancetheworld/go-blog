package task

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const TypeCommentModeration = "ai:comment:moderation"

type CommentModerationPayload struct { //payload载荷：盒子里面真正要传递的业务内容
	AITaskID uint `json:"ai_task_id"`
}

//这是一个工厂函数：专门用来创建一条「评论审核」的 Asynq 队列任务对象 `*asynq.Task`。 
//给它一个数据库 AITask 的 ID，它帮你把数据打包，生成可以丢进 Redis 异步队列的任务，给后面 worker 去消费执行评论审核。
func NewCommentModerationTask(aiTaskID uint) (*asynq.Task, error) { //`Task`：asynq 包里面定义的结构体类型，代表队列里的一条任务消息
	//返回值 1：`*asynq.Task`，Asynq 队列任务指针，拿到之后就可以调用 `client.Enqueue(...)` 丢 Redis 队列。
	payload, err := json.Marshal(CommentModerationPayload{ //构造 payload 并且序列化为字节
		AITaskID: aiTaskID,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCommentModeration, payload), nil //`NewTask` 需要两个入参：任务类型名、载荷字节数组 payload。
}
