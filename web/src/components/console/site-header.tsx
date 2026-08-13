'use client'

import type { User } from '@/models/user'
import { ArrowUpRight } from 'lucide-react'
import { useTranslations } from 'next-intl'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { LanguageToggle } from '@/components/i18n/language-toggle'
import { getConsolePageTitleKey } from './data'
import { NavUser } from './nav-user'
import { ThemeToggle } from './theme-toggle'

interface SiteHeaderProps {
  user: Pick<User, 'username' | 'nickname' | 'avatar'>
}

export function SiteHeader({
  user,
}: SiteHeaderProps) {
  const t = useTranslations('Console')
  const navigationT = useTranslations('Console.navigation')
  const pathname = usePathname()

  return (
    <header className="console-header">
      <div className="console-header-copy">
        <span>{t('workspace')}</span>
        <strong>{navigationT(getConsolePageTitleKey(pathname))}</strong>
      </div>
      <div className="console-header-actions">
        <Link href="/" className="console-site-link">
          <span>{t('backToSite')}</span>
          <ArrowUpRight aria-hidden="true" />
        </Link>
        <LanguageToggle className="console-icon-button" />
        <ThemeToggle />
        <NavUser user={user} />
      </div>
    </header>
  )
}
