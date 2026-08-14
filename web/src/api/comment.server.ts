import type {
  AdminCommentListRequest,
  CommentListResponse,
} from '@/models/comment'
import { serverGet } from '@/lib/api/server'

export async function listAdminComments(
  req: AdminCommentListRequest = {},
): Promise<CommentListResponse> {
  const params = new URLSearchParams()

  if (req.targetType)
    params.set('target_type', req.targetType)
  if (req.targetId !== undefined)
    params.set('target_id', String(req.targetId))
  if (req.authorId !== undefined)
    params.set('author_id', String(req.authorId))
  if (req.keyword)
    params.set('keyword', req.keyword)
  if (req.page !== undefined)
    params.set('page', String(req.page))
  if (req.pageSize !== undefined)
    params.set('page_size', String(req.pageSize))

  const query = params.toString()
  const path = query
    ? `/admin/comment/list?${query}`
    : '/admin/comment/list'
  return serverGet<CommentListResponse>(path)
}
