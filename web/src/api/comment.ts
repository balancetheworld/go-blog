import type {
  Comment,
  CommentListRequest,
  CommentListResponse,
  CreateCommentRequest,
} from '@/models/comment'
import type { Resp } from '@/models/resp'
import { axiosClient } from '@/lib/api/client'

function responseData<T>(response: Resp<T>): T {
  if (response.data === null)
    throw new Error(response.message)
  return response.data
}

export async function listComments(
  req: CommentListRequest,
): Promise<CommentListResponse> {
  const response = await axiosClient.get<Resp<CommentListResponse>>(
    '/comment/list',
    {
      params: req,
    },
  )

  return responseData(response.data)
}

export async function createComment(
  req: CreateCommentRequest,
): Promise<Comment> {
  const response = await axiosClient.post<Resp<Comment>>(
    '/comment',
    req,
  )

  return responseData(response.data)
}

export async function deleteComment(id: number): Promise<void> {
  await axiosClient.delete(`/comment/${id}`)
}
