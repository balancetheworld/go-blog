import type { Resp } from '@/models/resp'
import { axiosClient } from '@/lib/api/client'

export interface UploadedImage {
  url: string
  name: string
  size: number
}

export async function uploadImage(file: File): Promise<UploadedImage> {
  const form = new FormData()
  form.append('file', file)

  const response = await axiosClient.post<Resp<UploadedImage>>(
    '/file/image',
    form,
  )

  if (response.data.data === null)
    throw new Error(response.data.message)

  return response.data.data
}
