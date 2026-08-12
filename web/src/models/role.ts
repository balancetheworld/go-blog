export interface RoleOption {
  id: number
  code: string
  name: string
  description: string
}

export interface Role extends RoleOption {
  isSystem: boolean
  isDefault: boolean
  isRequestable: boolean
  enabled: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateRoleReq {
  code: string
  name: string
  description: string
  isRequestable: boolean
  enabled: boolean
}

export interface UpdateRoleReq {
  name?: string
  description?: string
  isRequestable?: boolean
  enabled?: boolean
}

export interface ListRolesRequest {
  page?: number
  pageSize?: number
  keyword?: string
}

export interface ListRolesResponse {
  items: Role[]
  total: number
  page: number
  pageSize: number
}

export type RoleApplicationStatus
  = | 'pending'
    | 'approved'
    | 'rejected'

export interface RoleApplicationUser {
  id: number
  username: string
  nickname: string
}

export interface RoleApplication {
  id: number
  user: RoleApplicationUser
  requestedRole: RoleOption
  status: RoleApplicationStatus
  reviewerId: number | null
  reviewedAt: string | null
  rejectReason: string
  createdAt: string
}

export interface ListRoleApplicationsRequest {
  page?: number
  pageSize?: number
  status?: RoleApplicationStatus
}

export interface ListRoleApplicationsResponse {
  items: RoleApplication[]
  total: number
  page: number
  pageSize: number
}
