package resps

import (
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zyj/my-blog/pkg/errs"
)

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func write(c *app.RequestContext, httpStatus int, code int, message string, data any) {
	c.JSON(httpStatus, Resp{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Ok(c *app.RequestContext, message string, data any) {
	write(c, consts.StatusOK, consts.StatusOK, message, data)
}

func BadRequest(c *app.RequestContext, message string) {
	write(c, consts.StatusBadRequest, consts.StatusBadRequest, message, nil)
}

func Unauthorized(c *app.RequestContext, message string) {
	write(c, consts.StatusUnauthorized, consts.StatusUnauthorized, message, nil)
}

func Forbidden(c *app.RequestContext, message string) {
	write(c, consts.StatusForbidden, consts.StatusForbidden, message, nil)
}

func InternalServerError(c *app.RequestContext, message string) {
	write(c, consts.StatusInternalServerError, consts.StatusInternalServerError, message, nil)
}

func Error(c *app.RequestContext, serviceErr error) {
	var target *errs.ServiceError
	if !errors.As(serviceErr, &target) {
		InternalServerError(c, "internal server error")
		return
	}
	write(
		c,
		target.HTTPStatus,
		target.Code,
		target.Message,
		nil,
	)
}
