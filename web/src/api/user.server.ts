import type { Resp } from '@/models/resp'
import type {
  UserListRequest,
  UserListResponse,
} from '@/models/user'
import process from 'node:process'
import { snakeToCamelObj } from 'field-conv'
import { cookies } from 'next/headers'

const backendUrl
  = process.env.BACKEND_URL ?? 'http://localhost:8888'

export async function listUsers(
  req: UserListRequest = {},
): Promise<UserListResponse> {
  const cookieStore = await cookies()
  const params = new URLSearchParams()

  if (req.page !== undefined)
    params.set('page', String(req.page))
  if (req.pageSize !== undefined)
    params.set('page_size', String(req.pageSize))

  const query = params.toString()
  const path = query
    ? `/api/v1/users?${query}`
    : '/api/v1/users'

  const response = await fetch(`${backendUrl}${path}`, {
    cache: 'no-store',
    headers: {
      cookie: cookieStore.toString(),
    },
  })

  const body = snakeToCamelObj(
    await response.json(),
  ) as Resp<UserListResponse>

  if (!response.ok || body.data === null) {
    throw new Error(
      body.message || `List users failed: ${response.status}`,
    )
  }

  return body.data
}
