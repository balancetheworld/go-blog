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
      <div className="flex items-center gap-3 text-sm">
        <Link href="/auth/login" className="font-medium">登录</Link>
        <Link href="/auth/register" className="font-medium">注册</Link>
      </div>
    )
  }

  return (
    <div className="flex items-center gap-3 text-sm">
      <span className="max-w-32 truncate">
        {currentUser.nickname || currentUser.username}
      </span>

      {(currentRole === 'editor' || currentRole === 'admin') && (
        <Link href="/console" className="font-medium">
          管理后台
        </Link>
      )}

      <button
        type="button"
        onClick={() => void logout()}
        className="font-medium"
      >
        退出登录
      </button>
    </div>
  )
}
