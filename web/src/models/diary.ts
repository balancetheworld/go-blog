import type { RoleOption } from './role'
import type { User } from './user'

export type DiaryStatus = 'published' | 'draft'
export type DiaryVisibility = 'public' | 'roles' | 'private'

export interface DiaryFolder {
  id: number
  name: string
  slug: string
  description: string
  cover: string
  sort: number
  visibility: DiaryVisibility
  visibleRoles: RoleOption[]
  createdAt: string
  updatedAt: string
}

export interface Diary {
  id: number
  title: string
  slug: string
  description: string
  cover: string
  content: string
  draftContent?: string
  author: User
  folder: DiaryFolder | null
  visibility: DiaryVisibility
  visibleRoles: RoleOption[]
  viewCount: number
  likeCount: number
  commentCount: number
  status: DiaryStatus
  publishedAt: string | null
  createdAt: string
  updatedAt: string
}

export interface ListDiariesRequest {
  page?: number
  pageSize?: number
  status?: DiaryStatus | 'all'
  keyword?: string
  folderId?: number
}

export interface ListDiariesResponse {
  items: Diary[]
  total: number
  page: number
  pageSize: number
}

export interface CreateDiaryRequest {
  title: string
  slug?: string
  description: string
  cover: string
  folderId?: number
  draftContent: string
  publish: boolean
  visibility: DiaryVisibility
  visibleRoleIds: number[]
}

export interface UpdateDiaryRequest {
  title?: string
  slug?: string
  description?: string
  cover?: string
  folderId?: number
  clearFolder?: boolean
  draftContent?: string
  publish?: boolean
  visibility?: DiaryVisibility
  visibleRoleIds?: number[]
}

export interface CreateDiaryFolderRequest {
  name: string
  slug?: string
  description: string
  cover: string
  sort: number
  visibility: DiaryVisibility
  visibleRoleIds: number[]
}

export interface UpdateDiaryFolderRequest {
  name?: string
  slug?: string
  description?: string
  cover?: string
  sort?: number
  visibility?: DiaryVisibility
  visibleRoleIds?: number[]
}

export interface ListDiaryFoldersResponse {
  items: DiaryFolder[]
}
