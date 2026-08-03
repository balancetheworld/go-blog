'use client'

  import type { FormEvent } from 'react'
  import { useRouter } from 'next/navigation'
  import { useState } from 'react'
  import { toast } from 'sonner'
  import {
    register,
    sendRegisterEmailCode,
  } from '@/lib/api/user'

  export default function RegisterPage() {
    const router = useRouter()
    const [email, setEmail] = useState('')
    const [sendingCode, setSendingCode] = useState(false)
    const [submitting, setSubmitting] = useState(false)

    async function handleSendCode() {
      if (!email) {
        toast.error('请先填写邮箱')
        return
      }

      setSendingCode(true)

      try {
        await sendRegisterEmailCode(email)
        toast.success('验证码已发送')
      }
      catch {
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
        })

        toast.success('注册成功')
        router.replace('/auth/login')
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
              disabled={!email || sendingCode}
              onClick={handleSendCode}
            >
              {sendingCode ? '发送中...' : '发送验证码'}
            </button>
          </div>

          <label htmlFor="emailCode">邮箱验证码</label>
          <input
            id="emailCode"
            name="emailCode"
            type="text"
            inputMode="numeric"
            autoComplete="one-time-code"
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
            required/>
          <div>人机验证暂未启用</div>

          <button type="submit" disabled={submitting}>
            {submitting ? '注册中...' : '注册'}
          </button>
        </form>
      </main>
    )
  }
