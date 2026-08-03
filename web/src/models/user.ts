export type UserRole = 'user' | 'editor' | 'admin'

export interface CurrentUser {
  id: number
  username: string
  nickname: string
  role: UserRole
}
