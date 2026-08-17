'use client'

import type { FormEvent } from 'react'
import { useTranslations } from 'next-intl'
import Link from 'next/link'
import { useRouter, useSearchParams } from 'next/navigation'
import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import { Captcha } from '@/components/auth/captcha'
import { GithubLoginButton } from '@/components/auth/github-login-button'
import { useAuth } from '@/contexts/auth-context'

export default function LoginPage() {
  const t = useTranslations('Auth.login')
  const router = useRouter()
  const searchParams = useSearchParams()
  const [account, setAccount] = useState('')
  const [password, setPassword] = useState('')
  const [rememberMe, setRememberMe] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [captchaToken, setCaptchaToken] = useState('')
  const [captchaRequired, setCaptchaRequired] = useState(true)
  const { login } = useAuth()

  useEffect(() => {
    if (searchParams.get('oauth_error') !== 'github')
      return
    toast.error(t('githubFailed'))
    router.replace('/auth/login')
  }, [router, searchParams, t])

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setSubmitting(true)

    try {
      await login({
        account,
        password,
        rememberMe,
        captchaToken: captchaToken || undefined,
      })
      toast.success(t('success'))
      router.replace('/')
      router.refresh()
    }
    catch {
      toast.error(t('failed'))
    }
    finally {
      setSubmitting(false)
    }
  }

  return (
    <main className="auth-page-main">
      <section className="section auth-section">
        <div className="auth-shell">
          <form className="auth-card" onSubmit={handleSubmit}>
            <div className="auth-card-head">
              <span>{t('title')}</span>
              <small>{t('subtitle')}</small>
            </div>
            <label className="auth-field" htmlFor="account">
              <span>{t('account')}</span>
              <input
                id="account"
                name="account"
                type="text"
                autoComplete="username"
                placeholder={t('accountPlaceholder')}
                required
                value={account}
                onChange={event => setAccount(event.target.value)}
              />
            </label>
            <label className="auth-field" htmlFor="password">
              <span>{t('password')}</span>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="current-password"
                placeholder={t('passwordPlaceholder')}
                required
                value={password}
                onChange={event => setPassword(event.target.value)}
              />
            </label>
            <div className="auth-row">
              <label className="auth-check">
                <input
                  type="checkbox"
                  checked={rememberMe}
                  onChange={event => setRememberMe(event.target.checked)}
                />
                <span>{t('remember')}</span>
              </label>
              <Link href="/auth/reset-password" className="auth-link">{t('forgotPassword')}</Link>
            </div>

            <Captcha
              token={captchaToken}
              onTokenChange={setCaptchaToken}
              onRequiredChange={setCaptchaRequired}
            />

            <button
              type="submit"
              disabled={submitting || (captchaRequired && !captchaToken)}
              className="auth-submit"
            >
              {submitting ? t('submitting') : t('submit')}
            </button>
            <div className="auth-divider"><span>{t('or')}</span></div>
            <GithubLoginButton label={t('github')} />
            <Link href="/auth/register" className="auth-submit auth-submit-ghost">{t('registerLink')}</Link>
          </form>
        </div>
      </section>
    </main>
  )
}
