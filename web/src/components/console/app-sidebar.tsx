'use client'

import type { UserRole } from '@/models/user'
import Link from 'next/link'
import { isAdmin } from '@/lib/permission'
import { consoleNavGroups } from './data'
import { NavGroup } from './nav-group'

interface AppSidebarProps {
  role: UserRole
}

export function AppSidebar({
  role,
}: AppSidebarProps) {
  return (
    <aside className="border-b border-black/10 p-4 dark:border-white/10 md:sticky md:top-0 md:h-screen md:border-r md:border-b-0">
      <Link href="/console" className="text-lg font-semibold">
        管理后台
      </Link>

      <div className="mt-6 space-y-6">
        {consoleNavGroups.map(group => (
          (!group.adminOnly || isAdmin(role)) && (
            <NavGroup
              key={group.title}
              title={group.title}
              items={group.items}
            />
          )
        ))}
      </div>
    </aside>
  )
}
