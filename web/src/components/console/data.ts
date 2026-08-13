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
  NotebookText,
  ShieldCheck,
  Tags,
  UserCheck,
  Users,
} from 'lucide-react'

export interface ConsoleNavItem {
  title: string
  href: string
  icon: LucideIcon
}

export interface ConsoleNavGroup {
  title: string
  items: ConsoleNavItem[]
  adminOnly?: boolean
}

export const consoleNavGroups: ConsoleNavGroup[] = [
  {
    title: '内容管理',
    items: [
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
      {
        title: '草稿箱',
        href: '/console/posts/drafts',
        icon: FilePenLine,
      },
      {
        title: '日记管理',
        href: '/console/diaries',
        icon: NotebookText,
      },
      {
        title: '分类管理',
        href: '/console/categories',
        icon: FolderTree,
      },
      {
        title: '标签管理',
        href: '/console/labels',
        icon: Tags,
      },
      {
        title: '文件管理',
        href: '/console/files',
        icon: Files,
      },
    ],
  },
  {
    title: '系统管理',
    adminOnly: true,
    items: [
      {
        title: '评论管理',
        href: '/console/comments',
        icon: MessageSquare,
      },
      {
        title: '用户管理',
        href: '/console/users',
        icon: Users,
      },
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
      {
        title: '存储管理',
        href: '/console/storages',
        icon: Cloud,
      },
      {
        title: 'OIDC 配置',
        href: '/console/oidc',
        icon: Boxes,
      },
      {
        title: '全局配置',
        href: '/console/global',
        icon: Globe2,
      },
    ],
  },
]

export function getConsolePageTitle(pathname: string): string {
  const items = consoleNavGroups
    .flatMap(group => group.items)
    .sort((a, b) => b.href.length - a.href.length)
  const item = items.find(item => item.href === '/console'
    ? pathname === item.href
    : pathname === item.href || pathname.startsWith(`${item.href}/`))

  return item?.title ?? '后台管理'
}
