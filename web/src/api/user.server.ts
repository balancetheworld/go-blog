import type {
  UserListRequest,
  UserListResponse,
} from '@/models/user'
import { serverGet } from '@/lib/api/server'

export async function listUsers(
  req: UserListRequest = {},
): Promise<UserListResponse> {
  const params = new URLSearchParams()

  if (req.page !== undefined)
    params.set('page', String(req.page))
  if (req.pageSize !== undefined)
    params.set('page_size', String(req.pageSize))

  const query = params.toString()
  const path = query
    ? `/users?${query}`
    : '/users'
  return serverGet<UserListResponse>(path)
}
