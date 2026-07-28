package constant

// Role 角色类型，底层基于string，模拟枚举，用于区分系统用户权限身份
type Role string

// Mode 运行环境模式，区分开发环境 / 生产环境
type Mode string

// CaptchaType 验证码类型，标识后端启用哪一种人机验证方案
type CaptchaType string

// TargetType 操作目标类型，标记当前操作针对哪一类业务资源
type TargetType string

// 角色枚举常量
const (
	RoleUser   Role = "user"    // 普通用户
	RoleEditor Role = "editor"  // 编辑者，拥有内容编辑权限
	RoleAdmin  Role = "admin"   // 管理员，最高权限
)

// 程序运行模式枚举
const (
	ModeDev  Mode = "dev"   // 开发模式：开启调试日志、热更新、允许更多调试接口
	ModeProd Mode = "prod"  // 生产模式：关闭调试，严格安全校验，优化性能
)

// 验证码方案枚举
const (
	CaptchaDisable   CaptchaType = "disable"    // 禁用验证码，本地开发常用
	CaptchaTurnstile CaptchaType = "turnstile"  // Cloudflare Turnstile 人机验证
	CaptchaRecaptcha CaptchaType = "recaptcha"  // Google reCAPTCHA
	CaptchaHcaptcha  CaptchaType = "hcaptcha"   // hCaptcha 验证
)

// 业务资源目标类型枚举
const (
	TargetPost    TargetType = "post"     // 博客文章
	TargetPage    TargetType = "page"     // 独立页面
	TargetComment TargetType = "comment"  // 评论
	TargetDiary   TargetType = "diary"    // 日记
)

// APIPrefix 全局API路由统一前缀
// 所有v1版本接口统一挂载在 /api/v1 下，方便路由分组、版本管理
const APIPrefix = "/api/v1"
