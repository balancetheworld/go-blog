package utils

//  这个文件负责：

//   - 生成随机 sessionID。
//   - 生成 access token。
//   - 生成 refresh token。
//   - 解析并校验两种 token。
//   - 判断 token 类型，避免拿 refresh token 当 access token 使用。

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zyj/my-blog/pkg/constant"
)

const (
	TokenTypeAccess  = "access"  //访问令牌（短时效）作用：调用业务接口用，请求头 Authorization 携带
	TokenTypeRefresh = "refresh" //刷新令牌（长时效）
)

// TokenClaims JWT令牌载荷结构体，存放加密在Token内的用户、会话、令牌类型信息
//
//	`TokenClaims` = JWT 的**载荷载体**
//
// JWT 分为三段：Header（加密算法）.Payload（载荷）.Signature（签名）
// 这个结构体，就是用来定义 `Payload` 里要存哪些业务数据。
// 当后端签发 token 时，会把这个结构体里所有字段加密打包进 token；
// 前端每次请求带上 token，后端解密后就能完整读出这套数据。
type TokenClaims struct {
	UserID               uint          `json:"user_id"`    // 当前登录用户ID
	Role                 constant.Role `json:"role"`       // 用户角色，枚举类型：guest/admin等
	SessionID            string        `json:"session_id"` // 关联数据库sessions表的会话唯一ID，用于撤销会话校验
	TokenType            string        `json:"token_type"` // 令牌类型，constant常量：access/refresh，区分访问令牌、刷新令牌
	jwt.RegisteredClaims               // 内嵌JWT官方标准声明，自带过期时间ExpiresAt、签发时间IssuedAt等标准字段
}

// TokenPair 登录成功后返回给前端的双令牌结构体
type TokenPair struct {
	AccessToken   string // 短时效业务访问令牌，前端调用普通接口携带
	RefreshToken  string // 长时效刷新令牌，仅用于过期后换新AccessToken，不可访问业务接口
	AccessMaxAge  int    // AccessToken有效期（单位秒），前端用于本地存储过期倒计时
	RefreshMaxAge int    // RefreshToken有效期（单位秒）
}

func GenerateSessionID() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func ValidateTokenConfig() error {
	if _, err := tokenSecret(); err != nil {
		return err
	}

	durations := []struct {
		key          string
		defaultValue string
	}{
		{constant.EnvKeyTokenDuration, "3600"},
		{constant.EnvKeyRefreshTokenDuration, "604800"},
		{constant.EnvKeyRefreshTokenDurationWithRemember, "2592000"},
	}
	for _, duration := range durations {
		value, err := strconv.Atoi(Get(duration.key, duration.defaultValue))
		if err != nil || value <= 0 {
			return fmt.Errorf("%s must be a positive integer", duration.key)
		}
	}

	return nil
}

// GenerateTokenPair 生成一对双令牌（AccessToken + RefreshToken），登录接口核心方法
// 参数：
//
//	userID: 当前登录用户ID
//	role: 用户角色 guest/user/editor/admin
//	sessionID: 本次登录会话唯一ID，关联sessions数据库表
//	remember: 是否勾选【记住我】，true则刷新令牌延长有效期
//
// 返回：TokenPair 双令牌结构体 + 生成失败错误
func GenerateTokenPair(
	userID uint,
	role constant.Role,
	sessionID string,
	remember bool,
) (TokenPair, error) {
	// 读取AccessToken过期时长：优先读取环境变量EnvKeyTokenDuration，无配置默认3600秒=1小时
	accessMaxAge := GetAsInt(
		constant.EnvKeyTokenDuration,
		3600,
	)

	// 初始化普通登录RefreshToken过期时长：默认7天 604800秒
	refreshMaxAge := GetAsInt(
		constant.EnvKeyRefreshTokenDuration,
		604800,
	)

	// 如果用户勾选记住我，覆盖为长时效刷新令牌，默认30天 2592000秒
	if remember {
		refreshMaxAge = GetAsInt(
			constant.EnvKeyRefreshTokenDurationWithRemember,
			2592000,
		)
	}

	// 调用通用生成函数，生成访问令牌，类型为access
	accessToken, err := generateToken(
		userID,
		role,
		sessionID,
		TokenTypeAccess,
		accessMaxAge,
	)
	// 生成失败直接返回空令牌+错误
	if err != nil {
		return TokenPair{}, err
	}

	// 生成刷新令牌，类型为refresh
	refreshToken, err := generateToken(
		userID,
		role,
		sessionID,
		TokenTypeRefresh,
		refreshMaxAge,
	)
	if err != nil {
		return TokenPair{}, err
	}

	// 组装双令牌结构体，同时返回两个令牌的过期时长给前端
	return TokenPair{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessMaxAge:  accessMaxAge,
		RefreshMaxAge: refreshMaxAge,
	}, nil
}

