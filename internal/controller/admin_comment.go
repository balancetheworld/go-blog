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

func ListAdminComments(ctx context.Context, c *app.RequestContext) {
	var req dto.AdminCommentListRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	comments, err := service.ListAdminComments(
		ctx,
		req,
		middleware.GetCurrentRole(c),
	)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, comments)
}

func DeleteAdminComment(ctx context.Context, c *app.RequestContext) {
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

	if err := service.DeleteComment(
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

func ModerateComment(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	var req dto.UpdateCommentModerationRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	comment, err := service.ModerateComment(
		ctx,
		uint(id),
		middleware.GetCurrentRole(c),
		req,
	)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, comment)
}
