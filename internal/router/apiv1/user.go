package apiv1

 import (
        "github.com/cloudwego/hertz/pkg/route"
        "github.com/zyj/my-blog/internal/controller"
  )

  func registerUserRoutes(group *route.RouterGroup) {
    users := group.Group("/users")
	
	// 注册用户相关路由
	users.GET("", controller.ListUsers)
	users.GET("/:id", controller.GetUser)
	users.POST("", controller.CreateUser)
	users.PUT("/:id", controller.UpdateUser)
	users.DELETE("/:id", controller.DeleteUser)
}
