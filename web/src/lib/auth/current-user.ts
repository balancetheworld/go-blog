import type { CurrentUserResp, User } from '@/models/user'
import { ServerApiError, serverGet } from '@/lib/api/server'

export async function getCurrentUser(): Promise<User | null> {
  try {
    const result = await serverGet<CurrentUserResp>('/user/me')
    return result.user ?? null
  }
  catch (error) {
    if (error instanceof ServerApiError && error.status === 401)
      return null
    throw error
  }
}
