import type {
  AdminCommentListRequest,
  CommentListResponse,
} from '@/models/comment'
import type { Resp } from '@/models/resp'
import process from 'node:process'
import { snakeToCamelObj } from 'field-conv'
import { cookies } from 'next/headers'

const backendUrl = process.env.BACKEND_URL ?? 'http://localhost:8888'

export async function listAdminComments(
  req: AdminCommentListRequest = {},
): Promise<CommentListResponse> {
  const cookieStore = await cookies()
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
    ? `/api/v1/admin/comment/list?${query}`
    : '/api/v1/admin/comment/list'
  const response = await fetch(`${backendUrl}${path}`, {
    cache: 'no-store',
    headers: {
      cookie: cookieStore.toString(),
    },
  })
  const body = snakeToCamelObj(
    await response.json(),
  ) as Resp<CommentListResponse>

  if (!response.ok || body.data === null) {
    throw new Error(
      body.message || `List comments failed: ${response.status}`,
    )
  }

  return body.data
}
