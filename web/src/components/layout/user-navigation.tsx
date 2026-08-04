'use client'

import Link from 'next/link'
import { useAuth } from '@/contexts/auth-context'

export function UserNavigation() {
  const {
    currentUser,
    currentRole,
    isLoading,
    logout,
  } = useAuth()

  if (isLoading) {
    return null
  }

  if (!currentUser) {
    return (
      <div className="flex items-center gap-4">
        <Link href="/auth/login">登录</Link>
        <Link href="/auth/register">注册</Link>
      </div>
    )
  }

  return (
    <div className="flex items-center gap-4">
      <span>
        {currentUser.nickname || currentUser.username}
      </span>

      {(currentRole === 'editor' || currentRole === 'admin') && (
        <Link href="/console">
          管理后台
        </Link>
      )}

      <button
        type="button"
        onClick={() => void logout()}
      >
        退出登录
      </button>
    </div>
  )
}
