import type { PostDetail } from '@/models/post'
import type { Resp } from '@/models/resp'
// 从 field-conv 工具包导入 对象转换函数 snakeToCamelObj
import { snakeToCamelObj } from 'field-conv'
import { BACKEND_URL } from '@/lib/api/client'

export class PostUnavailableError extends Error {}

export async function getPublicPost(
  id: string,
): Promise<PostDetail | null> {
  const response = await fetch(
    `${BACKEND_URL}/api/v1/posts/${encodeURIComponent(id)}`,
    {
      cache: 'no-store',
    },
  )

  if (response.status === 404) {
    return null
  }

  if (response.status === 401 || response.status === 403) {
    throw new PostUnavailableError()
  }

  if (!response.ok) {
    throw new Error(`Get post failed: ${response.status}`)
  }

  const body = snakeToCamelObj(
    await response.json(),
  ) as Resp<PostDetail>

  return body.data
}
