'use client'

import { Globe2 } from 'lucide-react'
import { useLocale, useTranslations } from 'next-intl'
import { useRouter } from 'next/navigation'
import { useTransition } from 'react'
import { localeCookieName } from '@/i18n/config'

interface LanguageToggleProps {
  className: string
  iconClassName?: string
}

export function LanguageToggle({
  className,
  iconClassName,
}: LanguageToggleProps) {
  const locale = useLocale()
  const t = useTranslations('Common')
  const router = useRouter()
  const [pending, startTransition] = useTransition()
  const nextLocale = locale === 'en' ? 'zh-CN' : 'en'
  const label = nextLocale === 'en'
    ? t('switchToEnglish')
    : t('switchToChinese')

  function changeLanguage() {
    document.cookie = `${localeCookieName}=${nextLocale}; Path=/; Max-Age=31536000; SameSite=Lax`
    startTransition(() => {
      router.refresh()
    })
  }

  return (
    <button
      type="button"
      className={`${className}${locale === 'en' ? ' active' : ''}`}
      aria-label={label}
      title={label}
      disabled={pending}
      onClick={changeLanguage}
    >
      <Globe2 className={iconClassName} aria-hidden="true" />
    </button>
  )
}
