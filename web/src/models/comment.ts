import type { User } from './user'

export type CommentTargetType = 'post' | 'page' | 'comment' | 'diary'

export type CommentContentTargetType = Exclude<CommentTargetType, 'comment'>

export interface Comment {
  id: number
  postId?: number
  targetType: CommentContentTargetType
  targetId: number
  parentId: number | null
  rootId: number | null
  replyToUser?: User
  content: string
  author: User
  depth: number
  replyCount: number
  likeCount: number
  createdAt: string
  updatedAt: string
}

export interface CommentListRequest {
  postId?: number
  targetType?: CommentContentTargetType
  targetId?: number
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
  postId?: number
  targetType?: CommentTargetType
  targetId?: number
  parentId?: number
  content: string
}

export interface AdminCommentListRequest {
  targetType?: CommentContentTargetType
  targetId?: number
  authorId?: number
  keyword?: string
  page?: number
  pageSize?: number
}