// ParseAccessToken 专门解析业务访问AccessToken，封装通用parseToken，强制校验tokenType=access
// 参数value：前端传来的token字符串
// 返回：解析后的载荷TokenClaims / 解析失败错误
func ParseAccessToken(value string) (*TokenClaims, error) {
	// 传入固定类型TokenTypeAccess，校验令牌类型不能是refresh
	return parseToken(value, TokenTypeAccess)
}

// ParseRefreshToken 专门解析刷新RefreshToken，强制校验tokenType=refresh
func ParseRefreshToken(value string) (*TokenClaims, error) {
	return parseToken(value, TokenTypeRefresh)
}

// generateToken 通用JWT生成工具函数，统一生成access/refresh两种令牌
// 参数：
//
//	userID 用户ID
//	role 用户角色
//	sessionID 会话ID
//	tokenType 令牌类型 access / refresh
//	maxAge 令牌有效时长（秒）
//
// 返回：加密后的JWT字符串，生成错误
func generateToken(
	userID uint,
	role constant.Role,
	sessionID string,
	tokenType string,
	maxAge int,
) (string, error) {
	// 参数合法性基础校验：用户ID、会话ID、有效期不能非法
	if userID == 0 || sessionID == "" || maxAge <= 0 {
		return "", errors.New("invalid token parameters")
	}

	// 校验角色是否为合法登录角色（排除guest游客）
	if !isAuthenticatedRole(role) {
		return "", errors.New("invalid user role")
	}

	// 获取JWT加密密钥，密钥不合法直接返回错误
	secret, err := tokenSecret()
	if err != nil {
		return "", err
	}

	// 获取当前系统时间，用于签发、过期时间计算
	now := time.Now()

	// 组装JWT载荷Claims，存放用户身份、会话、令牌类型、标准时间字段
	claims := TokenClaims{
		UserID:    userID,
		Role:      role,
		SessionID: sessionID,
		TokenType: tokenType,
		// 内嵌JWT官方标准注册声明
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(userID), 10),                           // 主题，存入用户ID字符串
			IssuedAt:  jwt.NewNumericDate(now),                                          // 签发时间iat
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(maxAge) * time.Second)), // 过期时间exp
		},
	}

	// 使用HS256对称加密算法，根据载荷创建JWT对象
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// 使用密钥对JWT签名，生成最终token字符串返回
	return token.SignedString(secret)
}

// parseToken 通用JWT解析校验核心函数
// 参数：
//
//	value 前端传入的token字符串
//	tokenType 期望校验的令牌类型 access/refresh
//
// 返回：解析完成的TokenClaims载荷，各类校验失败/解析错误
func parseToken(
	value string,
	tokenType string,
) (*TokenClaims, error) {
	// 校验token非空
	if value == "" {
		return nil, errors.New("token is empty")
	}

	// 读取JWT加密密钥
	secret, err := tokenSecret()
	if err != nil {
		return nil, err
	}

	// 声明结构体接收解析后的载荷数据
	claims := &TokenClaims{}

	// 解析token，填充claims，同时配置解析校验规则
	token, err := jwt.ParseWithClaims(
		value,
		claims,
		// 回调函数，返回解密用的密钥
		func(token *jwt.Token) (any, error) {
			return secret, nil
		},
		// 强制校验加密算法只能是HS256，防止算法篡改攻击
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		// 强制要求token必须携带过期时间，不允许永不过期令牌
		jwt.WithExpirationRequired(),
	)
	// 解析失败（签名错误、过期、格式错误）直接返回
	if err != nil {
		return nil, err
	}

	// 校验token整体合法（签名有效、未过期）
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 关键安全校验：令牌类型匹配，禁止用refresh访问业务接口，禁止access刷新令牌
	if claims.TokenType != tokenType {
		return nil, errors.New("invalid token type")
	}

	// 二次校验载荷关键字段完整性，防止伪造残缺token
	if claims.UserID == 0 || claims.SessionID == "" {
		return nil, errors.New("invalid token claims")
	}

	// 校验token内携带的角色是合法登录角色
	if !isAuthenticatedRole(claims.Role) {
		return nil, errors.New("invalid token role")
	}

	// 全部校验通过，返回载荷数据供业务读取用户信息
	return claims, nil
}

// tokenSecret 读取环境变量中的JWT加密密钥，并校验密钥安全长度
// 返回：字节数组密钥，密钥长度不足32位时报错
func tokenSecret() ([]byte, error) {
	// 读取配置文件/环境变量中JWT密钥字符串
	secret := []byte(Get(constant.EnvKeyJWTSecret))

	// 安全规范：密钥至少32位字符，短密钥极易被暴力破解
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET must contain at least 32 characters")
	}

	return secret, nil
}

// isAuthenticatedRole 校验角色是否为已登录授权角色（排除游客guest）
// 仅user、editor、admin三种角色允许签发、解析登录令牌
func isAuthenticatedRole(role constant.Role) bool {
	return role == constant.RoleUser ||
		role == constant.RoleEditor ||
		role == constant.RoleAdmin
}
