'use client'

import type { FormEvent } from 'react'
import type { CommentTargetType } from '@/models/comment'
import { useTranslations } from 'next-intl'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { createComment } from '@/api/comment'
import { useAuth } from '@/contexts/auth-context'

interface CommentInputProps {
  targetType: CommentTargetType
  targetID: number
  onCreated: () => void | Promise<void>
  onCancel?: () => void
  replyToName?: string
  autoFocus?: boolean
}

export function CommentInput({
  targetType,
  targetID,
  onCreated,
  onCancel,
  replyToName,
  autoFocus = false,
}: CommentInputProps) {
  const t = useTranslations('Comments')
  const pathname = usePathname()
  const { currentUser, isLoading } = useAuth()
  const [content, setContent] = useState('')
  const [submitting, setSubmitting] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const value = content.trim()
    if (!value || submitting)
      return

    setSubmitting(true)

    try {
      await createComment({
        targetType,
        targetId: targetID,
        content: value,
      })
      setContent('')
      await onCreated()
      toast.success(replyToName ? t('replyPublished') : t('published'))
    }
    catch {
      toast.error(replyToName ? t('replyFailed') : t('publishFailed'))
    }
    finally {
      setSubmitting(false)
    }
  }

  if (isLoading)
    return <div className="comment-state" />

  if (!currentUser) {
    return (
      <div className="comment-login-prompt">
        <span>{t('loginPrompt')}</span>
        <Link
          href={`/auth/login?next=${encodeURIComponent(pathname)}`}
        >
          {t('goLogin')}
        </Link>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="comment-form">
      <textarea
        aria-label={replyToName ? t('replyTo', { name: replyToName }) : t('write')}
        placeholder={replyToName ? t('replyPlaceholder', { name: replyToName }) : t('placeholder')}
        autoFocus={autoFocus}
        rows={replyToName ? 3 : 4}
        required
        maxLength={2000}
        value={content}
        onChange={event => setContent(event.target.value)}
      />
      <div className="comment-form-actions">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            disabled={submitting}
            className="comment-secondary-action"
          >
            {t('cancel')}
          </button>
        )}
        <button
          type="submit"
          disabled={submitting || !content.trim()}
        >
          {submitting ? t('submitting') : replyToName ? t('submitReply') : t('submit')}
        </button>
      </div>
    </form>
  )
}
