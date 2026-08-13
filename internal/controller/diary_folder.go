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

func ListDiaryFolders(ctx context.Context, c *app.RequestContext) {
	var req dto.ListDiaryFoldersRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	viewerRoleID, _ := middleware.GetCurrentRoleID(c)
	result, err := service.ListDiaryFolders(
		ctx,
		middleware.GetCurrentRole(c),
		viewerRoleID,
		req.All,
	)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, result)
}

func CreateDiaryFolder(ctx context.Context, c *app.RequestContext) {
	var req dto.CreateDiaryFolderRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	folder, err := service.CreateDiaryFolder(ctx, middleware.GetCurrentRole(c), req)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, folder)
}

func UpdateDiaryFolder(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	var req dto.UpdateDiaryFolderRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	folder, err := service.UpdateDiaryFolder(
		ctx,
		uint(id),
		middleware.GetCurrentRole(c),
		req,
	)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, folder)
}

func DeleteDiaryFolder(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	if err := service.DeleteDiaryFolder(
		ctx,
		uint(id),
		middleware.GetCurrentRole(c),
	); err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, nil)
}
