package apiv1

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/zyj/my-blog/internal/controller"
	"github.com/zyj/my-blog/internal/middleware"
	"github.com/zyj/my-blog/pkg/constant"
)

func registerFileRoutes(group *route.RouterGroup) {
	group.GET("/file/content/*filepath", controller.GetImage)
	group.HEAD("/file/content/*filepath", controller.GetImage)

	editor := group.Group(
		"/file",
		middleware.UseAuth(true),
		middleware.UseRole(constant.RoleEditor),
	)
	editor.POST("/image", controller.UploadImage)
}
