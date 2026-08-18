'use client'

import type { FormEvent } from 'react'
import type { RoleOption } from '@/models/role'
import axios from 'axios'
import { useTranslations } from 'next-intl'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import { getRequestableRoles } from '@/api/role'
import {
  register,
  requestEmailVerify,
} from '@/api/user'
import { Captcha } from '@/components/auth/captcha'
import { GithubLoginButton } from '@/components/auth/github-login-button'
import { FilterSelect } from '@/components/design/filter-select'
import { useAuth } from '@/contexts/auth-context'

export default function RegisterPage() {
  const t = useTranslations('Auth.register')
  const router = useRouter()
  const { refreshUser } = useAuth()
  const [email, setEmail] = useState('')
  const [captchaToken, setCaptchaToken] = useState('')
  const [captchaRequired, setCaptchaRequired] = useState(true)
  const [sendingCode, setSendingCode] = useState(false)
  const [codeCooldown, setCodeCooldown] = useState(0)
  const [submitting, setSubmitting] = useState(false)
  const [roles, setRoles] = useState<RoleOption[]>([])
  const [loadingRoles, setLoadingRoles] = useState(true)

  useEffect(() => {
    let cancelled = false
    async function loadRoles() {
      try {
        const data = await getRequestableRoles()
        if (!cancelled)
          setRoles(data)
      }
      catch {
        if (!cancelled)
          toast.error(t('rolesFailed'))
      }
      finally {
        if (!cancelled)
          setLoadingRoles(false)
      }
    }
    void loadRoles()

    return () => {
      cancelled = true
    }
  }, [])

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
      await requestEmailVerify({ email, purpose: 'register' })
      setCodeCooldown(60)
      toast.success(t('codeSent'))
    }
    catch (error) {
      if (
        axios.isAxiosError(error)
        && error.response?.status === 429
      ) {
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
    const requestedRoleValue = String(
      formData.get('requestedRoleId') ?? '',
    )

    setSubmitting(true)

    try {
      await register({
        email,
        emailCode: String(formData.get('emailCode')),
        username: String(formData.get('username')),
        password: String(formData.get('password')),
        captchaToken: captchaToken || undefined,
        requestedRoleId: requestedRoleValue && requestedRoleValue !== 'none'
          ? Number(requestedRoleValue)
          : undefined,
      })
      await refreshUser()

      toast.success(t('success'))
      router.replace('/')
      router.refresh()
    }
    catch (error) {
      if (axios.isAxiosError(error)) {
        const message = error.response?.data?.message
        if (message === 'email already exists') {
          toast.error(t('emailExists'))
          return
        }
        if (message === 'username already exists') {
          toast.error(t('usernameExists'))
          return
        }
        if (message === 'invalid or expired email verification code') {
          toast.error(t('invalidEmailCode'))
          return
        }
        if (message === 'requested role is not available') {
          toast.error(t('roleUnavailable'))
          return
        }
        if (message === 'registration is disabled') {
          toast.error(t('registrationDisabled'))
          return
        }
      }
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
            <GithubLoginButton label={t('github')} />
            <div className="auth-divider"><span>{t('or')}</span></div>

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

            <label className="auth-field" htmlFor="emailCode">
              <span>{t('emailCode')}</span>
              <input
                id="emailCode"
                name="emailCode"
                type="text"
                inputMode="numeric"
                autoComplete="one-time-code"
                placeholder={t('emailCodePlaceholder')}
                minLength={6}
                maxLength={6}
                required
              />
            </label>

            <label className="auth-field" htmlFor="username">
              <span>{t('username')}</span>
              <input
                id="username"
                name="username"
                type="text"
                autoComplete="username"
                placeholder={t('usernamePlaceholder')}
                minLength={3}
                maxLength={32}
                required
              />
            </label>

            <div className="auth-field">
              <span>{t('requestedRole')}</span>
              <FilterSelect
                name="requestedRoleId"
                defaultValue="none"
                ariaLabel={t('requestedRole')}
                disabled={loadingRoles}
                options={[
                  { value: 'none', label: t('defaultRole') },
                  ...roles.map(role => ({
                    value: String(role.id),
                    label: role.name,
                  })),
                ]}
              />
            </div>

            <label className="auth-field" htmlFor="password">
              <span>{t('password')}</span>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="new-password"
                placeholder={t('passwordPlaceholder')}
                minLength={8}
                maxLength={72}
                required
              />
            </label>

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
            <Link href="/auth/login" className="auth-submit auth-submit-ghost">{t('backToLogin')}</Link>
          </form>
        </div>
      </section>
    </main>
  )
}
