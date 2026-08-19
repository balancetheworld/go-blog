package apiv1

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/zyj/my-blog/internal/controller"
	"github.com/zyj/my-blog/internal/middleware"
	"github.com/zyj/my-blog/pkg/constant"
)

func registerCommentRoutes(group *route.RouterGroup) {
	comment := group.Group("/comment")
	comment.GET("/list", middleware.UseAuth(false), controller.ListComments)
	comment.GET("/:id/replies", middleware.UseAuth(false), controller.ListCommentReplies)
	comment.GET("/:id/moderation", middleware.UseAuth(true), controller.GetCommentModeration)
	comment.POST("", middleware.UseAuth(true), controller.CreateComment)
	comment.DELETE("/:id", middleware.UseAuth(true), controller.DeleteComment)

	admin := group.Group(
		"/admin/comment",
		middleware.UseAuth(true),
		middleware.UseRole(constant.RoleAdmin),
	)
	admin.GET("/list", controller.ListAdminComments)
	admin.PATCH("/:id/moderation", controller.ModerateComment)
	admin.DELETE("/:id", controller.DeleteAdminComment)
}
