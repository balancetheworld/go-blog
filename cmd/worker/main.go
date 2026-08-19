package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/internal/task"
)

func main() {
	if err := repo.InitDatabase(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := repo.CloseDatabase(); err != nil {
			log.Printf("close database failed: %v", err)
		}
	}()
	mux, err := task.NewCommentModerationMuxFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	server := task.NewServer()

	errCh := make(chan error, 1)
	// 把阻塞的 server.Run(mux) 放到新goroutine后台运行
	go func() {
		// Run是阻塞的，它的返回值（错误）发送到errCh通道
		errCh <- server.Run(mux)
	}()

	//创建信号通道
	// 创建信号通道，接收操作系统信号：Ctrl+C(SIGINT)、容器终止信号(SIGTERM)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM) //原本信号来了操作系统直接杀进程；调用 Notify 之后，操作系统不再粗暴杀死，把信号投递进你给的 channel。
	defer signal.Stop(signalCh)

	// select 等待：要么worker出错，要么收到关闭信号
	select {
	case err := <-errCh:
		// worker内部发生致命错误，Run返回了错误，直接退出程序
		if err != nil {
			log.Fatal(err)
		}
	case <-signalCh:
		// ✅用户按Ctrl‑C / 容器要关闭，走到这里
		server.Shutdown() // asynq优雅关闭：停止拉新任务，等待正在跑的任务执行完毕
		// 等待Run返回，Shutdown完成之后server.Run才会返回
		if err := <-errCh; err != nil {
			log.Fatal(err)
		}
	}

}
