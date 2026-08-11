package controller

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/service"
	"github.com/zyj/my-blog/pkg/resps"
)

func ListRequestableRoles(
	ctx context.Context,
	c *app.RequestContext,
) {
	roles, err := service.ListRequestableRoles(ctx)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, roles)
}

func ListRoleApplications(
	ctx context.Context,
	c *app.RequestContext,
) {
	var req dto.ListRoleApplicationsRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	applications, err := service.ListRoleApplications(ctx, req)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, applications)
}
