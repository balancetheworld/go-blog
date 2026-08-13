package apiv1

import "github.com/cloudwego/hertz/pkg/route"

// RegisterRoutes 统一注册所有业务路由
func RegisterRoutes(group *route.RouterGroup) {
    // 把子模块：用户相关路由注册到当前路由分组下
	// user 模块
    registerUserRoutes(group)
	registerRoleRoutes(group)
    registerPostRoutes(group)
	registerCommentRoutes(group)
	registerDiaryRoutes(group)
	registerFileRoutes(group)
}
