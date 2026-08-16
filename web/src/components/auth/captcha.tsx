'use client'

import type { CaptchaConfig } from '@/models/user'
import { Turnstile } from '@marsidev/react-turnstile'
import { CircleCheck, ShieldCheck, ShieldOff } from 'lucide-react'
import { useTranslations } from 'next-intl'
import { useEffect, useState } from 'react'
import { getCaptchaConfig } from '@/api/user'

interface CaptchaProps {
  token: string
  onTokenChange: (token: string) => void
  onRequiredChange: (required: boolean) => void
}

export function Captcha({
  token,
  onTokenChange,
  onRequiredChange,
}: CaptchaProps) {
  const t = useTranslations('Auth.captcha')
  const [config, setConfig] = useState<CaptchaConfig | null>(null)
  const [loadFailed, setLoadFailed] = useState(false)

  useEffect(() => {
    let active = true

    void getCaptchaConfig()
      .then((result) => {
        if (!active)
          return

        setConfig(result)
        onRequiredChange(result.provider !== 'disable')
      })
      .catch(() => {
        if (!active)
          return

        setLoadFailed(true)
        onRequiredChange(true)
      })

    return () => {
      active = false
    }
  }, [onRequiredChange])

  if (loadFailed)
    return <p role="alert">{t('loadFailed')}</p>

  if (config === null)
    return <p>{t('loading')}</p>

  if (config.provider === 'disable') {
    return (
      <p className="captcha-status is-disabled" role="status">
        <ShieldOff className="size-4" aria-hidden="true" />
        <span>{t('disabled')}</span>
      </p>
    )
  }

  if (config.provider !== 'turnstile' || !config.siteKey)
    return <p role="alert">{t('unavailable')}</p>

  return (
    <div className="captcha-field">
      <Turnstile
        siteKey={config.siteKey}
        onSuccess={onTokenChange}
        onExpire={() => onTokenChange('')}
        onError={() => onTokenChange('')}
      />
      <p className={`captcha-status${token ? ' is-verified' : ''}`} role="status">
        {token
          ? <CircleCheck className="size-4" aria-hidden="true" />
          : <ShieldCheck className="size-4" aria-hidden="true" />}
        <span>{token ? t('verified') : t('required')}</span>
      </p>
    </div>
  )
}
