'use client'

import { Moon, Sun } from 'lucide-react'
import { useTranslations } from 'next-intl'
import { useTheme } from 'next-themes'
import { useEffect, useState } from 'react'

export function ThemeToggle() {
  const t = useTranslations('Common')
  const { resolvedTheme, setTheme } = useTheme()
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  const dark = mounted && resolvedTheme === 'dark'

  return (
    <button
      type="button"
      onClick={() => setTheme(dark ? 'light' : 'dark')}
      disabled={!mounted}
      aria-label={dark ? t('switchToLight') : t('switchToDark')}
      title={dark ? t('switchToLight') : t('switchToDark')}
      className="console-icon-button"
    >
      {dark
        ? <Sun aria-hidden="true" />
        : <Moon aria-hidden="true" />}
    </button>
  )
}
