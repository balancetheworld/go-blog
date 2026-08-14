package router

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/zyj/my-blog/internal/router/apiv1"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/resps"
	"github.com/zyj/my-blog/pkg/utils"
)

func newServer() (*server.Hertz, error) {
	// 创建 hertz 服务实例，监听 8888 端口
	port := utils.Get(constant.EnvKeyPort, "8888")
	shutdownTimeout := utils.GetAsInt(constant.EnvKeyShutdownTimeout, 30)
	if shutdownTimeout <= 0 {
		shutdownTimeout = 30
	}
	h := server.New(
		server.WithHostPorts(":"+port),
		server.WithMaxRequestBodySize(11*1024*1024),
		server.WithExitWaitTime(time.Duration(shutdownTimeout)*time.Second),
	)
	clientIPFunc, err := buildClientIPFunc(
		utils.Get(constant.EnvKeyTrustedProxyCIDRs),
	)
	if err != nil {
		return nil, err
	}
	h.SetClientIPFunc(clientIPFunc)
	//`Use 所有 HTTP 请求，在进入路由处理函数之前都会先走注册的中间件。
	//recovery.Recovery() 捕获 handler 中发生的 panic，防止整个服务进程崩溃！
	h.Use(recovery.Recovery())

	v1 := h.Group(constant.APIPrefix)
	apiv1.RegisterRoutes(v1)

	//创建健康检查路由
	h.GET("/healthz", func(ctx context.Context, c *app.RequestContext) {
		resps.Ok(c, resps.Success, map[string]string{"status": "ok"})
	})

	return h, nil
}

func Run() error {
	h, err := newServer()
	if err != nil {
		return err
	}

	h.Spin()
	return nil
}
