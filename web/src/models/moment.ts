import type { User } from './user'

export interface Moment {
  id: number
  content: string
  author: User
  createdAt: string
  updatedAt: string
}

export interface ListMomentsRequest {
  page?: number
  pageSize?: number
}

export interface ListMomentsResponse {
  items: Moment[]
  total: number
  page: number
  pageSize: number
}

export interface CreateMomentRequest {
  content: string
}
