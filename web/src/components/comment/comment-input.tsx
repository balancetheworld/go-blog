'use client'

import type { FormEvent } from 'react'
import Link from 'next/link'
import { useState } from 'react'
import { toast } from 'sonner'
import { createComment } from '@/api/comment'
import { useAuth } from '@/contexts/auth-context'

interface CommentInputProps {
  postID: number
  onCreated: () => void | Promise<void>
}

export function CommentInput({ postID, onCreated }: CommentInputProps) {
  const {
    currentUser,
    isLoading,
  } = useAuth()
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
        postId: postID,
        content: value,
      })
      setContent('')
      await onCreated()
      toast.success('评论已发布')
    }
    catch {
      toast.error('评论发布失败')
    }
    finally {
      setSubmitting(false)
    }
  }

  if (isLoading) {
    return (
      <div className="min-h-28 border-t border-black/10 py-6 dark:border-white/10" />
    )
  }
  if (!currentUser) {
    return (
      <section className="border-t border-black/10 py-6 dark:border-white/10">
        <p className="text-neutral-600 dark:text-neutral-400">
          登录后可以参与评论。
        </p>
        <Link
          href="/auth/login"
          className="mt-3 inline-block font-medium"
        >
          前往登录
        </Link>
      </section>
    )
  }
  return (
    <form onSubmit={handleSubmit} className="py-6">
      <label htmlFor="comment" className="font-medium">
        发表评论
      </label>
      <textarea
        id="comment"
        name="comment"
        rows={5}
        required
        maxLength={2000}
        value={content}
        onChange={event => setContent(event.target.value)}
        className="mt-3 w-full resize-y border border-black/20 p-3 dark:border-white/20"
      />
      <button
        type="submit"
        disabled={submitting || !content.trim()}
        className="mt-3 min-h-10 rounded-md bg-black px-4 text-sm text-white disabled:opacity-50 dark:bg-white dark:text-black"
      >
        {submitting ? '提交中' : '提交评论'}
      </button>
    </form>
  )
}
