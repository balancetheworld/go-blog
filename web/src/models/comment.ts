import type { User } from './user'

export interface Comment {
  id: number
  postId: number
  content: string
  author: User
  createdAt: string
  updatedAt: string
}

export interface CommentListRequest {
  postId: number
  page?: number
  pageSize?: number
}

export interface CommentListResponse {
  items: Comment[]
  total: number
  page: number
  pageSize: number
}

export interface CreateCommentRequest {
  postId: number
  content: string
}
