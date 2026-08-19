import type {
  AdminCommentListRequest,
  Comment,
  CommentListRequest,
  CommentListResponse,
  CommentModerationResponse,
  CreateCommentRequest,
} from '@/models/comment'
import type { Resp } from '@/models/resp'
import { axiosClient } from '@/lib/api/client'
import { unwrapResponse } from '@/lib/api/response'

export async function listComments(
  req: CommentListRequest,
): Promise<CommentListResponse> {
  const response = await axiosClient.get<Resp<CommentListResponse>>(
    '/comment/list',
    { params: req },
  )

  return unwrapResponse(response.data)
}

export async function listCommentReplies(id: number): Promise<Comment[]> {
  const response = await axiosClient.get<Resp<Comment[]>>(
    `/comment/${id}/replies`,
  )

  return unwrapResponse(response.data)
}

export async function createComment(
  req: CreateCommentRequest,
): Promise<Comment> {
  const response = await axiosClient.post<Resp<Comment>>(
    '/comment',
    req,
  )

  return unwrapResponse(response.data)
}

export async function getCommentModeration(
  id: number,
): Promise<CommentModerationResponse> {
  const response = await axiosClient.get<Resp<CommentModerationResponse>>(
    `/comment/${id}/moderation`,
  )
  return unwrapResponse(response.data)
}

export async function deleteComment(id: number): Promise<void> {
  await axiosClient.delete(`/comment/${id}`)
}

export async function listAdminComments(
  req: AdminCommentListRequest = {},
): Promise<CommentListResponse> {
  const response = await axiosClient.get<Resp<CommentListResponse>>(
    '/admin/comment/list',
    { params: req },
  )

  return unwrapResponse(response.data)
}

export async function deleteAdminComment(id: number): Promise<void> {
  await axiosClient.delete(`/admin/comment/${id}`)
}
