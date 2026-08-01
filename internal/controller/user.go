package controller

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/service"
	"github.com/zyj/my-blog/pkg/resps"
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
	if err := c.BindAndValidate(&req);err != nil {
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
