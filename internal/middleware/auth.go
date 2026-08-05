package middleware

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/resps"
	"github.com/zyj/my-blog/pkg/utils"
	"gorm.io/gorm"
)

const (
	currentUserIDKey    = "current_user_id"
	currentRoleKey      = "current_role"
	currentSessionIDKey = "current_session_id"
)

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

func SetGuest(c *app.RequestContext) {
	//往本次请求的上下文里存入一个键值对，只在当前这一次请求内生效
	c.Set(currentUserIDKey, uint(0))
	c.Set(currentRoleKey, constant.RoleGuest)
	c.Set(currentSessionIDKey, "")
}

func SetCurrentUser(
	c *app.RequestContext,
	userID uint,
	role constant.Role,
	sessionID string,
) {
	c.Set(currentUserIDKey, userID)
	c.Set(currentRoleKey, role)
	c.Set(currentSessionIDKey, sessionID)
}

func GetCurrentUserID(
	c *app.RequestContext,
) (uint, bool) {
	value, exists := c.Get(currentUserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint)
	return userID, ok && userID > 0
}

func GetCurrentSessionID(
	c *app.RequestContext,
) (string, bool) {
	value, exists := c.Get(currentSessionIDKey)
	if !exists {
		return "", false
	}

	sessionID, ok := value.(string)
	return sessionID, ok && sessionID != ""
}

func GetCurrentRole(
	c *app.RequestContext,
) constant.Role {
	value, exists := c.Get(currentRoleKey)
	if !exists {
		return constant.RoleGuest
	}

	role, ok := value.(constant.Role)
	if !ok {
		return constant.RoleGuest
	}

	return role
}

func SetTokenCookies(
	c *app.RequestContext,
	accessToken string,
	refreshToken string,
	accessMaxAge int,
	refreshMaxAge int,
) {
	secure := constant.Mode(
		utils.Get(
			constant.EnvKeyMode,
			string(constant.ModeDev),
		),
	) == constant.ModeProd

	c.SetCookie(
		AccessTokenCookieName,
		accessToken,
		accessMaxAge,
		"/",
		"",
		protocol.CookieSameSiteLaxMode,
		secure,
		true,
	)
	c.SetCookie(
		RefreshTokenCookieName,
		refreshToken,
		refreshMaxAge,
		"/",
		"",
		protocol.CookieSameSiteLaxMode,
		secure,
		true,
	)
}

func ClearTokenCookies(c *app.RequestContext) {
	SetTokenCookies(c, "", "", -1, -1)
}

func UseAuth(block bool) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		accessClaims, err := utils.ParseAccessToken(readAccessToken(c))
		if err == nil {
			SetCurrentUser(
				c,
				accessClaims.UserID,
				accessClaims.Role,
				accessClaims.SessionID,
			)
			c.Next(ctx)
			return
		}

		refreshClaims, err := utils.ParseRefreshToken(
			string(c.Cookie(RefreshTokenCookieName)),
		)
		if err == nil {
			valid, sessionErr := repo.IsSessionValidForUser(
				ctx,
				refreshClaims.SessionID,
				refreshClaims.UserID,
			)
			if sessionErr != nil {
				abortInternalServerError(c)
				return
			}

			if valid {
				user, userErr := repo.GetUserByID(
					ctx,
					uint64(refreshClaims.UserID),
				)
				if userErr == nil {
					tokens, tokenErr := utils.GenerateTokenPair(
						user.ID,
						user.Role,
						refreshClaims.SessionID,
						isRememberedRefresh(refreshClaims),
					)
					if tokenErr != nil {
						abortInternalServerError(c)
						return
					}

					SetTokenCookies(
						c,
						tokens.AccessToken,
						tokens.RefreshToken,
						tokens.AccessMaxAge,
						tokens.RefreshMaxAge,
					)
					SetCurrentUser(
						c,
						user.ID,
						user.Role,
						refreshClaims.SessionID,
					)
					c.Next(ctx)
					return
				}

				if !errors.Is(userErr, gorm.ErrRecordNotFound) {
					abortInternalServerError(c)
					return
				}
			}
		}

		ClearTokenCookies(c)
		SetGuest(c)
		if block {
			c.Abort()
			resps.Unauthorized(c, resps.ErrUnauthorized)
			return
		}

		c.Next(ctx)
	}
}

func readAccessToken(c *app.RequestContext) string {
	if value := c.Cookie(AccessTokenCookieName); len(value) > 0 {
		return string(value)
	}

	fields := strings.Fields(string(c.GetHeader("Authorization")))
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
		return ""
	}

	return fields[1]
}

func isRememberedRefresh(claims *utils.TokenClaims) bool {
	if claims.IssuedAt == nil || claims.ExpiresAt == nil {
		return false
	}

	regularDuration := time.Duration(
		utils.GetAsInt(
			constant.EnvKeyRefreshTokenDuration,
			604800,
		),
	) * time.Second
	return claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time) > regularDuration
}

func abortInternalServerError(c *app.RequestContext) {
	c.Abort()
	resps.InternalServerError(c, "internal server error")
}

func roleLevel(role constant.Role) int {
	switch role {
	case constant.RoleGuest:
		return 0
	case constant.RoleUser:
		return 1
	case constant.RoleEditor:
		return 2
	case constant.RoleAdmin:
		return 3
	default:
		return -1
	}
}

func UseRole(
	requiredRole constant.Role,
) app.HandlerFunc {
	requiredLevel := roleLevel(requiredRole)

	return func(
		ctx context.Context,
		c *app.RequestContext,
	) {
		currentLevel := roleLevel(
			GetCurrentRole(c),
		)

		if requiredLevel < 0 ||
			currentLevel < requiredLevel {
			c.Abort()
			resps.Forbidden(
				c,
				resps.ErrForbidden,
			)
			return
		}

		c.Next(ctx)
	}
}
