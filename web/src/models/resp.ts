export interface Resp<T = unknown> {
  code: number
  message: string
  data: T | null
}
