import type { Resp } from '@/models/resp'
import type { CurrentUser } from '@/models/user'
import process from 'node:process'
import { snakeToCamelObj } from 'field-conv'
import { cookies } from 'next/headers'

const backendUrl
  = process.env.BACKEND_URL ?? 'http://localhost:8888'

export async function getCurrentUser(): Promise<CurrentUser | null> {
  const cookieStore = await cookies()

  const response = await fetch(`${backendUrl}/api/v1/auth/me`, {
    cache: 'no-store',
    headers: {
      cookie: cookieStore.toString(),
    },
  })

  if (response.status === 401) {
    return null
  }

  if (!response.ok) {
    throw new Error(`Get current user failed: ${response.status}`)
  }

  const body = snakeToCamelObj(
    await response.json(),
  ) as Resp<CurrentUser>

  return body.data
}
