'use client'

import type { FormEvent } from 'react'
import axios from 'axios'
import { useTranslations } from 'next-intl'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import { requestEmailVerify, resetPassword } from '@/api/user'

export default function ResetPasswordPage() {
  const t = useTranslations('Auth.resetPassword')
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [sendingCode, setSendingCode] = useState(false)
  const [codeCooldown, setCodeCooldown] = useState(0)
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (codeCooldown <= 0)
      return

    const timer = window.setTimeout(() => {
      setCodeCooldown(value => value - 1)
    }, 1000)

    return () => window.clearTimeout(timer)
  }, [codeCooldown])

  async function handleSendCode() {
    if (!email) {
      toast.error(t('emailRequired'))
      return
    }

    setSendingCode(true)
    try {
      await requestEmailVerify({ email })
      setCodeCooldown(60)
      toast.success(t('codeSent'))
    }
    catch (error) {
      if (axios.isAxiosError(error) && error.response?.status === 429) {
        setCodeCooldown(60)
        toast.error(t('tooManyRequests'))
        return
      }
      toast.error(t('codeFailed'))
    }
    finally {
      setSendingCode(false)
    }
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)
    const newPassword = String(formData.get('newPassword'))
    const confirmPassword = String(formData.get('confirmPassword'))
    if (newPassword !== confirmPassword) {
      toast.error(t('passwordMismatch'))
      return
    }

    setSubmitting(true)
    try {
      await resetPassword({
        email,
        code: String(formData.get('code')),
        newPassword,
      })
      toast.success(t('success'))
      router.replace('/auth/login')
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
          <form className="auth-card auth-card-register" onSubmit={handleSubmit}>
            <div className="auth-card-head">
              <span>{t('title')}</span>
              <small>{t('subtitle')}</small>
            </div>

            <label className="auth-field" htmlFor="email">
              <span>{t('email')}</span>
              <span className="auth-inline-field">
                <input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="email"
                  placeholder="you@example.com"
                  required
                  value={email}
                  onChange={event => setEmail(event.target.value)}
                />
                <button
                  type="button"
                  className="auth-inline-action"
                  disabled={!email || sendingCode || codeCooldown > 0}
                  onClick={handleSendCode}
                >
                  {sendingCode
                    ? t('sendingCode')
                    : codeCooldown > 0
                      ? t('cooldown', { seconds: codeCooldown })
                      : t('sendCode')}
                </button>
              </span>
            </label>

            <label className="auth-field" htmlFor="code">
              <span>{t('code')}</span>
              <input
                id="code"
                name="code"
                type="text"
                inputMode="numeric"
                autoComplete="one-time-code"
                placeholder={t('codePlaceholder')}
                minLength={6}
                maxLength={6}
                required
              />
            </label>

            <label className="auth-field" htmlFor="newPassword">
              <span>{t('newPassword')}</span>
              <input
                id="newPassword"
                name="newPassword"
                type="password"
                autoComplete="new-password"
                placeholder={t('passwordPlaceholder')}
                minLength={8}
                maxLength={72}
                required
              />
            </label>

            <label className="auth-field" htmlFor="confirmPassword">
              <span>{t('confirmPassword')}</span>
              <input
                id="confirmPassword"
                name="confirmPassword"
                type="password"
                autoComplete="new-password"
                placeholder={t('passwordPlaceholder')}
                minLength={8}
                maxLength={72}
                required
              />
            </label>

            <button type="submit" disabled={submitting} className="auth-submit">
              {submitting ? t('submitting') : t('submit')}
            </button>
            <Link href="/auth/login" className="auth-submit auth-submit-ghost">{t('backToLogin')}</Link>
          </form>
        </div>
      </section>
    </main>
  )
}
