import type { LucideIcon } from 'lucide-react'
import {
  Boxes,
  Cloud,
  FilePenLine,
  Files,
  FileText,
  FolderTree,
  Globe2,
  LayoutDashboard,
  MessageSquare,
  MessageSquareText,
  NotebookText,
  ShieldCheck,
  Tags,
  UserCheck,
  Users,
} from 'lucide-react'

export interface ConsoleNavItem {
  titleKey: string
  href: string
  icon: LucideIcon
}

export interface ConsoleNavGroup {
  titleKey: string
  items: ConsoleNavItem[]
  adminOnly?: boolean
}

export const consoleNavGroups: ConsoleNavGroup[] = [
  {
    titleKey: 'content',
    items: [
      {
        titleKey: 'dashboard',
        href: '/console',
        icon: LayoutDashboard,
      },
      {
        titleKey: 'posts',
        href: '/console/posts',
        icon: FileText,
      },
      {
        titleKey: 'drafts',
        href: '/console/posts/drafts',
        icon: FilePenLine,
      },
      {
        titleKey: 'diaries',
        href: '/console/diaries',
        icon: NotebookText,
      },
      {
        titleKey: 'moments',
        href: '/console/moments',
        icon: MessageSquareText,
      },
      {
        titleKey: 'categories',
        href: '/console/categories',
        icon: FolderTree,
      },
      {
        titleKey: 'labels',
        href: '/console/labels',
        icon: Tags,
      },
      {
        titleKey: 'files',
        href: '/console/files',
        icon: Files,
      },
    ],
  },
  {
    titleKey: 'system',
    adminOnly: true,
    items: [
      {
        titleKey: 'comments',
        href: '/console/comments',
        icon: MessageSquare,
      },
      {
        titleKey: 'users',
        href: '/console/users',
        icon: Users,
      },
      {
        titleKey: 'roles',
        href: '/console/roles',
        icon: ShieldCheck,
      },
      {
        titleKey: 'roleApplications',
        href: '/console/roles/applications',
        icon: UserCheck,
      },
      {
        titleKey: 'storages',
        href: '/console/storages',
        icon: Cloud,
      },
      {
        titleKey: 'oidc',
        href: '/console/oidc',
        icon: Boxes,
      },
      {
        titleKey: 'global',
        href: '/console/global',
        icon: Globe2,
      },
    ],
  },
]

export function getConsolePageTitleKey(pathname: string): string {
  const items = consoleNavGroups
    .flatMap(group => group.items)
    .sort((a, b) => b.href.length - a.href.length)
  const item = items.find(item => item.href === '/console'
    ? pathname === item.href
    : pathname === item.href || pathname.startsWith(`${item.href}/`))

  return item?.titleKey ?? 'dashboard'
}
