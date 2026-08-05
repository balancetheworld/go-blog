'use client'

import type { FormEvent } from 'react'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { Captcha } from '@/components/auth/captcha'
import { useAuth } from '@/contexts/auth-context'

export default function LoginPage() {
  const router = useRouter()
  const [account, setAccount] = useState('')
  const [password, setPassword] = useState('')
  const [rememberMe, setRememberMe] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [captchaToken, setCaptchaToken] = useState('')
  const [captchaRequired, setCaptchaRequired] = useState(true)
  const { login } = useAuth()

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
      toast.success('登录成功')
      router.replace('/')
      router.refresh()
    }
    catch {
      toast.error('登录失败，请检查用户名和密码')
    }
    finally {
      setSubmitting(false)
    }
  }

  return (
    <main>
      <form onSubmit={handleSubmit}>
        <h1>登录</h1>
        <div>
          <label htmlFor="account">
            用户名或邮箱
          </label>
          <input
            id="account"
            name="account"
            type="text"
            autoComplete="username"
            required
            value={account}
            onChange={event => setAccount(event.target.value)}
            className="mt-2 h-10 w-full border border-black/20 px-3 dark:border-white/20"
          />
        </div>
        <div>
          <label htmlFor="password">
            密码
          </label>
          <input
            id="password"
            name="password"
            type="password"
            autoComplete="current-password"
            required
            value={password}
            onChange={event => setPassword(event.target.value)}
            className="mt-2 h-10 w-full border border-black/20 px-3 dark:border-white/20"
          />
        </div>
        <label>
          <input
            type="checkbox"
            checked={rememberMe}
            onChange={event => setRememberMe(event.target.checked)}
          />
          记住我
        </label>

        <Captcha
          onTokenChange={setCaptchaToken}
          onRequiredChange={setCaptchaRequired}
        />

        <button
          type="submit"
          disabled={submitting || (captchaRequired && !captchaToken)}
        >
          {submitting ? '登录中...' : '登录'}
        </button>
      </form>
    </main>
  )
}
