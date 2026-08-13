'use client'

import type { ConsoleNavItem } from './data'
import Link from 'next/link'
import { usePathname } from 'next/navigation'

interface NavGroupProps {
  title: string
  items: ConsoleNavItem[]
}

export function NavGroup({
  title,
  items,
}: NavGroupProps) {
  const pathname = usePathname()

  return (
    <div className="space-y-2">
      <p className="px-2 text-xs text-neutral-500">
        {title}
      </p>

      <nav aria-label={title} className="space-y-1">
        {items.map((item) => {
          const active = item.href === '/console'
            ? pathname === item.href
            : pathname === item.href
              || pathname.startsWith(`${item.href}/`)
          const Icon = item.icon

          return (
            <Link
              key={item.href}
              href={item.href}
              aria-current={active ? 'page' : undefined}
              className={`flex min-h-10 items-center gap-3 rounded-md px-2 text-sm ${
                active
                  ? 'bg-black text-white dark:bg-white dark:text-black'
                  : 'text-neutral-600 hover:bg-black/5 dark:text-neutral-300 dark:hover:bg-white/10'
              }`}
            >
              <Icon className="size-4" aria-hidden="true" />
              {item.title}
            </Link>
          )
        })}
      </nav>
    </div>
  )
}
