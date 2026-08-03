import type { Resp } from '@/models/resp'
import type { CurrentUser } from '@/models/user'
import axios from 'axios'
import { axiosClient } from '@/lib/api/client'

export interface LoginRequest {
  account: string
  password: string
  rememberMe: boolean
  captchaToken?: string
}

export interface RegisterRequest {
  email: string
  emailCode: string
  username: string
  nickname: string
  password: string
  captchaToken?: string
}

export async function login(request: LoginRequest): Promise<void> {
  await axiosClient.post('/auth/login', request)
}

export async function register(request: RegisterRequest): Promise<void> {
  await axiosClient.post('/auth/register', request)
}

export async function sendRegisterEmailCode(email: string): Promise<void> {
  await axiosClient.post('/auth/register/email-code', {
    email,
  })
}

export async function logout(): Promise<void> {
  await axiosClient.post('/auth/logout')
}

export async function fetchCurrentUser(): Promise<CurrentUser | null> {
  try {
    const response = await axiosClient.get<Resp<CurrentUser>>('/auth/me')
    return response.data.data
  }
  catch (error) {
    if (
      axios.isAxiosError(error)
      && error.response?.status === 401
    ) {
      return null
    }

    throw error
  }
}
