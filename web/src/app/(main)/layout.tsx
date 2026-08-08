import type { ReactNode } from 'react'
import Link from 'next/link'
import { UserNavigation } from '@/components/layout/user-navigation'

interface MainLayoutProps {
  children: ReactNode
}

export default function MainLayout({
  children,
}: MainLayoutProps) {
  return (
    <div className="flex min-h-screen flex-col">
      <header className="border-b border-black/10 dark:border-white/10">
        <div className="mx-auto flex min-h-16 w-full max-w-6xl flex-wrap items-center gap-x-8 gap-y-3 px-4 py-3 sm:px-6">
          <Link href="/" className="text-lg font-semibold">
            My Blog
          </Link>

          <nav
            aria-label="前台导航"
            className="flex flex-1 flex-wrap items-center gap-x-5 gap-y-2 text-sm"
          >
            <Link href="/">首页</Link>
            <Link href="/?sort=latest">文章</Link>
          </nav>
          <UserNavigation />
        </div>
      </header>

      <main className="mx-auto w-full max-w-6xl flex-1 px-4 py-10 sm:px-6">
        {children}
      </main>

      <footer className="border-t border-black/10 dark:border-white/10">
        <div className="mx-auto w-full max-w-6xl px-4 py-6 sm:px-6">
          <p className="text-sm text-gray-500">
            &copy; 2026 My Blog. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  )
}
