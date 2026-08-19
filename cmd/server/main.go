package main

import (
	"context"
	"log"

	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/internal/router"
	"github.com/zyj/my-blog/internal/service"
	"github.com/zyj/my-blog/internal/task"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/utils"
)

func main() {
	if err := utils.ValidateTokenConfig(); err != nil {
		panic(err)
	}
	if err := utils.ValidateCaptchaConfig(); err != nil {
		panic(err)
	}

	//初始化数据库
	if err := repo.InitDatabase(); err != nil {
		panic(err)
	}
	defer func() {
		if err := repo.CloseDatabase(); err != nil {
			log.Printf("close database failed: %v", err)
		}
	}()
	if err := service.EnsureRootUser(context.Background()); err != nil {
		panic(err)
	}
	if utils.GetAsBool(constant.EnvKeySeedDemoContent, false) {
		if err := service.EnsureDemoContent(context.Background()); err != nil {
			panic(err)
		}
	}
	if err := repo.InitRedis(context.Background()); err != nil {
		panic(err)
	}

	task.InitClient() // 全局client被初始化，拿到*asynq.Client，建立redis连接
	defer func() {
		if err := task.CloseClient(); err != nil {
			log.Printf("close asynq client failed: %v", err)
		}
	}()

	defer func() {
		if err := repo.CloseRedis(); err != nil {
			log.Printf("close redis failed: %v", err)
		}
	}()
	//初始化存储
	// 作用：系统统一的文件存储入口，用于处理图片、附件、上传文件等资源读写

	// 为什么必须程序启动时初始化：
	// 1. 存储方案存在多种实现：本地磁盘存储 / 对象存储OSS（阿里云/腾讯云）
	//    项目启动一次性根据环境变量加载对应驱动，业务代码不用关心底层是本地还是云存储
	// 2. 提前创建存储客户端、打开目录、校验权限，提前发现配置错误（路径不存在、无写入权限）
	//    避免等到用户上传文件时才报错，提前拦截启动阶段异常
	// 3. 全局只初始化一次，后续上传接口直接调用已经准备好的存储实例，不需要每次请求重复创建
	// 4. 统一封装，业务handler只调用统一接口 SaveFile / GetFile，底层存储实现可以无缝切换
	// tasks.InitStorageProvider()

	// cache.InitFileCache()

	if err := router.Run(); err != nil {
		panic(err)
	}
}
