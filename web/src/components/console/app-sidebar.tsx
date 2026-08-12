'use client'

import type { UserRole } from '@/models/user'
import {
  FileText,
  LayoutDashboard,
  ShieldCheck,
  UserCheck,
} from 'lucide-react'
import Link from 'next/link'
import { NavGroup } from './nav-group'

interface AppSidebarProps {
  role: UserRole
}

const contentItems = [
  {
    title: '控制台',
    href: '/console',
    icon: LayoutDashboard,
  },
  {
    title: '文章管理',
    href: '/console/posts',
    icon: FileText,
  },
]

const adminItems = [
  {
    title: '身份管理',
    href: '/console/roles',
    icon: ShieldCheck,
  },
  {
    title: '身份审核',
    href: '/console/roles/applications',
    icon: UserCheck,
  },
]

export function AppSidebar({
  role,
}: AppSidebarProps) {
  return (
    <aside className="border-b border-black/10 p-4 dark:border-white/10 md:sticky md:top-0 md:h-screen md:border-r md:border-b-0">
      <Link href="/console" className="text-lg font-semibold">
        管理后台
      </Link>

      <div className="mt-6 space-y-6">
        <NavGroup title="内容管理" items={contentItems} />

        {role === 'admin' && (
          <NavGroup title="系统管理" items={adminItems} />
        )}
      </div>
    </aside>
  )
}
