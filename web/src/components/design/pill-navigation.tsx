'use client'

import { gsap } from 'gsap'
import {
  BookOpen,
  ChevronRight,
  Clock3,
  LayoutDashboard,
  LogIn,
  LogOut,
  Menu,
  MessageSquareText,
  Moon,
  NotebookTabs,
  Sparkles,
  X,
} from 'lucide-react'
import { useTranslations } from 'next-intl'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useEffect, useRef, useState } from 'react'
import { LogoutDialog } from '@/components/auth/logout-dialog'
import { LanguageToggle } from '@/components/i18n/language-toggle'
import { useAuth } from '@/contexts/auth-context'

export function PillNavigation() {
  const t = useTranslations('Navigation')
  const common = useTranslations('Common')
  const pathname = usePathname()
  const { currentUser, currentRole, isLoading } = useAuth()
  const [collapsed, setCollapsed] = useState(false)
  const [mobile, setMobile] = useState(false)
  const [mobileOpen, setMobileOpen] = useState(false)
  const themeTransitioning = useRef(false)

  useEffect(() => {
    const nav = document.querySelector<HTMLElement>('.pill-nav')
    const links = document.querySelector<HTMLElement>('.pill-links')
    const expand = document.querySelector<HTMLElement>('.pill-expand')
    if (!nav || !links || !expand)
      return
    const linksElement = links

    let isCollapsed = false

    function collapse() {
      if (isCollapsed)
        return
      isCollapsed = true
      gsap.to(linksElement, {
        maxWidth: 0,
        opacity: 0,
        padding: 0,
        margin: 0,
        duration: 0.4,
        ease: 'expo.inOut',
        onComplete: () => {
          setCollapsed(true)
        },
      })
      gsap.fromTo(expand, {
        rotate: -90,
        opacity: 0,
        scale: 0.8,
      }, {
        rotate: 0,
        opacity: 1,
        scale: 1,
        duration: 0.4,
        ease: 'back.out(1.7)',
        delay: 0.1,
      })
    }

    function expandNavigation() {
      if (!isCollapsed)
        return
      isCollapsed = false
      setCollapsed(false)
      gsap.to(linksElement, {
        maxWidth: linksElement.scrollWidth,
        opacity: 1,
        duration: 0.5,
        ease: 'expo.out',
      })
      gsap.fromTo('.pill-link', {
        x: -20,
        opacity: 0,
      }, {
        x: 0,
        opacity: 1,
        duration: 0.4,
        stagger: 0.05,
        ease: 'back.out(1.4)',
        delay: 0.1,
      })
    }

    function handleScroll() {
      if (window.innerWidth <= 768) {
        isCollapsed = false
        setCollapsed(false)
        gsap.set(linksElement, { clearProps: 'maxWidth,opacity,padding,margin' })
        return
      }

      if (window.scrollY > 100)
        collapse()
      else
        expandNavigation()
    }

    function handleResize() {
      const nextMobile = window.innerWidth <= 768
      setMobile(nextMobile)
      if (!nextMobile)
        setMobileOpen(false)
      handleScroll()
    }

    function handleManualExpand() {
      expandNavigation()
    }

    handleResize()
    window.addEventListener('scroll', handleScroll, { passive: true })
    window.addEventListener('resize', handleResize, { passive: true })
    window.addEventListener('design-nav-expand', handleManualExpand)
    return () => {
      window.removeEventListener('scroll', handleScroll)
      window.removeEventListener('resize', handleResize)
      window.removeEventListener('design-nav-expand', handleManualExpand)
      gsap.killTweensOf([linksElement, expand, '.pill-link'])
    }
  }, [])

  useEffect(() => {
    setMobileOpen(false)
  }, [pathname])

  useEffect(() => {
    if (!mobileOpen)
      return

    const previousOverflow = document.body.style.overflow

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape')
        setMobileOpen(false)
    }

    document.body.style.overflow = 'hidden'
    window.addEventListener('keydown', handleKeyDown)

    return () => {
      document.body.style.overflow = previousOverflow
      window.removeEventListener('keydown', handleKeyDown)
    }
  }, [mobileOpen])

  function toggleTheme() {
    if (themeTransitioning.current)
      return
    themeTransitioning.current = true

    const current = document.documentElement.dataset.theme
    const next = current === 'cat' ? 'glass' : 'cat'

    function applyTheme() {
      document.documentElement.dataset.theme = next
      window.localStorage.setItem('blog-theme', next)
    }

    function createRipple() {
      const ripple = document.createElement('div')
      ripple.className = 'theme-ripple'
      ripple.dataset.theme = next
      document.body.appendChild(ripple)
      ripple.addEventListener('animationend', () => ripple.remove(), { once: true })
    }

    const viewTransitionDocument = document as Document & {
      startViewTransition?: (callback: () => void) => {
        ready: Promise<void>
        finished: Promise<void>
      }
    }

    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches || !viewTransitionDocument.startViewTransition) {
      applyTheme()
      createRipple()
      themeTransitioning.current = false
      window.dispatchEvent(new CustomEvent('design-theme-burst'))
      return
    }

    const endRadius = Math.hypot(window.innerWidth, window.innerHeight)
    const transition = viewTransitionDocument.startViewTransition(applyTheme)
    void transition.ready.then(() => {
      createRipple()
      document.documentElement.animate({
        clipPath: [
          'circle(0px at left bottom)',
          `circle(${endRadius}px at left bottom)`,
        ],
      }, {
        duration: 560,
        easing: 'linear',
        pseudoElement: '::view-transition-new(root)',
      })
    }).catch(() => {
      applyTheme()
    })
    void transition.finished.finally(() => {
      themeTransitioning.current = false
      window.dispatchEvent(new CustomEvent('design-theme-burst'))
    })
  }

  const canAccessConsole = currentRole === 'admin' || currentRole === 'editor'

  function closeMobileNavigation() {
    setMobileOpen(false)
  }

  return (
    <>
      <nav id="pillNav" className={`pill-nav${collapsed ? ' collapsed' : ''}${mobileOpen ? ' mobile-open' : ''}`}>
        <div className="pill-nav-inner">
          <Link href="/" className="pill-brand">
            <span className="pill-avatar">
              <svg width="32" height="32" viewBox="0 0 40 40" fill="none" aria-hidden="true">
                <circle cx="20" cy="20" r="20" fill="var(--accent)" opacity="0.2" />
                <text x="20" y="26" textAnchor="middle" fill="var(--accent)" fontSize="14" fontWeight="bold" fontFamily="monospace">Z</text>
              </svg>
            </span>
            <span className="pill-name">{common('siteName')}</span>
          </Link>

          <span className="pill-divider" />

          <div id="pillLinks" className="pill-links">
            <Link href="/#articles" className={`pill-link${pathname === '/' ? ' active' : ''}`} onClick={closeMobileNavigation}>
              <BookOpen className="pill-icon" aria-hidden="true" />
              <span>{t('articles')}</span>
            </Link>
            <Link href="/#diary" className={`pill-link${pathname.startsWith('/diary') ? ' active' : ''}`} onClick={closeMobileNavigation}>
              <NotebookTabs className="pill-icon" aria-hidden="true" />
              <span>{t('diaries')}</span>
            </Link>
            <Link href="/#moments" className="pill-link" onClick={closeMobileNavigation}>
              <Sparkles className="pill-icon" aria-hidden="true" />
              <span>{t('moments')}</span>
            </Link>
            <Link href="/#timeline" className="pill-link" onClick={closeMobileNavigation}>
              <Clock3 className="pill-icon" aria-hidden="true" />
              <span>{t('timeline')}</span>
            </Link>
            <Link href="/#guestbook" className="pill-link" onClick={closeMobileNavigation}>
              <MessageSquareText className="pill-icon" aria-hidden="true" />
              <span>{t('guestbook')}</span>
            </Link>
            {!isLoading && (
              <div className="pill-mobile-account">
                {currentUser
                  ? (
                      <>
                        {canAccessConsole && (
                          <Link href="/console" className="pill-link" onClick={closeMobileNavigation}>
                            <LayoutDashboard className="pill-icon" aria-hidden="true" />
                            <span>{common('console')}</span>
                          </Link>
                        )}
                        <LogoutDialog
                          onLoggedOut={closeMobileNavigation}
                          trigger={(
                            <button type="button" className="pill-link" onClick={closeMobileNavigation}>
                              <LogOut className="pill-icon" aria-hidden="true" />
                              <span>{common('logout')}</span>
                            </button>
                          )}
                        />
                      </>
                    )
                  : (
                      <Link href="/auth/login" className="pill-link" onClick={closeMobileNavigation}>
                        <LogIn className="pill-icon" aria-hidden="true" />
                        <span>{common('login')}</span>
                      </Link>
                    )}
              </div>
            )}
          </div>

          <button
            id="pillExpand"
            type="button"
            className="pill-expand"
            aria-label={mobileOpen ? t('close') : t('open')}
            aria-expanded={mobile ? mobileOpen : !collapsed}
            onClick={() => {
              if (window.innerWidth <= 768) {
                setMobileOpen(value => !value)
                return
              }
              window.dispatchEvent(new CustomEvent('design-nav-expand'))
            }}
          >
            <ChevronRight className="pill-expand-desktop" aria-hidden="true" />
            {mobileOpen
              ? <X className="pill-expand-mobile" aria-hidden="true" />
              : <Menu className="pill-expand-mobile" aria-hidden="true" />}
          </button>

          <span className="pill-divider" />

          <button id="themeToggle" type="button" className="pill-action" aria-label={t('theme')} onClick={toggleTheme}>
            <Moon className="pill-icon" aria-hidden="true" />
          </button>

          <LanguageToggle className="pill-action" iconClassName="pill-icon" />

          {!isLoading && (
            <div className="pill-account-actions">
              <span className="pill-divider" />
              {currentUser
                ? (
                    <>
                      {canAccessConsole
                        ? (
                            <Link href="/console" className="pill-cta pill-cta-console" title={currentUser.nickname || currentUser.username}>
                              <LayoutDashboard className="pill-icon" aria-hidden="true" />
                              <span>{common('console')}</span>
                            </Link>
                          )
                        : (
                            <Link href="/" className="pill-cta">
                              <span>{currentUser.nickname || currentUser.username}</span>
                            </Link>
                          )}
                      <LogoutDialog
                        trigger={(
                          <button type="button" className="pill-action" aria-label={common('logout')}>
                            <LogOut className="pill-icon" aria-hidden="true" />
                          </button>
                        )}
                      />
                    </>
                  )
                : (
                    <Link href="/auth/login" className="pill-cta pill-cta-login">
                      <LogIn className="pill-icon" aria-hidden="true" />
                      <span>{common('login')}</span>
                    </Link>
                  )}
            </div>
          )}
        </div>
      </nav>
      <button
        type="button"
        className={`pill-nav-overlay${mobileOpen ? ' visible' : ''}`}
        aria-label={t('close')}
        tabIndex={mobileOpen ? 0 : -1}
        onClick={closeMobileNavigation}
      />
    </>
  )
}
