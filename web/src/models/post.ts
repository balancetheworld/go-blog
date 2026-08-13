import type { RoleOption } from './role'
import type { User } from './user'

export type PostPublishStatus = 'published' | 'draft'

export type PostStatus = PostPublishStatus | 'all'

export type PostSort = 'latest' | 'oldest' | 'hot'

export type PostVisibility = 'public' | 'roles' | 'private'

export interface Category {
  id: number
  name: string
  slug: string
  description: string
  createdAt: string
  updatedAt: string
}

export interface Label {
  id: number
  name: string
  slug: string
}

export interface Post {
  id: number
  title: string
  content: string
  draftContent?: string
  description: string
  cover: string
  type: string
  slug: string
  categoryId: number | null
  category?: Category
  labels: Label[]
  author: User
  isPrivate: boolean
  visibility: PostVisibility
  visibleRoles: RoleOption[]
  top: boolean
  likeCount: number
  commentCount: number
  viewCount: number
  heat: number
  status: PostPublishStatus
  publishedAt: string | null
  createdAt: string
  updatedAt: string
}

export interface PostListItem {
  id: number
  title: string
  content: string
  description: string
  cover: string
  type: string
  slug: string
  category?: Category
  labels: Label[]
  author: User
  isPrivate: boolean
  visibility: PostVisibility
  visibleRoles: RoleOption[]
  top: boolean
  likeCount: number
  commentCount: number
  viewCount: number
  heat: number
  status: PostPublishStatus
  publishedAt: string | null
  createdAt: string
}

export interface PostListRequest {
  page?: number
  pageSize?: number
  keyword?: string
  type?: string
  categoryId?: number
  labelId?: number
  authorId?: number
  status?: PostStatus
  sort?: PostSort
}

export interface PostListResponse {
  items: PostListItem[]
  total: number
  page: number
  pageSize: number
}

export interface CreatePostReq {
  title: string
  draftContent: string
  description?: string
  cover?: string
  type?: string
  categoryId?: number
  labelIds?: number[]
  isPrivate?: boolean
  visibility?: PostVisibility
  visibleRoleIds?: number[]
  top?: boolean
  publish?: boolean
}

export interface UpdatePostReq {
  title?: string
  draftContent?: string
  description?: string
  cover?: string
  type?: string
  categoryId?: number
  labelIds?: number[]
  isPrivate?: boolean
  visibility?: PostVisibility
  visibleRoleIds?: number[]
  top?: boolean
  publish?: boolean
}

export interface CreateCategoryReq {
  name: string
  slug: string
  description?: string
}

export interface UpdateCategoryReq {
  name?: string
  slug?: string
  description?: string
}

export interface CreateLabelReq {
  name: string
  slug: string
}

export interface UpdateLabelReq {
  name?: string
  slug?: string
}
