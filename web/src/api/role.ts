import type { Resp } from '@/models/resp'
import type { RoleOption } from '@/models/role'
import { axiosClient } from '@/lib/api/client'

function responseData<T>(response: Resp<T>): T {
  if (response.data === null)
    throw new Error(response.message)

  return response.data
}

export async function getRequestableRoles(): Promise<RoleOption[]> {
  const response = await axiosClient.get<Resp<RoleOption[]>>(
    '/role/requestable',
  )
  return responseData(response.data)
}
