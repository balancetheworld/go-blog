package middleware

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/resps"
	"github.com/zyj/my-blog/pkg/utils"
)

func UseEmailVerify() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !utils.GetAsBool(constant.EnvKeyEnableEmailVerify, true) {
			c.Next(ctx)
			return
		}

		var req struct {
			Email string `json:"email"`
			Code  string `json:"code"`
		}
		if err := json.Unmarshal(c.Request.Body(), &req); err != nil {
			c.Abort()
			resps.BadRequest(c, resps.ErrParamInvalid)
			return
		}
		valid, err := repo.VerifyEmailCode(ctx, req.Email, req.Code)
		if err != nil {
			c.Abort()
			resps.InternalServerError(c, "verify email code failed")
			return
		}
		if !valid {
			c.Abort()
			resps.Unauthorized(c, "invalid or expired email verification code")
			return
		}

		c.Next(ctx)
	}
}
