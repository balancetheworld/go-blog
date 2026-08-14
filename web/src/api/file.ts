import type { Resp } from '@/models/resp'
import { axiosClient } from '@/lib/api/client'
import { unwrapResponse } from '@/lib/api/response'

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

  return unwrapResponse(response.data)
}
