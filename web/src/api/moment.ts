import type {
  CreateMomentRequest,
  ListMomentsRequest,
  ListMomentsResponse,
  Moment,
} from '@/models/moment'
import type { Resp } from '@/models/resp'
import { axiosClient } from '@/lib/api/client'

function responseData<T>(response: Resp<T>): T {
  if (response.data === null)
    throw new Error(response.message)

  return response.data
}

export async function listMoments(
  req: ListMomentsRequest = {},
): Promise<ListMomentsResponse> {
  const response = await axiosClient.get<Resp<ListMomentsResponse>>(
    '/moment/list',
    { params: req },
  )

  return responseData(response.data)
}

export async function createMoment(
  req: CreateMomentRequest,
): Promise<Moment> {
  const response = await axiosClient.post<Resp<Moment>>('/moment', req)
  return responseData(response.data)
}

export async function deleteMoment(id: number): Promise<void> {
  await axiosClient.delete(`/moment/${id}`)
}
