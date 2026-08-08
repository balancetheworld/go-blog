package apiv1

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/zyj/my-blog/internal/controller"
	"github.com/zyj/my-blog/internal/middleware"
)

func registerCommentRoutes(group *route.RouterGroup) {
	comment := group.Group("/comment")
	comment.GET("/list", middleware.UseAuth(false), controller.ListComments)
	comment.POST("", middleware.UseAuth(true), controller.CreateComment)
	comment.DELETE("/:id", middleware.UseAuth(true), controller.DeleteComment)
}
