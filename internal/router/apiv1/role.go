package apiv1

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/zyj/my-blog/internal/controller"
	"github.com/zyj/my-blog/internal/middleware"
	"github.com/zyj/my-blog/pkg/constant"
)

func registerRoleRoutes(group *route.RouterGroup) {
	role := group.Group("/role")
	role.GET("/requestable", controller.ListRequestableRoles)
	editor := group.Group(
		"/role",
		middleware.UseAuth(true),
		middleware.UseRole(constant.RoleEditor),
	)
	editor.GET("/options", controller.ListEnabledRoleOptions)
	admin := group.Group(
		"/role",
		middleware.UseAuth(true),
		middleware.UseRole(constant.RoleAdmin),
	)
	admin.GET("", controller.ListRoles)
	admin.POST("", controller.CreateRole)
	admin.GET("/applications", controller.ListRoleApplications)
	admin.PUT(
		"/applications/:id/approve",
		controller.ApproveRoleApplication,
	)
	admin.PUT(
		"/applications/:id/reject",
		controller.RejectRoleApplication,
	)
	admin.PUT("/:id", controller.UpdateRole)
	admin.DELETE("/:id", controller.DeleteRole)
}
