package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/resps"
	"github.com/zyj/my-blog/pkg/utils"
)

func UseCaptcha() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := string(c.GetHeader("X-Captcha-Token"))
		mode := constant.Mode(
			utils.Get(
				constant.EnvKeyMode,
				string(constant.ModeDev),
			),
		)
		devPasscode := utils.Get(constant.EnvKeyCaptchaDevPasscode)
		if mode == constant.ModeDev && devPasscode != "" && token == devPasscode {
			c.Next(ctx)
			return
		}

		valid, err := utils.VerifyCaptcha(ctx, token, c.ClientIP())
		if err != nil {
			c.Abort()
			resps.InternalServerError(c, "captcha verification failed")
			return
		}
		if !valid {
			c.Abort()
			resps.Unauthorized(c, "invalid captcha")
			return
		}

		c.Next(ctx)
	}
}
