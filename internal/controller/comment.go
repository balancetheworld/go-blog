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

func ListComments(ctx context.Context, c *app.RequestContext) {
	var req dto.CommentListRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	viewerID, _ := middleware.GetCurrentUserID(c)
	viewerRoleID, _ := middleware.GetCurrentRoleID(c)
	comments, err := service.ListComments(
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

	resps.Ok(c, resps.Success, comments)
}

func CreateComment(ctx context.Context, c *app.RequestContext) {
	userID, ok := middleware.GetCurrentUserID(c)
	viewerRoleID, _ := middleware.GetCurrentRoleID(c)
	if !ok {
		resps.Unauthorized(c, resps.ErrUnauthorized)
		return
	}

	var req dto.CreateCommentRequest
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ErrParamInvalid)
		return
	}

	comment, err := service.CreateComment(
		ctx,
		userID,
		middleware.GetCurrentRole(c),
		viewerRoleID,
		req,
	)
	if err != nil {
		resps.Error(c, err)
		return
	}

	resps.Ok(c, resps.Success, comment)
}

func DeleteComment(ctx context.Context, c *app.RequestContext) {
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
