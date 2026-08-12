package controller

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/middleware"
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

func ListEnabledRoleOptions(
	ctx context.Context,
	c *app.RequestContext,
) {
	roles, err := service.ListEnabledRoleOptions(ctx)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, roles)
}

func ListRoles(
	ctx context.Context,
	c *app.RequestContext,
) {
	var req dto.ListRolesRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	roles, err := service.ListRoles(ctx, req)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, roles)
}

func CreateRole(
	ctx context.Context,
	c *app.RequestContext,
) {
	var req dto.CreateRoleRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	role, err := service.CreateRole(ctx, req)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, role)
}

func UpdateRole(
	ctx context.Context,
	c *app.RequestContext,
) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	role, err := service.UpdateRole(ctx, uint(id), req)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, role)
}

func DeleteRole(
	ctx context.Context,
	c *app.RequestContext,
) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	if err := service.DeleteRole(ctx, uint(id)); err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, nil)
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

func ApproveRoleApplication(
	ctx context.Context,
	c *app.RequestContext,
) {
	applicationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	reviewerID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		resps.Unauthorized(c, resps.ErrUnauthorized)
		return
	}

	if err := service.ApproveRoleApplication(
		ctx,
		uint(applicationID),
		reviewerID,
	); err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, nil)
}

func RejectRoleApplication(
	ctx context.Context,
	c *app.RequestContext,
) {
	applicationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	var req dto.RejectRoleApplicationRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	reviewerID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		resps.Unauthorized(c, resps.ErrUnauthorized)
		return
	}

	if err := service.RejectRoleApplication(
		ctx,
		uint(applicationID),
		reviewerID,
		req.Reason,
	); err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, nil)
}
