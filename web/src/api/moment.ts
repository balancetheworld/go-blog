import type {
  CreateMomentRequest,
  ListMomentsRequest,
  ListMomentsResponse,
  Moment,
} from '@/models/moment'
import type { Resp } from '@/models/resp'
import { axiosClient } from '@/lib/api/client'
import { unwrapResponse } from '@/lib/api/response'

export async function listMoments(
  req: ListMomentsRequest = {},
): Promise<ListMomentsResponse> {
  const response = await axiosClient.get<Resp<ListMomentsResponse>>(
    '/moment/list',
    { params: req },
  )

  return unwrapResponse(response.data)
}

export async function createMoment(
  req: CreateMomentRequest,
): Promise<Moment> {
  const response = await axiosClient.post<Resp<Moment>>('/moment', req)
  return unwrapResponse(response.data)
}

export async function deleteMoment(id: number): Promise<void> {
  await axiosClient.delete(`/moment/${id}`)
}
