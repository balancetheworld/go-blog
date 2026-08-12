import type { User } from '@/models/user'
import { NavUser } from './nav-user'

interface SiteHeaderProps {
  user: Pick<User, 'username' | 'nickname'>
}

export function SiteHeader({
  user,
}: SiteHeaderProps) {
  return (
    <header className="flex min-h-16 items-center justify-between gap-4 border-b border-black/10 px-4 sm:px-6 dark:border-white/10">
      <span className="font-medium">后台管理</span>
      <NavUser user={user} />
    </header>
  )
}
