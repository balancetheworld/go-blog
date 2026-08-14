import type {
  RoleApplicationStatus,
  RoleOption,
} from './role'

export type UserRole = 'user' | 'editor' | 'admin'

export type CurrentRole = UserRole | 'guest'

export interface CurrentUser {
  id: number
  username: string
  nickname: string
  role: UserRole
}

/**
 * 导出联合类型 CaptchaProvider：人机验证码服务商枚举约束
 * 联合类型 Union Type，限定变量只能赋值下方固定字符串，传入其他字符串会触发TS类型报错
 */
export type CaptchaProvider
// 不启用验证码，登录/注册页面跳过人机校验
  = | 'disable'
    // Cloudflare Turnstile 轻量人机验证，无追踪Cookie，项目常用
    | 'turnstile'
    // Google reCAPTCHA 谷歌官方人机验证
    | 'recaptcha'
    // hCaptcha 海外第三方人机验证，reCAPTCHA替代方案
    | 'hcaptcha'

export interface User {
  id: number
  username: string
  nickname: string
  avatar: string
  email?: string
  role: UserRole
  createdAt: string
  updatedAt: string
}

export interface LoginReq {
  account: string
  password: string
  rememberMe: boolean
  captchaToken?: string
}

export type LoginResp = User

export interface RegisterReq {
  email: string
  emailCode: string
  username: string
  nickname: string
  password: string
  requestedRoleId?: number
  captchaToken?: string
}

export type RegisterResp = User

export interface VerifyEmailReq {
  email: string
}

export interface ResetPasswordReq {
  email: string
  code: string
  newPassword: string
}

export interface CaptchaConfig {
  provider: CaptchaProvider
  siteKey: string
}

export interface UserRoleApplication {
  id: number
  requestedRole: RoleOption
  status: RoleApplicationStatus
  rejectReason: string
  reviewedAt: string | null
  createdAt: string
}

export interface CurrentUserResp {
  user: User | null
  role: CurrentRole
  identity: RoleOption | null
  roleApplication: UserRoleApplication | null
}

export interface UserListRequest {
  page?: number
  pageSize?: number
}

export interface UserListResponse {
  items: User[]
  total: number
  page: number
  pageSize: number
}
