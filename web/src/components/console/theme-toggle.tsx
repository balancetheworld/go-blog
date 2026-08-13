'use client'

import { Moon, Sun } from 'lucide-react'
import { useTheme } from 'next-themes'
import { useEffect, useState } from 'react'

export function ThemeToggle() {
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
      aria-label={dark ? '切换浅色主题' : '切换深色主题'}
      title={dark ? '切换浅色主题' : '切换深色主题'}
      className="inline-flex size-9 items-center justify-center rounded-md border border-black/10 disabled:opacity-50 dark:border-white/10"
    >
      {dark
        ? <Sun className="size-4" aria-hidden="true" />
        : <Moon className="size-4" aria-hidden="true" />}
    </button>
  )
}
