import type { Resp } from '@/models/resp'
import process from 'node:process'
import { snakeToCamelObj } from 'field-conv'
import { cookies } from 'next/headers'

type ServerErrorConstructor = new (status: number, message: string) => Error

export class ServerApiError extends Error {
  constructor(
    public readonly status: number,
    message: string,
  ) {
    super(message)
    this.name = 'ServerApiError'
  }
}

export async function serverGet<T>(
  path: string,
  ErrorType: ServerErrorConstructor = ServerApiError,
  backendUrl = process.env.BACKEND_URL ?? 'http://localhost:8888',
): Promise<T> {
  const cookieStore = await cookies()
  const response = await fetch(`${backendUrl}/api/v1${path}`, {
    cache: 'no-store',
    headers: {
      cookie: cookieStore.toString(),
    },
  })
  const json = await response.json().catch(() => null)
  const body = json && typeof json === 'object'
    ? snakeToCamelObj(json) as Resp<T>
    : null

  if (!response.ok) {
    throw new ErrorType(
      response.status,
      body?.message || `Request failed: ${response.status}`,
    )
  }
  if (!body || body.data == null) {
    throw new ErrorType(
      response.status,
      body?.message || 'Response data is empty',
    )
  }

  return body.data
}
