import type { Resp } from '@/models/resp'

export function unwrapResponse<T>(response: Resp<T> | null | undefined): T {
  if (!response || response.data == null)
    throw new Error(response?.message || 'Response data is empty')

  return response.data
}
