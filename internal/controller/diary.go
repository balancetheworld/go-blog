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

func ListDiaries(ctx context.Context, c *app.RequestContext) {
	var req dto.ListDiariesRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	viewerID, _ := middleware.GetCurrentUserID(c)
	viewerRoleID, _ := middleware.GetCurrentRoleID(c)
	result, err := service.ListDiaries(
		ctx,
		req,
		viewerID,
		middleware.GetCurrentRole(c),
		viewerRoleID,
	)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, result)
}

func GetDiary(ctx context.Context, c *app.RequestContext) {
	identifier := c.Param("id")
	if identifier == "" {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	viewerID, _ := middleware.GetCurrentUserID(c)
	viewerRoleID, _ := middleware.GetCurrentRoleID(c)
	diary, err := service.GetDiaryByIdentifier(
		ctx,
		identifier,
		viewerID,
		middleware.GetCurrentRole(c),
		viewerRoleID,
	)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, diary)
}

func CreateDiary(ctx context.Context, c *app.RequestContext) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		resps.Unauthorized(c, resps.ErrUnauthorized)
		return
	}

	var req dto.CreateDiaryRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	diary, err := service.CreateDiary(
		ctx,
		userID,
		middleware.GetCurrentRole(c),
		req,
	)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, diary)
}

func UpdateDiary(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		resps.Unauthorized(c, resps.ErrUnauthorized)
		return
	}

	var req dto.UpdateDiaryRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	diary, err := service.UpdateDiary(
		ctx,
		uint(id),
		userID,
		middleware.GetCurrentRole(c),
		req,
	)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, diary)
}

func DeleteDiary(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		resps.Unauthorized(c, resps.ErrUnauthorized)
		return
	}

	if err := service.DeleteDiary(
		ctx,
		uint(id),
		userID,
		middleware.GetCurrentRole(c),
	); err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, nil)
}
