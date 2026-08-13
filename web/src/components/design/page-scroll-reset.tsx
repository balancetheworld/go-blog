'use client'

import { usePathname } from 'next/navigation'
import { useLayoutEffect } from 'react'

export function PageScrollReset() {
  const pathname = usePathname()

  useLayoutEffect(() => {
    const restoreValue = window.sessionStorage.getItem('article-list-restore-scroll')
    const restorePosition = restoreValue === null ? null : Number(restoreValue)
    if (restorePosition === null && window.location.hash)
      return

    if (restorePosition !== null)
      window.sessionStorage.removeItem('article-list-restore-scroll')

    const root = document.documentElement
    const scrollBehavior = root.style.scrollBehavior
    const targetPosition = restorePosition !== null && Number.isFinite(restorePosition)
      ? restorePosition
      : 0
    let firstFrame = 0
    let secondFrame = 0

    root.style.scrollBehavior = 'auto'
    window.scrollTo(0, targetPosition)

    firstFrame = window.requestAnimationFrame(() => {
      window.scrollTo(0, targetPosition)
      secondFrame = window.requestAnimationFrame(() => {
        window.scrollTo(0, targetPosition)
        root.style.scrollBehavior = scrollBehavior
      })
    })

    return () => {
      window.cancelAnimationFrame(firstFrame)
      window.cancelAnimationFrame(secondFrame)
      root.style.scrollBehavior = scrollBehavior
    }
  }, [pathname])

  return null
}
