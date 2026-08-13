import type {
  Category,
  CreateCategoryReq,
  CreateLabelReq,
  CreatePostReq,
  Label,
  Post,
  PostListRequest,
  PostListResponse,
  UpdateCategoryReq,
  UpdateLabelReq,
  UpdatePostReq,
} from '@/models/post'
import type { Resp } from '@/models/resp'
import { axiosClient } from '@/lib/api/client'

function responseData<T>(response: Resp<T>): T {
  if (response.data === null)
    throw new Error(response.message)
  return response.data
}

export async function getPost(slugOrID: string): Promise<Post> {
  const response = await axiosClient.get<Resp<Post>>(
    `/post/p/${encodeURIComponent(slugOrID)}`,
  )
  return responseData(response.data)
}

export async function listPosts(req: PostListRequest = {}): Promise<PostListResponse> {
  const response = await axiosClient.get<Resp<PostListResponse>>(
    '/post/list',
    {
      params: req,
    },
  )
  return responseData(response.data)
}

export async function getRandomPost(): Promise<Post> {
  const response = await axiosClient.get<Resp<Post>>(
    '/post/random',
  )

  return responseData(response.data)
}

export async function createPost(req: CreatePostReq): Promise<Post> {
  const response = await axiosClient.post<Resp<Post>>(
    '/post/p',
    req,
  )
  return responseData(response.data)
}

export async function updatePost(
  id: number,
  req: UpdatePostReq,
): Promise<Post> {
  const response = await axiosClient.put<Resp<Post>>(
    `/post/p/${id}`,
    req,
  )
  return responseData(response.data)
}

export async function deletePost(id: number): Promise<void> {
  await axiosClient.delete(`/post/p/${id}`)
}

export async function createCategory(
  req: CreateCategoryReq,
): Promise<Category> {
  const response = await axiosClient.post<Resp<Category>>(
    '/post/c',
    req,
  )

  return responseData(response.data)
}

export async function updateCategory(
  id: number,
  req: UpdateCategoryReq,
): Promise<Category> {
  const response = await axiosClient.put<Resp<Category>>(
    `/post/c/${id}`,
    req,
  )

  return responseData(response.data)
}

export async function deleteCategory(id: number): Promise<void> {
  await axiosClient.delete(`/post/c/${id}`)
}

export async function createLabel(req: CreateLabelReq): Promise<Label> {
  const response = await axiosClient.post<Resp<Label>>(
    '/post/l',
    req,
  )

  return responseData(response.data)
}

export async function updateLabel(
  id: number,
  req: UpdateLabelReq,
): Promise<Label> {
  const response = await axiosClient.put<Resp<Label>>(
    `/post/l/${id}`,
    req,
  )

  return responseData(response.data)
}

export async function deleteLabel(id: number): Promise<void> {
  await axiosClient.delete(`/post/l/${id}`)
}
