'use client'

import { useTranslations } from 'next-intl'
import Link from 'next/link'
import { LogoutDialog } from '@/components/auth/logout-dialog'
import { useAuth } from '@/contexts/auth-context'

export function UserNavigation() {
  const t = useTranslations('Common')
  const applicationT = useTranslations('Auth.roleApplication')
  const {
    currentUser,
    currentRole,
    currentIdentity,
    roleApplication,
    isLoading,
  } = useAuth()

  if (isLoading) {
    return null
  }

  if (!currentUser) {
    return (
      <div className="flex items-center gap-3 text-sm">
        <Link href="/auth/login" className="font-medium">{t('login')}</Link>
        <Link href="/auth/register" className="font-medium">{t('register')}</Link>
      </div>
    )
  }

  return (
    <div className="flex items-center gap-3 text-sm">
      <span className="max-w-32 truncate">
        {currentUser.nickname || currentUser.username}
      </span>

      {currentIdentity && (
        <span className="text-neutral-500">
          {currentIdentity.name}
        </span>
      )}

      {roleApplication?.status === 'pending' && (
        <span className="text-amber-600">
          {applicationT('pending', { role: roleApplication.requestedRole.name })}
        </span>
      )}

      {roleApplication?.status === 'rejected' && (
        <span
          className="text-red-600"
          title={roleApplication.rejectReason || undefined}
        >
          {applicationT('rejected', { role: roleApplication.requestedRole.name })}
        </span>
      )}

      {(currentRole === 'editor' || currentRole === 'admin') && (
        <Link href="/console" className="font-medium">
          {t('console')}
        </Link>
      )}

      <LogoutDialog
        trigger={(
          <button type="button" className="font-medium">
            {t('logout')}
          </button>
        )}
      />
    </div>
  )
}
