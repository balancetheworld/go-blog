import type {
  ListMomentsRequest,
  ListMomentsResponse,
} from '@/models/moment'
import type { Resp } from '@/models/resp'
import process from 'node:process'
import { snakeToCamelObj } from 'field-conv'
import { cookies } from 'next/headers'

const backendUrl = process.env.BACKEND_URL ?? 'http://localhost:8888'

function createMomentListPath(req: ListMomentsRequest): string {
  const params = new URLSearchParams()
  if (req.page !== undefined)
    params.set('page', String(req.page))
  if (req.pageSize !== undefined)
    params.set('page_size', String(req.pageSize))

  const query = params.toString()
  return `/moment/list${query ? `?${query}` : ''}`
}

export async function listMoments(
  req: ListMomentsRequest = {},
): Promise<ListMomentsResponse> {
  const cookieStore = await cookies()
  const response = await fetch(
    `${backendUrl}/api/v1${createMomentListPath(req)}`,
    {
      cache: 'no-store',
      headers: {
        cookie: cookieStore.toString(),
      },
    },
  )
  const json = await response.json().catch(() => null)
  const body = json ? snakeToCamelObj(json) as Resp<ListMomentsResponse> : null

  if (!response.ok || !body || body.data === null)
    throw new Error(body?.message ?? `Request failed: ${response.status}`)

  return body.data
}
