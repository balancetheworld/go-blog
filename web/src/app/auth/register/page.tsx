'use client'

import type { FormEvent } from 'react'
import axios from 'axios'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import {
  register,
  requestEmailVerify,
} from '@/api/user'
import { Captcha } from '@/components/auth/captcha'
import { useAuth } from '@/contexts/auth-context'

export default function RegisterPage() {
  const router = useRouter()
  const { refreshUser } = useAuth()
  const [email, setEmail] = useState('')
  const [captchaToken, setCaptchaToken] = useState('')
  const [captchaRequired, setCaptchaRequired] = useState(true)
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
      toast.error('请先填写邮箱')
      return
    }

    setSendingCode(true)

    try {
      await requestEmailVerify({ email })
      setCodeCooldown(60)
      toast.success('验证码已发送')
    }
    catch (error) {
      if (
        axios.isAxiosError(error)
        && error.response?.status === 429
      ) {
        setCodeCooldown(60)
        toast.error('请求过于频繁，请稍后再试')
        return
      }

      toast.error('验证码发送失败')
    }
    finally {
      setSendingCode(false)
    }
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const formData = new FormData(event.currentTarget)

    setSubmitting(true)

    try {
      await register({
        email,
        emailCode: String(formData.get('emailCode')),
        username: String(formData.get('username')),
        nickname: String(formData.get('nickname')),
        password: String(formData.get('password')),
        captchaToken: captchaToken || undefined,
      })
      await refreshUser()

      toast.success('注册成功')
      router.replace('/')
      router.refresh()
    }
    catch {
      toast.error('注册失败，请检查填写内容')
    }
    finally {
      setSubmitting(false)
    }
  }

  return (
    <main>
      <form onSubmit={handleSubmit}>
        <h1>注册</h1>

        <label htmlFor="email">邮箱</label>
        <div>
          <input
            id="email"
            name="email"
            type="email"
            autoComplete="email"
            required
            value={email}
            onChange={event => setEmail(event.target.value)}
          />
          <button
            type="button"
            disabled={!email || sendingCode || codeCooldown > 0}
            onClick={handleSendCode}
          >
            {sendingCode
              ? '发送中...'
              : codeCooldown > 0
                ? `${codeCooldown} 秒后重试`
                : '发送验证码'}
          </button>
        </div>

        <label htmlFor="emailCode">邮箱验证码</label>
        <input
          id="emailCode"
          name="emailCode"
          type="text"
          inputMode="numeric"
          autoComplete="one-time-code"
          minLength={6}
          maxLength={6}
          required
        />

        <label htmlFor="username">用户名</label>
        <input
          id="username"
          name="username"
          type="text"
          autoComplete="username"
          minLength={3}
          maxLength={32}
          required
        />

        <label htmlFor="nickname">昵称</label>
        <input
          id="nickname"
          name="nickname"
          type="text"
          maxLength={64}
        />

        <label htmlFor="password">密码</label>
        <input
          id="password"
          name="password"
          type="password"
          autoComplete="new-password"
          minLength={8}
          maxLength={72}
          required
        />

        <Captcha
          onTokenChange={setCaptchaToken}
          onRequiredChange={setCaptchaRequired}
        />

        <button
          type="submit"
          disabled={submitting || (captchaRequired && !captchaToken)}
        >
          {submitting ? '注册中...' : '注册'}
        </button>
      </form>
    </main>
  )
}
