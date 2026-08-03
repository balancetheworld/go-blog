import type { ReactNode } from 'react'
import Link from 'next/link'

interface MainLayoutProps {
  children: ReactNode
}

export default function MainLayout({
  children,
}: MainLayoutProps) {
  return (
    <div>
      <header>
        <div>
          <Link href="/" className="text-lg font-semibold">
            My Blog
          </Link>

          <nav
            aria-label="前台导航"
            className="flex flex-wrap items-center gap-x-5 gap-y-2"
          >
            <Link href="/">首页</Link>
            <Link href="/posts">文章</Link>
            <Link href="/diaries">日记</Link>
            <Link href="/about">关于</Link>
          </nav>
        </div>
      </header>

      <main>{children}</main>

      <footer>
        <div>
          <p className="text-sm text-gray-500">
            &copy; 2024 My Blog. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  )
}
