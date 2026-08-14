package controller

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/middleware"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/internal/service"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/resps"
	"github.com/zyj/my-blog/pkg/utils"
)

func ListUsers(ctx context.Context, c *app.RequestContext) {
	var req dto.ListUsersRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	users, err := service.ListUsers(ctx, req)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, users)
}

func GetUser(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	user, err := service.GetUser(ctx, id)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, user)
}

func CreateUser(ctx context.Context, c *app.RequestContext) {
	var req dto.CreateUserRequest
	//Hertz 提供的一体化方法，做两件事：
	// 	1. **Bind 绑定**：读取前端 POST JSON 请求体，自动映射到 `req` 结构体对应字段；
	// 2. **Validate 校验**：根据结构体 `binding` 标签校验参数合法性。
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	user, err := service.CreateUser(ctx, req)
	if err != nil {
		resps.Error(c, err)
		return
	}
	resps.Ok(c, resps.Success, user)
}

func UpdateUser(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}
	var req dto.UpdateUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	user, err := service.UpdateUser(ctx, id, req)
	if err != nil {
		resps.Error(c, err)
		return
	}
	resps.Ok(c, resps.Success, user)
}

func DeleteUser(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	if err := service.DeleteUser(ctx, id); err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, nil)
}

func Login(ctx context.Context, c *app.RequestContext) {
	var req dto.UserLoginReq
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}
	req.UserIP = c.ClientIP()
	req.UserAgent = string(c.UserAgent())

	result, err := service.UserLogin(ctx, &req)
	if err != nil {
		resps.Error(c, err)
		return
	}

	middleware.SetTokenCookies(
		c,
		result.AccessToken,
		result.RefreshToken,
		result.AccessMaxAge,
		result.RefreshMaxAge,
	)
	resps.Ok(c, resps.Success, result.User)
}

func Register(ctx context.Context, c *app.RequestContext) {
	var req dto.UserRegisterReq
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}
	req.UserIP = c.ClientIP()
	req.UserAgent = string(c.UserAgent())

	result, err := service.UserRegister(ctx, &req)
	if err != nil {
		resps.Error(c, err)
		return
	}

	middleware.SetTokenCookies(
		c,
		result.AccessToken,
		result.RefreshToken,
		result.AccessMaxAge,
		result.RefreshMaxAge,
	)
	resps.Ok(c, resps.Success, result.User)
}

func Logout(ctx context.Context, c *app.RequestContext) {
	sessionID, ok := middleware.GetCurrentSessionID(c)
	middleware.ClearTokenCookies(c)
	if !ok {
		resps.Unauthorized(c, resps.ErrUnauthorized)
		return
	}

	if err := repo.RevokeSession(ctx, sessionID); err != nil {
		resps.InternalServerError(c, "logout failed")
		return
	}

	resps.Ok(c, resps.Success, nil)
}

func GetLoginUser(ctx context.Context, c *app.RequestContext) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		resps.Ok(c, resps.Success, dto.LoginUserResponse{
			User:            nil,
			Role:            constant.RoleGuest,
			Identity:        nil,
			RoleApplication: nil,
		})
		return
	}

	result, err := service.GetLoginUser(ctx, userID)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, result)
}

func VerifyEmail(ctx context.Context, c *app.RequestContext) {
	var req dto.VerifyEmailReq
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}
	req.UserIP = c.ClientIP()

	if err := service.RequestVerifyEmail(ctx, &req); err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, nil)
}

func ChangePassword(ctx context.Context, c *app.RequestContext) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		resps.Unauthorized(c, resps.ErrUnauthorized)
		return
	}

	var req dto.UpdatePasswordReq
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	if err := service.UpdatePassword(ctx, userID, &req); err != nil {
		resps.Error(c, err)
		return
	}

	middleware.ClearTokenCookies(c)
	resps.Ok(c, resps.Success, nil)
}

func ResetPassword(ctx context.Context, c *app.RequestContext) {
	var req dto.ResetPasswordReq
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	if err := service.ResetPassword(ctx, &req); err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, nil)
}

func ChangeEmail(ctx context.Context, c *app.RequestContext) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		resps.Unauthorized(c, resps.ErrUnauthorized)
		return
	}

	var req dto.UpdateEmailReq
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	if err := service.UpdateEmail(ctx, userID, &req); err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, nil)
}

func GetCaptchaConfig(ctx context.Context, c *app.RequestContext) {
	provider := constant.CaptchaType(
		utils.Get(
			constant.EnvKeyCaptchaProvider,
			string(constant.CaptchaDisable),
		),
	)
	resps.Ok(c, resps.Success, dto.CaptchaConfigResponse{
		Provider: provider,
		SiteKey:  utils.Get(constant.EnvKeyCaptchaSiteKey),
	})
}
