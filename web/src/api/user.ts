import type { Resp } from '@/models/resp'
import type {
  CaptchaConfig,
  CurrentUserResp,
  LoginReq,
  LoginResp,
  RegisterReq,
  RegisterResp,
  VerifyEmailReq,
} from '@/models/user'
import { axiosClient } from '@/lib/api/client'

function responseData<T>(response: Resp<T>): T {
  if (response.data === null)
    throw new Error(response.message)

  return response.data
}

function captchaHeaders(token?: string) {
  // 如果没有传入验证码token，直接返回空，不附加头部
  if (!token)
    return undefined
  // 返回HTTP请求头，把验证码凭证放到 X-Captcha-Token 自定义请求头中传给后端
  return {
    'X-Captcha-Token': token,
  }
}

export async function login(req: LoginReq): Promise<LoginResp> {
  const { captchaToken, ...data } = req
  const response = await axiosClient.post<Resp<LoginResp>>(
    '/user/login',
    data,
    { headers: captchaHeaders(captchaToken) },
  )
  return responseData(response.data)
}

export async function register(req: RegisterReq): Promise<RegisterResp> {
  const { captchaToken, ...data } = req
  const response = await axiosClient.post<Resp<RegisterResp>> (
    '/user/register',
    data,
    { headers: captchaHeaders(captchaToken) },
  )
  return responseData(response.data)
}

export async function logout(): Promise<void> {
  await axiosClient.post('/user/logout')
}

export async function getCurrentUser(): Promise<CurrentUserResp> {
  const response = await axiosClient.get<Resp<CurrentUserResp>>('/user/me')
  return responseData(response.data)
}

export async function requestEmailVerify(req: VerifyEmailReq): Promise<void> {
  await axiosClient.post('/user/email/verify', req)
}

export async function getCaptchaConfig(): Promise<CaptchaConfig> {
  const response
    = await axiosClient.get<Resp<CaptchaConfig>>('/user/captcha')
  return responseData(response.data)
}
