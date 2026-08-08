'use client'

import type { Comment, CommentListResponse } from '@/models/comment'
import { Trash2 } from 'lucide-react'
import { useCallback, useEffect, useState } from 'react'
import { toast } from 'sonner'
import { deleteComment, listComments } from '@/api/comment'
import { useAuth } from '@/contexts/auth-context'
import { CommentInput } from './comment-input'

interface CommentSectionProps {
  postID: number
  postAuthorID: number
}

const dateFormatter = new Intl.DateTimeFormat('zh-CN', {
  dateStyle: 'medium',
  timeStyle: 'short',
})

export function CommentSection({
  postID,
  postAuthorID,
}: CommentSectionProps) {
  const { currentUser, currentRole } = useAuth()
  const [page, setPage] = useState(1)
  const [result, setResult] = useState<CommentListResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [deletingID, setDeletingID] = useState<number | null>(null)
  const [error, setError] = useState('')

  const loadComments = useCallback(async () => {
    setLoading(true)
    setError('')

    try {
      const comments = await listComments({
        postId: postID,
        page,
        pageSize: 20,
      })
      setResult(comments)
    }
    catch {
      setError('评论暂时无法加载')
    }
    finally {
      setLoading(false)
    }
  }, [page, postID])

  useEffect(() => {
    void loadComments()
  }, [loadComments])

  async function handleCreated() {
    if (page === 1) {
      await loadComments()
      return
    }

    setPage(1)
  }

  async function handleDelete(comment: Comment) {
    if (deletingID !== null)
      return

    setDeletingID(comment.id)

    try {
      await deleteComment(comment.id)
      toast.success('评论已删除')

      if (result?.items.length === 1 && page > 1) {
        setPage(currentPage => currentPage - 1)
      }
      else {
        await loadComments()
      }
    }
    catch {
      toast.error('删除评论失败')
    }
    finally {
      setDeletingID(null)
    }
  }

  function canDelete(comment: Comment) {
    if (!currentUser)
      return false

    return currentRole === 'admin'
      || comment.author.id === currentUser.id
      || (currentRole === 'editor' && currentUser.id === postAuthorID)
  }

  const totalPages = result
    ? Math.max(1, Math.ceil(result.total / result.pageSize))
    : 1

  return (
    <section aria-labelledby="comments-title" className="mx-auto mt-10 w-full max-w-3xl">
      <header className="border-t border-black/10 pt-8 dark:border-white/10">
        <h2 id="comments-title" className="text-xl font-semibold">
          评论
          {result && ` (${result.total})`}
        </h2>
      </header>

      <CommentInput postID={postID} onCreated={handleCreated} />

      {loading && (
        <div className="min-h-24 py-6 text-sm text-neutral-500">正在加载评论</div>
      )}

      {!loading && error && (
        <p className="py-6 text-sm text-red-600">{error}</p>
      )}

      {!loading && !error && result?.items.length === 0 && (
        <p className="py-6 text-sm text-neutral-500">暂无评论</p>
      )}

      {!loading && !error && result && result.items.length > 0 && (
        <div className="divide-y divide-black/10 border-y border-black/10 dark:divide-white/10 dark:border-white/10">
          {result.items.map(comment => (
            <article key={comment.id} className="py-5">
              <header className="flex items-start justify-between gap-4">
                <div>
                  <p className="font-medium">
                    {comment.author.nickname || comment.author.username}
                  </p>
                  <time className="text-xs text-neutral-500" dateTime={comment.createdAt}>
                    {dateFormatter.format(new Date(comment.createdAt))}
                  </time>
                </div>

                {canDelete(comment) && (
                  <button
                    type="button"
                    aria-label="删除评论"
                    title="删除评论"
                    disabled={deletingID !== null}
                    onClick={() => void handleDelete(comment)}
                    className="inline-flex size-9 items-center justify-center disabled:opacity-50"
                  >
                    <Trash2 className="size-4" aria-hidden="true" />
                  </button>
                )}
              </header>

              <p className="mt-3 whitespace-pre-wrap break-words leading-7">
                {comment.content}
              </p>
            </article>
          ))}
        </div>
      )}

      {result && totalPages > 1 && (
        <nav aria-label="评论分页" className="flex items-center justify-between pt-5 text-sm">
          <button
            type="button"
            disabled={page <= 1 || loading}
            onClick={() => setPage(currentPage => currentPage - 1)}
            className="disabled:text-neutral-400"
          >
            上一页
          </button>
          <span>
            {page}
            {' '}
            /
            {' '}
            {totalPages}
          </span>
          <button
            type="button"
            disabled={page >= totalPages || loading}
            onClick={() => setPage(currentPage => currentPage + 1)}
            className="disabled:text-neutral-400"
          >
            下一页
          </button>
        </nav>
      )}
    </section>
  )
}
