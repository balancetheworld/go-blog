'use client'

import type { User } from '@/models/user'
import { LogOut, UserRound } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { useAuth } from '@/contexts/auth-context'

interface NavUserProps {
  user: Pick<User, 'username' | 'nickname'>
}

export function NavUser({
  user,
}: NavUserProps) {
  const router = useRouter()
  const { logout } = useAuth()
  const [loggingOut, setLoggingOut] = useState(false)

  async function handleLogout() {
    if (loggingOut)
      return

    setLoggingOut(true)

    try {
      await logout()
      router.replace('/auth/login')
      router.refresh()
    }
    finally {
      setLoggingOut(false)
    }
  }

  return (
    <div className="flex items-center gap-3">
      <div className="flex min-w-0 items-center gap-2">
        <UserRound
          className="size-4 shrink-0 text-neutral-500"
          aria-hidden="true"
        />
        <span className="max-w-32 truncate text-sm">
          {user.nickname || user.username}
        </span>
      </div>

      <button
        type="button"
        disabled={loggingOut}
        onClick={() => void handleLogout()}
        aria-label="退出登录"
        title="退出登录"
        className="inline-flex size-9 items-center justify-center rounded-md border border-black/10 disabled:opacity-50 dark:border-white/10"
      >
        <LogOut className="size-4" aria-hidden="true" />
      </button>
    </div>
  )
}
