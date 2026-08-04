'use client'
import Link from 'next/link'
import { useAuth } from '@/contexts/auth-context'

export function CommentInput() {
  const {
    currentUser,
    isLoading,
  } = useAuth()
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
    <section className="border-t border-black/10 py-6 dark:border-white/10">
      <label htmlFor="comment" className="font-medium">
        发表评论
      </label>
      <textarea
        id="comment"
        name="comment"
        rows={5}
        className="mt-3 w-full resize-y border border-black/20 p-3 dark:border-white/20"
      />
      <button
        type="button"
        disabled
        className="mt-3 px-4 py-2 disabled:opacity-50"
      >
        提交评论
      </button>
    </section>
  )
}
