package apiv1

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/zyj/my-blog/internal/controller"
	"github.com/zyj/my-blog/internal/middleware"
	"github.com/zyj/my-blog/pkg/constant"
)

func registerMomentRoutes(group *route.RouterGroup) {
	moment := group.Group("/moment")
	moment.GET("/list", middleware.UseAuth(false), controller.ListMoments)

	editor := group.Group(
		"/moment",
		middleware.UseAuth(true),
		middleware.UseRole(constant.RoleEditor),
	)
	editor.POST("", controller.CreateMoment)
	editor.DELETE("/:id", controller.DeleteMoment)
}
