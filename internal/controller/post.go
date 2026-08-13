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

  func GetPostDetail(ctx context.Context, c *app.RequestContext) {
	slugOrID := c.Param("slug_or_id")
        viewerID, _ := middleware.GetCurrentUserID(c)
        role := middleware.GetCurrentRole(c)
		viewerRoleID, _ := middleware.GetCurrentRoleID(c)

        post, err := service.GetPostDetail(
                ctx,
                slugOrID,
                viewerID,
                role,
				viewerRoleID,
        )
        if err != nil {
                resps.Error(c, err)
                return
        }

        resps.Ok(c, resps.Success, post)
  }

  func ListPosts(
        ctx context.Context,
        c *app.RequestContext,
  ) {
        var req dto.PostListRequest
        if err := c.BindAndValidate(&req); err != nil {
                resps.BadRequest(c, resps.ErrParamInvalid)
                return
        }

        viewerID, _ := middleware.GetCurrentUserID(c)
        role := middleware.GetCurrentRole(c)
		viewerRoleID, _ := middleware.GetCurrentRoleID(c)

        posts, err := service.ListPosts(
                ctx,
                req,
                viewerID,
                role,
				viewerRoleID,
        )
        if err != nil {
                resps.Error(c, err)
                return
        }

        resps.Ok(c, resps.Success, posts)
  }

  func GetRandomPost(
        ctx context.Context,
        c *app.RequestContext,
  ) {
		viewerID, _ := middleware.GetCurrentUserID(c)
		viewerRoleID, _ := middleware.GetCurrentRoleID(c)
		post, err := service.GetRandomPost(
			ctx,
			viewerID,
			middleware.GetCurrentRole(c),
			viewerRoleID,
		)
        if err != nil {
                resps.Error(c, err)
                return
        }

        resps.Ok(c, resps.Success, post)
  }

	func ListCategories(
        ctx context.Context,
        c *app.RequestContext,
  ) {
        categories, err := service.ListCategories(ctx)
        if err != nil {
                resps.Error(c, err)
                return
        }

		resps.Ok(c, resps.Success, categories)
	  }

	func ListLabels(
		ctx context.Context,
		c *app.RequestContext,
	) {
		labels, err := service.ListLabels(ctx)
		if err != nil {
			resps.Error(c, err)
			return
		}

		resps.Ok(c, resps.Success, labels)
	}

	func CreateLabel(
		ctx context.Context,
		c *app.RequestContext,
	) {
		if _, ok := middleware.GetCurrentUserID(c); !ok {
			resps.Unauthorized(c, resps.ErrUnauthorized)
			return
		}

		var req dto.CreateLabelRequest
		if err := c.BindAndValidate(&req); err != nil {
			resps.BadRequest(c, resps.ErrParamInvalid)
			return
		}

		label, err := service.CreateLabel(
			ctx,
			middleware.GetCurrentRole(c),
			req,
		)
		if err != nil {
			resps.Error(c, err)
			return
		}

		resps.Ok(c, resps.Success, label)
	}

	func UpdateLabel(
		ctx context.Context,
		c *app.RequestContext,
	) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil || id == 0 {
			resps.BadRequest(c, resps.ErrParamInvalid)
			return
		}

		if _, ok := middleware.GetCurrentUserID(c); !ok {
			resps.Unauthorized(c, resps.ErrUnauthorized)
			return
		}

		var req dto.UpdateLabelRequest
		if err := c.BindAndValidate(&req); err != nil {
			resps.BadRequest(c, resps.ErrParamInvalid)
			return
		}

		label, err := service.UpdateLabel(
			ctx,
			uint(id),
			middleware.GetCurrentRole(c),
			req,
		)
		if err != nil {
			resps.Error(c, err)
			return
		}

		resps.Ok(c, resps.Success, label)
	}

	func DeleteLabel(
		ctx context.Context,
		c *app.RequestContext,
	) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil || id == 0 {
			resps.BadRequest(c, resps.ErrParamInvalid)
			return
		}

		if _, ok := middleware.GetCurrentUserID(c); !ok {
			resps.Unauthorized(c, resps.ErrUnauthorized)
			return
		}

		if err := service.DeleteLabel(
			ctx,
			uint(id),
			middleware.GetCurrentRole(c),
		); err != nil {
			resps.Error(c, err)
			return
		}

		resps.Ok(c, resps.Success, nil)
	}

  func CreatePost(
        ctx context.Context,
        c *app.RequestContext,
  ) {
        userID, ok := middleware.GetCurrentUserID(c)
        if !ok {
                resps.Unauthorized(c, resps.ErrUnauthorized)
                return
        }

        var req dto.CreatePostRequest
        if err := c.BindAndValidate(&req); err != nil {
                resps.BadRequest(c, resps.ErrParamInvalid)
                return
        }

        post, err := service.CreatePost(
                ctx,
                userID,
                middleware.GetCurrentRole(c),
                req,
        )
        if err != nil {
                resps.Error(c, err)
                return
        }

        resps.Ok(c, resps.Success, post)
  }

  func UpdatePost(
        ctx context.Context,
        c *app.RequestContext,
  ) {
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

        var req dto.UpdatePostRequest
        if err := c.BindAndValidate(&req); err != nil {
                resps.BadRequest(c, resps.ErrParamInvalid)
                return
        }

        post, err := service.UpdatePost(
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

        resps.Ok(c, resps.Success, post)
  }

  func DeletePost(
        ctx context.Context,
        c *app.RequestContext,
  ) {
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

        if err := service.DeletePost(
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

  func CreateCategory(
        ctx context.Context,
        c *app.RequestContext,
  ) {
        if _, ok := middleware.GetCurrentUserID(c); !ok {
                resps.Unauthorized(c, resps.ErrUnauthorized)
                return
        }

        var req dto.CreateCategoryRequest
        if err := c.BindAndValidate(&req); err != nil {
                resps.BadRequest(c, resps.ErrParamInvalid)
                return
        }

        category, err := service.CreateCategory(
                ctx,
                middleware.GetCurrentRole(c),
                req,
        )
        if err != nil {
                resps.Error(c, err)
                return
        }

        resps.Ok(c, resps.Success, category)
  }

  func UpdateCategory(
        ctx context.Context,
        c *app.RequestContext,
  ) {
        id, err := strconv.ParseUint(c.Param("id"), 10, 64)
        if err != nil || id == 0 {
                resps.BadRequest(c, resps.ErrParamInvalid)
                return
        }

        if _, ok := middleware.GetCurrentUserID(c); !ok {
                resps.Unauthorized(c, resps.ErrUnauthorized)
                return
        }

        var req dto.UpdateCategoryRequest
        if err := c.BindAndValidate(&req); err != nil {
                resps.BadRequest(c, resps.ErrParamInvalid)
                return
        }

        category, err := service.UpdateCategory(
                ctx,
                uint(id),
                middleware.GetCurrentRole(c),
                req,
        )
        if err != nil {
                resps.Error(c, err)
                return
        }

        resps.Ok(c, resps.Success, category)
  }

  func DeleteCategory(
        ctx context.Context,
        c *app.RequestContext,
  ) {
        id, err := strconv.ParseUint(c.Param("id"), 10, 64)
        if err != nil || id == 0 {
                resps.BadRequest(c, resps.ErrParamInvalid)
                return
        }

        if _, ok := middleware.GetCurrentUserID(c); !ok {
                resps.Unauthorized(c, resps.ErrUnauthorized)
                return
        }

        if err := service.DeleteCategory(
                ctx,
                uint(id),
                middleware.GetCurrentRole(c),
        ); err != nil {
                resps.Error(c, err)
                return
        }

        resps.Ok(c, resps.Success, nil)
  }
