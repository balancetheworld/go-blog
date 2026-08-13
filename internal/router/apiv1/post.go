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

        editor := group.Group(
                "/post",
                middleware.UseAuth(true),
                middleware.UseRole(constant.RoleEditor),
        )

        editor.POST("/p", controller.CreatePost)
        editor.PUT("/p/:id", controller.UpdatePost)
        editor.DELETE("/p/:id", controller.DeletePost)

        editor.POST("/c", controller.CreateCategory)
        editor.PUT("/c/:id", controller.UpdateCategory)
        editor.DELETE("/c/:id", controller.DeleteCategory)

		editor.POST("/l", controller.CreateLabel)
		editor.PUT("/l/:id", controller.UpdateLabel)
		editor.DELETE("/l/:id", controller.DeleteLabel)
  }
