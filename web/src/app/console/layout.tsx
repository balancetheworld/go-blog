import type { ReactNode } from 'react'
import Link from 'next/link'
import { notFound, redirect } from 'next/navigation'
import { getCurrentUser } from '@/lib/auth/current-user'

interface ConsoleLayoutProps {
  children: ReactNode
}

export default async function ConsoleLayout({
  children,
}: ConsoleLayoutProps) {
  const user = await getCurrentUser()

  if (!user) {
    redirect('/auth/login?next=/console')
  }

  if (user.role !== 'admin' && user.role !== 'editor') {
    notFound()
  }

  return (
    <div className="grid min-h-screen md:grid-cols-[220px_minmax(0,1fr)]">
      <aside className="border-b border-black/10 p-4 dark:border-white/10 md:border-r md:border-b-0">
        <Link href="/console" className="text-lg font-semibold">
          管理后台
        </Link>

        <nav
          aria-label="后台导航"
          className="mt-6 flex gap-4 md:flex-col"
        >
          <Link href="/console">控制台</Link>
          <Link href="/console/posts">文章管理</Link>
        </nav>
      </aside>

      <div className="min-w-0">
        <header className="flex min-h-16 items-center justify-between border-b border-black/10 px-6 dark:border-white/10">
          <span>后台管理</span>
          <span className="text-sm text-neutral-500">
            {user.nickname || user.username}
          </span>
        </header>

        <main className="p-6">
          {children}
        </main>
      </div>
    </div>
  )
}
