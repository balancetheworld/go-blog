'use client'

import type { ReactNode } from 'react'
import { GitBranch, Mail } from 'lucide-react'
import { useTranslations } from 'next-intl'
import { useEffect } from 'react'
import { DesignAnimations } from './design-animations'
import { DesignBackground } from './design-background'
import { PageScrollReset } from './page-scroll-reset'
import { ParticleCanvas } from './particle-canvas'
import { PillNavigation } from './pill-navigation'

interface DesignShellProps {
  children: ReactNode
  auth?: boolean
  footer?: boolean
}

export function DesignShell({
  children,
  auth = false,
  footer = true,
}: DesignShellProps) {
  const t = useTranslations('Footer')

  useEffect(() => {
    const saved = window.localStorage.getItem('blog-theme')
    document.documentElement.dataset.theme = saved === 'cat' ? 'cat' : 'glass'
    document.body.classList.toggle('auth-page', auth)
    document.body.classList.add('loaded')

    return () => {
      document.body.classList.remove('auth-page')
      document.body.classList.remove('loaded')
    }
  }, [auth])

  return (
    <div className={auth ? 'auth-page' : undefined}>
      <ParticleCanvas />
      <DesignBackground />
      <PageScrollReset />
      <DesignAnimations />
      <PillNavigation />
      {children}
      {footer && (
        <footer className="footer">
          <div className="footer-links" aria-label={t('contact')}>
            <a
              className="footer-link"
              href="https://github.com/balancetheworld"
              target="_blank"
              rel="noreferrer"
            >
              <GitBranch aria-hidden="true" />
              <span>GitHub</span>
            </a>
            <a className="footer-link" href="mailto:2539888062@qq.com">
              <Mail aria-hidden="true" />
              <span>2539888062@qq.com</span>
            </a>
          </div>
          <div className="footer-bottom">
            <span className="copyright">{t('copyright')}</span>
            <span className="footer-slogan">{t('slogan')}</span>
          </div>
        </footer>
      )}
    </div>
  )
}
