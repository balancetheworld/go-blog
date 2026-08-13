import type { CurrentRole } from '@/models/user'

export function isAdmin(role: CurrentRole | null | undefined): boolean {
  return role === 'admin'
}

export function isEditor(role: CurrentRole | null | undefined): boolean {
  return role === 'editor' || isAdmin(role)
}
