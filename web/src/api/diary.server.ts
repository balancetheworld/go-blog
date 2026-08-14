import type {
  Diary,
  DiaryFolder,
  ListDiariesRequest,
  ListDiariesResponse,
  ListDiaryFoldersResponse,
} from '@/models/diary'
import { cache } from 'react'
import { ServerApiError, serverGet } from '@/lib/api/server'

export { ServerApiError as DiaryServerError }

function createDiaryListPath(req: ListDiariesRequest): string {
  const params = new URLSearchParams()
  if (req.page !== undefined)
    params.set('page', String(req.page))
  if (req.pageSize !== undefined)
    params.set('page_size', String(req.pageSize))
  if (req.status)
    params.set('status', req.status)
  if (req.keyword)
    params.set('keyword', req.keyword)
  if (req.folderId !== undefined)
    params.set('folder_id', String(req.folderId))

  const query = params.toString()
  return `/diary/list${query ? `?${query}` : ''}`
}

export async function listDiaries(
  req: ListDiariesRequest = {},
): Promise<ListDiariesResponse> {
  return serverGet<ListDiariesResponse>(createDiaryListPath(req))
}

export const getDiary = cache(async (
  idOrSlug: string | number,
): Promise<Diary> => {
  return serverGet<Diary>(
    `/diary/${encodeURIComponent(String(idOrSlug))}`,
  )
})

export async function listDiaryFolders(all = false): Promise<DiaryFolder[]> {
  const result = await serverGet<ListDiaryFoldersResponse>(
    `/diary/folders${all ? '?all=true' : ''}`,
  )
  return result.items
}
