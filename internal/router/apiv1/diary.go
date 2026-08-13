package apiv1

import (
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/zyj/my-blog/internal/controller"
	"github.com/zyj/my-blog/internal/middleware"
	"github.com/zyj/my-blog/pkg/constant"
)

func registerDiaryRoutes(group *route.RouterGroup) {
	diary := group.Group("/diary")
	diary.GET("/list", middleware.UseAuth(false), controller.ListDiaries)
	diary.GET("/folders", middleware.UseAuth(false), controller.ListDiaryFolders)
	diary.GET("/:id", middleware.UseAuth(false), controller.GetDiary)

	editor := group.Group(
		"/diary",
		middleware.UseAuth(true),
		middleware.UseRole(constant.RoleEditor),
	)
	editor.POST("/folders", controller.CreateDiaryFolder)
	editor.PUT("/folders/:id", controller.UpdateDiaryFolder)
	editor.DELETE("/folders/:id", controller.DeleteDiaryFolder)
	editor.POST("", controller.CreateDiary)
	editor.PUT("/:id", controller.UpdateDiary)
	editor.DELETE("/:id", controller.DeleteDiary)
}
