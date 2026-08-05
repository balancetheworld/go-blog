package apiv1

 import (
        "github.com/cloudwego/hertz/pkg/route"
        "github.com/zyj/my-blog/internal/controller"
		"github.com/zyj/my-blog/internal/middleware"
		"github.com/zyj/my-blog/pkg/constant"
  )

  func registerUserRoutes(group *route.RouterGroup) {
	user := group.Group("/user")
	user.POST("/login", middleware.UseCaptcha(), controller.Login)
	user.POST("/register", middleware.UseCaptcha(), controller.Register)
	user.POST("/email/verify", controller.VerifyEmail)
	user.GET("/captcha", controller.GetCaptchaConfig)
	user.POST("/logout", middleware.UseAuth(true), controller.Logout)
	user.GET("/me", middleware.UseAuth(false), controller.GetLoginUser)
	user.PUT("/password/edit", middleware.UseAuth(true), controller.ChangePassword)
	user.PUT("/password/reset", middleware.UseEmailVerify(), controller.ResetPassword)
	user.PUT(
		"/email/edit",
		middleware.UseAuth(true),
		middleware.UseEmailVerify(),
		controller.ChangeEmail,
	)

	users := group.Group(
		"/users",
		middleware.UseAuth(true),
		middleware.UseRole(constant.RoleAdmin),
	)

		// 注册用户相关路由
	users.GET("", controller.ListUsers)
	users.GET("/:id", controller.GetUser)
	users.POST("", controller.CreateUser)
	users.PUT("/:id", controller.UpdateUser)
	users.DELETE("/:id", controller.DeleteUser)
}
