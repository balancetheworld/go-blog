import type {
  CreateDiaryFolderRequest,
  CreateDiaryRequest,
  Diary,
  DiaryFolder,
  ListDiariesRequest,
  ListDiariesResponse,
  ListDiaryFoldersResponse,
  UpdateDiaryFolderRequest,
  UpdateDiaryRequest,
} from '@/models/diary'
import type { Resp } from '@/models/resp'
import { axiosClient } from '@/lib/api/client'
import { unwrapResponse } from '@/lib/api/response'

export async function listDiaries(
  req: ListDiariesRequest = {},
): Promise<ListDiariesResponse> {
  const response = await axiosClient.get<Resp<ListDiariesResponse>>(
    '/diary/list',
    { params: req },
  )

  return unwrapResponse(response.data)
}

export async function getDiary(idOrSlug: string | number): Promise<Diary> {
  const response = await axiosClient.get<Resp<Diary>>(
    `/diary/${encodeURIComponent(String(idOrSlug))}`,
  )
  return unwrapResponse(response.data)
}

export async function createDiary(
  req: CreateDiaryRequest,
): Promise<Diary> {
  const response = await axiosClient.post<Resp<Diary>>('/diary', req)
  return unwrapResponse(response.data)
}

export async function updateDiary(
  id: number,
  req: UpdateDiaryRequest,
): Promise<Diary> {
  const response = await axiosClient.put<Resp<Diary>>(`/diary/${id}`, req)
  return unwrapResponse(response.data)
}

export async function deleteDiary(id: number): Promise<void> {
  await axiosClient.delete(`/diary/${id}`)
}

export async function listDiaryFolders(all = false): Promise<DiaryFolder[]> {
  const response = await axiosClient.get<Resp<ListDiaryFoldersResponse>>(
    '/diary/folders',
    { params: all ? { all: true } : undefined },
  )
  return unwrapResponse(response.data).items
}

export async function createDiaryFolder(
  req: CreateDiaryFolderRequest,
): Promise<DiaryFolder> {
  const response = await axiosClient.post<Resp<DiaryFolder>>(
    '/diary/folders',
    req,
  )
  return unwrapResponse(response.data)
}

export async function updateDiaryFolder(
  id: number,
  req: UpdateDiaryFolderRequest,
): Promise<DiaryFolder> {
  const response = await axiosClient.put<Resp<DiaryFolder>>(
    `/diary/folders/${id}`,
    req,
  )
  return unwrapResponse(response.data)
}

export async function deleteDiaryFolder(id: number): Promise<void> {
  await axiosClient.delete(`/diary/folders/${id}`)
}
