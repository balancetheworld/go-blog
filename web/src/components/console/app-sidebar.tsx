'use client'

import type { UserRole } from '@/models/user'
import { ArrowUpRight, Menu, X } from 'lucide-react'
import { useTranslations } from 'next-intl'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useEffect, useState } from 'react'
import { isAdmin } from '@/lib/permission'
import { consoleNavGroups } from './data'
import { NavGroup } from './nav-group'

interface AppSidebarProps {
  role: UserRole
}

export function AppSidebar({
  role,
}: AppSidebarProps) {
  const t = useTranslations('Console')
  const pathname = usePathname()
  const [open, setOpen] = useState(false)

  useEffect(() => {
    setOpen(false)
  }, [pathname])

  useEffect(() => {
    if (!open)
      return

    const previousOverflow = document.body.style.overflow

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape')
        setOpen(false)
    }

    document.body.style.overflow = 'hidden'
    window.addEventListener('keydown', handleKeyDown)

    return () => {
      document.body.style.overflow = previousOverflow
      window.removeEventListener('keydown', handleKeyDown)
    }
  }, [open])

  return (
    <>
      <button
        type="button"
        className="console-sidebar-toggle"
        aria-label={open ? t('closeMenu') : t('openMenu')}
        aria-expanded={open}
        onClick={() => setOpen(value => !value)}
      >
        {open ? <X aria-hidden="true" /> : <Menu aria-hidden="true" />}
      </button>
      <button
        type="button"
        className={`console-sidebar-overlay${open ? ' is-open' : ''}`}
        aria-label={t('closeMenu')}
        tabIndex={open ? 0 : -1}
        onClick={() => setOpen(false)}
      />
      <aside className={`console-sidebar${open ? ' is-open' : ''}`}>
        <div className="console-brand">
          <Link href="/console" className="console-brand-link">
            <span className="console-brand-mark">Z</span>
            <span>
              <strong>Caitria Console</strong>
              <small>{t('brandSubtitle')}</small>
            </span>
          </Link>
        </div>

        <div className="console-sidebar-nav">
          {consoleNavGroups.map(group => (
            (!group.adminOnly || isAdmin(role)) && (
              <NavGroup
                key={group.titleKey}
                titleKey={group.titleKey}
                items={group.items}
              />
            )
          ))}
        </div>

        <Link href="/" className="console-view-site">
          <span>{t('viewSite')}</span>
          <ArrowUpRight aria-hidden="true" />
        </Link>
      </aside>
    </>
  )
}
