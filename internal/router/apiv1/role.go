package apiv1

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/zyj/my-blog/internal/controller"
)

func registerRoleRoutes(group *route.RouterGroup) {
	role := group.Group("/role")
	role.GET("/requestable", controller.ListRequestableRoles)
}
