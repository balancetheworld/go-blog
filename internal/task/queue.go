package task

import (
	"github.com/hibiken/asynq"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/utils"
)

// NewRedisClientOpt 构造Asynq连接Redis的配置对象
// 返回 asynq.RedisClientOpt：Asynq库定义的结构体，存放Redis连接参数（地址、密码、db库号）
// 不从代码写死硬编码，全部从环境变量读取，同时设置默认兜底值
func NewRedisClientOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		// Addr：Redis服务地址，格式 host:port
		// utils.Get(环境变量key, 默认值)：读取环境变量EnvKeyRedisAddr，没有读到就使用默认 127.0.0.1:6379
		Addr: utils.Get(constant.EnvKeyRedisAddr, "127.0.0.1:6379"),
		// Password：Redis密码，没有设置密码环境变量时为空字符串
		Password: utils.Get(constant.EnvKeyRedisPassword),
		// DB：Redis选择第几个数据库，数字；读环境变量，获取失败默认使用db 0
		DB: utils.GetAsInt(constant.EnvKeyRedisDB, 0),
	}
}
// Client（发任务）、Server（worker 消费任务）两边都要连接 Redis。
// 把连接配置抽出来，只写一处，两边复用。以后改 Redis 地址，只改环境变量，不用改两处代码。



// NewClient 创建 Asynq 的客户端实例
// *asynq.Client：生产者客户端，用于【往队列投递任务】Enqueue，业务接口层使用
// 内部复用上面NewRedisClientOpt拿到redis连接配置
func NewClient() *asynq.Client {
	// asynq.NewClient(redis配置) 根据redis配置生成客户端对象
	return asynq.NewClient(NewRedisClientOpt())
}
// 给 HTTP 接口用的，生产者，用来 Enqueue 投递任务到 Redis。
// 没有它，就没办法把任务丢进队列。


// NewServer 创建 Asynq 的服务端（Worker工作节点）实例
// *asynq.Server：worker服务，负责监听Redis队列、取出任务、执行业务handler，后台常驻运行
func NewServer() *asynq.Server {
	return asynq.NewServer(
		// 传入Redis连接配置，worker也要连接同一个Redis，才能消费队列任务
		NewRedisClientOpt(),
		// asynq.Config：worker的运行配置结构体
		asynq.Config{
			// Concurrency：并发数，同时最多可以并行处理2个任务
			// 即同时跑2个handler；不要设置过大，受CPU、Redis、LLM接口限流制约
			Concurrency: 2,
			// Queues：队列权重配置 map[队列名称]权重
			// "default"队列，权重1；权重越大，被worker取到任务的优先级越高
			Queues: map[string]int{
				"default": 1,
			},
		},
	)
}
// 启动 Worker，后台消费 Redis 里面的任务，真正执行 Agent、调用大模型。
// 没有 Server，就算你把任务丢进 Redis，也没有人去干活。
