'use client'

import type { User } from '@/models/user'
import * as Avatar from '@radix-ui/react-avatar'
import { LogOut, UserRound } from 'lucide-react'
import { useTranslations } from 'next-intl'
import { useRouter } from 'next/navigation'
import { LogoutDialog } from '@/components/auth/logout-dialog'

interface NavUserProps {
  user: Pick<User, 'username' | 'nickname' | 'avatar'>
}

export function NavUser({
  user,
}: NavUserProps) {
  const t = useTranslations('Common')
  const router = useRouter()

  return (
    <div className="console-user">
      <div className="console-user-profile">
        <Avatar.Root className="console-user-avatar">
          <Avatar.Image
            src={user.avatar}
            alt={user.nickname || user.username}
            className="console-user-image"
          />
          <Avatar.Fallback delayMs={200}>
            <UserRound
              className="console-user-fallback"
              aria-hidden="true"
            />
          </Avatar.Fallback>
        </Avatar.Root>
        <span className="console-user-name">
          {user.nickname || user.username}
        </span>
      </div>

      <LogoutDialog
        onLoggedOut={() => {
          router.replace('/auth/login')
          router.refresh()
        }}
        trigger={(
          <button
            type="button"
            aria-label={t('logout')}
            title={t('logout')}
            className="console-icon-button"
          >
            <LogOut aria-hidden="true" />
          </button>
        )}
      />
    </div>
  )
}
