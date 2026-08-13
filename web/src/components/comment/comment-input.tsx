'use client'

import type { FormEvent } from 'react'
import type { CommentTargetType } from '@/models/comment'
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
      toast.success(replyToName ? '回复已发布' : '评论已发布')
    }
    catch {
      toast.error(replyToName ? '回复发布失败' : '评论发布失败')
    }
    finally {
      setSubmitting(false)
    }
  }

  if (isLoading)
    return <div className="min-h-24" />

  if (!currentUser) {
    return (
      <div className="border-y border-black/10 py-5 text-sm dark:border-white/10">
        <span className="text-neutral-600 dark:text-neutral-400">
          登录后可以参与评论。
        </span>
        <Link
          href={`/auth/login?next=${encodeURIComponent(pathname)}`}
          className="ml-2 font-medium"
        >
          前往登录
        </Link>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3 py-5">
      <label className="block text-sm font-medium">
        {replyToName ? `回复 ${replyToName}` : '发表评论'}
        <textarea
          autoFocus={autoFocus}
          rows={replyToName ? 3 : 5}
          required
          maxLength={2000}
          value={content}
          onChange={event => setContent(event.target.value)}
          className="mt-2 w-full resize-y rounded-md border border-black/20 p-3 font-normal dark:border-white/20"
        />
      </label>
      <div className="flex justify-end gap-2">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            disabled={submitting}
            className="min-h-9 rounded-md border border-black/15 px-3 text-sm disabled:opacity-50 dark:border-white/15"
          >
            取消
          </button>
        )}
        <button
          type="submit"
          disabled={submitting || !content.trim()}
          className="min-h-9 rounded-md bg-black px-4 text-sm text-white disabled:opacity-50 dark:bg-white dark:text-black"
        >
          {submitting ? '提交中' : replyToName ? '提交回复' : '提交评论'}
        </button>
      </div>
    </form>
  )
}
