import type {
  ListMomentsRequest,
  ListMomentsResponse,
} from '@/models/moment'
import { serverGet } from '@/lib/api/server'

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
  return serverGet<ListMomentsResponse>(createMomentListPath(req))
}
