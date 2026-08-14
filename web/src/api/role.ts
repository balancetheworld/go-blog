import type { Resp } from '@/models/resp'
import type {
  CreateRoleReq,
  ListRoleApplicationsRequest,
  ListRoleApplicationsResponse,
  ListRolesRequest,
  ListRolesResponse,
  Role,
  RoleOption,
  UpdateRoleReq,
} from '@/models/role'
import { axiosClient } from '@/lib/api/client'
import { unwrapResponse } from '@/lib/api/response'

export async function getRequestableRoles(): Promise<RoleOption[]> {
  const response = await axiosClient.get<Resp<RoleOption[]>>(
    '/role/requestable',
  )
  return unwrapResponse(response.data)
}

export async function listRoles(
  req: ListRolesRequest = {},
): Promise<ListRolesResponse> {
  const response = await axiosClient.get<Resp<ListRolesResponse>>(
    '/role',
    {
      params: req,
    },
  )
  return unwrapResponse(response.data)
}

export async function createRole(req: CreateRoleReq): Promise<Role> {
  const response = await axiosClient.post<Resp<Role>>('/role', req)
  return unwrapResponse(response.data)
}

export async function updateRole(
  id: number,
  req: UpdateRoleReq,
): Promise<Role> {
  const response = await axiosClient.put<Resp<Role>>(`/role/${id}`, req)
  return unwrapResponse(response.data)
}

export async function deleteRole(id: number): Promise<void> {
  await axiosClient.delete(`/role/${id}`)
}

export async function listRoleApplications(
  req: ListRoleApplicationsRequest = {},
): Promise<ListRoleApplicationsResponse> {
  const response = await axiosClient.get<Resp<ListRoleApplicationsResponse>>(
    '/role/applications',
    {
      params: req,
    },
  )
  return unwrapResponse(response.data)
}

export async function approveRoleApplication(
  id: number,
): Promise<void> {
  await axiosClient.put(
    `/role/applications/${id}/approve`,
  )
}

export async function rejectRoleApplication(
  id: number,
  reason: string,
): Promise<void> {
  await axiosClient.put(
    `/role/applications/${id}/reject`,
    { reason },
  )
}
