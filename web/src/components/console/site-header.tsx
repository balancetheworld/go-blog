'use client'

import type { User } from '@/models/user'
import { usePathname } from 'next/navigation'
import { getConsolePageTitle } from './data'
import { NavUser } from './nav-user'
import { ThemeToggle } from './theme-toggle'

interface SiteHeaderProps {
  user: Pick<User, 'username' | 'nickname' | 'avatar'>
}

export function SiteHeader({
  user,
}: SiteHeaderProps) {
  const pathname = usePathname()

  return (
    <header className="flex min-h-16 items-center justify-between gap-4 border-b border-black/10 px-4 sm:px-6 dark:border-white/10">
      <span className="font-medium">{getConsolePageTitle(pathname)}</span>
      <div className="flex items-center gap-2">
        <ThemeToggle />
        <NavUser user={user} />
      </div>
    </header>
  )
}
