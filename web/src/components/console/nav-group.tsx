'use client'

import type { ConsoleNavItem } from './data'
import { useTranslations } from 'next-intl'
import Link from 'next/link'
import { usePathname } from 'next/navigation'

interface NavGroupProps {
  titleKey: string
  items: ConsoleNavItem[]
}

export function NavGroup({
  titleKey,
  items,
}: NavGroupProps) {
  const groupT = useTranslations('Console.groups')
  const navigationT = useTranslations('Console.navigation')
  const pathname = usePathname()
  const title = groupT(titleKey)

  return (
    <div className="console-nav-group">
      <p className="console-nav-title">
        {title}
      </p>

      <nav aria-label={title} className="console-nav-list">
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
              className={`console-nav-link${active ? ' is-active' : ''}`}
            >
              <Icon aria-hidden="true" />
              <span>{navigationT(item.titleKey)}</span>
            </Link>
          )
        })}
      </nav>
    </div>
  )
}
