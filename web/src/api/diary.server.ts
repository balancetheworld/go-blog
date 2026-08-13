import type {
  Diary,
  DiaryFolder,
  ListDiariesRequest,
  ListDiariesResponse,
  ListDiaryFoldersResponse,
} from '@/models/diary'
import type { Resp } from '@/models/resp'
import process from 'node:process'
import { snakeToCamelObj } from 'field-conv'
import { cookies } from 'next/headers'
import { cache } from 'react'

const backendUrl = process.env.BACKEND_URL ?? 'http://localhost:8888'

export class DiaryServerError extends Error {
  constructor(
    public readonly status: number,
    message: string,
  ) {
    super(message)
    this.name = 'DiaryServerError'
  }
}

async function serverGet<T>(path: string): Promise<T> {
  const cookieStore = await cookies()
  const response = await fetch(`${backendUrl}/api/v1${path}`, {
    cache: 'no-store',
    headers: {
      cookie: cookieStore.toString(),
    },
  })
  const json = await response.json().catch(() => null)
  const body = json ? snakeToCamelObj(json) as Resp<T> : null

  if (!response.ok || !body || body.data === null) {
    throw new DiaryServerError(
      response.status,
      body?.message ?? `Request failed: ${response.status}`,
    )
  }

  return body.data
}

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
