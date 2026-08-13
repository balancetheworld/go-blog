'use client'

import type { ReactNode } from 'react'
import * as Dialog from '@radix-ui/react-dialog'
import { LogOut } from 'lucide-react'
import { useTranslations } from 'next-intl'
import { useState } from 'react'
import { useAuth } from '@/contexts/auth-context'

interface LogoutDialogProps {
  trigger: ReactNode
  onLoggedOut?: () => void
}

export function LogoutDialog({
  trigger,
  onLoggedOut,
}: LogoutDialogProps) {
  const t = useTranslations('Auth.logoutDialog')
  const { logout } = useAuth()
  const [open, setOpen] = useState(false)
  const [loggingOut, setLoggingOut] = useState(false)

  async function handleLogout() {
    if (loggingOut)
      return

    setLoggingOut(true)
    try {
      await logout()
      setOpen(false)
      onLoggedOut?.()
    }
    finally {
      setLoggingOut(false)
    }
  }

  return (
    <Dialog.Root
      open={open}
      onOpenChange={(nextOpen) => {
        if (!loggingOut)
          setOpen(nextOpen)
      }}
    >
      <Dialog.Trigger asChild>
        {trigger}
      </Dialog.Trigger>
      <Dialog.Portal>
        <Dialog.Overlay className="logout-dialog-overlay" />
        <Dialog.Content className="logout-dialog-content">
          <div className="logout-dialog-icon">
            <LogOut aria-hidden="true" />
          </div>
          <Dialog.Title className="logout-dialog-title">
            {t('title')}
          </Dialog.Title>
          <Dialog.Description className="logout-dialog-description">
            {t('description')}
          </Dialog.Description>
          <div className="logout-dialog-actions">
            <Dialog.Close asChild>
              <button type="button" className="logout-dialog-cancel" disabled={loggingOut}>
                {t('cancel')}
              </button>
            </Dialog.Close>
            <button
              type="button"
              className="logout-dialog-confirm"
              disabled={loggingOut}
              onClick={() => void handleLogout()}
            >
              {loggingOut ? t('submitting') : t('confirm')}
            </button>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}
