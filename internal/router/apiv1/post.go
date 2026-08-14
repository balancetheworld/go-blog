package apiv1

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/zyj/my-blog/internal/controller"
	"github.com/zyj/my-blog/internal/middleware"
	"github.com/zyj/my-blog/pkg/constant"
)

func registerPostRoutes(group *route.RouterGroup) {
	post := group.Group("/post")

	post.GET(
		"/p/:slug_or_id",
		middleware.UseAuth(false),
		controller.GetPostDetail,
	)
	post.GET(
		"/list",
		middleware.UseAuth(false),
		controller.ListPosts,
	)
	post.GET(
		"/categories",
		middleware.UseAuth(false),
		controller.ListCategories,
	)
	post.GET(
		"/labels",
		middleware.UseAuth(false),
		controller.ListLabels,
	)
	post.GET(
		"/random",
		middleware.UseAuth(false),
		controller.GetRandomPost,
	)
	post.POST(
		"/p/:id/like",
		middleware.UseAuth(true),
		controller.TogglePostLike,
	)

	admin := group.Group(
		"/post",
		middleware.UseAuth(true),
		middleware.UseRole(constant.RoleAdmin),
	)

	admin.POST("/p", controller.CreatePost)
	admin.PUT("/p/:id", controller.UpdatePost)
	admin.DELETE("/p/:id", controller.DeletePost)

	admin.POST("/c", controller.CreateCategory)
	admin.PUT("/c/:id", controller.UpdateCategory)
	admin.DELETE("/c/:id", controller.DeleteCategory)

	admin.POST("/l", controller.CreateLabel)
	admin.PUT("/l/:id", controller.UpdateLabel)
	admin.DELETE("/l/:id", controller.DeleteLabel)
}
