package controller

import (
	"context"
	"errors"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zyj/my-blog/pkg/resps"
	"github.com/zyj/my-blog/pkg/utils"
)

func GetImage(ctx context.Context, c *app.RequestContext) {
	path, err := utils.ResolveImagePath(c.Param("filepath"))
	if err != nil {
		c.AbortWithStatus(consts.StatusNotFound)
		return
	}

	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() {
		c.AbortWithStatus(consts.StatusNotFound)
		return
	}

	c.File(path)
}

func UploadImage(ctx context.Context, c *app.RequestContext) {
	file, err := c.FormFile("file")
	if err != nil {
		resps.BadRequest(c, "image file is required")
		return
	}

	image, err := utils.SaveImage(file)
	if errors.Is(err, utils.ErrImageTooLarge) {
		resps.BadRequest(c, "image must not exceed 10 MB")
		return
	}
	if errors.Is(err, utils.ErrInvalidImage) {
		resps.BadRequest(c, "only JPEG, PNG, GIF and WebP images are supported")
		return
	}
	if err != nil {
		resps.InternalServerError(c, "save image failed")
		return
	}

	resps.Ok(c, resps.Success, image)
}
