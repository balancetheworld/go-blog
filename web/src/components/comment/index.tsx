'use client'

import type {
  CommentContentTargetType,
  CommentListResponse,
} from '@/models/comment'
import { useCallback, useEffect, useState } from 'react'
import { listComments } from '@/api/comment'
import { CommentInput } from './comment-input'
import { CommentItem } from './comment-item'

interface CommentSectionProps {
  targetType: CommentContentTargetType
  targetID: number
  targetAuthorID: number
}

export function CommentSection({
  targetType,
  targetID,
  targetAuthorID,
}: CommentSectionProps) {
  const [page, setPage] = useState(1)
  const [result, setResult] = useState<CommentListResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const loadComments = useCallback(async () => {
    setLoading(true)
    setError('')

    try {
      setResult(await listComments({
        targetType,
        targetId: targetID,
        page,
        pageSize: 20,
      }))
    }
    catch {
      setError('评论暂时无法加载')
    }
    finally {
      setLoading(false)
    }
  }, [page, targetID, targetType])

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

  async function handleDeleted() {
    if (result?.items.length === 1 && page > 1) {
      setPage(current => current - 1)
      return
    }

    await loadComments()
  }

  const totalPages = result
    ? Math.max(1, Math.ceil(result.total / result.pageSize))
    : 1

  return (
    <section aria-labelledby="comments-title" className="mx-auto mt-10 w-full max-w-3xl">
      <header className="border-t border-black/10 pt-8 dark:border-white/10">
        <h2 id="comments-title" className="text-xl font-semibold">评论</h2>
      </header>

      <CommentInput
        targetType={targetType}
        targetID={targetID}
        onCreated={handleCreated}
      />

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
            <CommentItem
              key={comment.id}
              comment={comment}
              targetAuthorID={targetAuthorID}
              onDeleted={handleDeleted}
            />
          ))}
        </div>
      )}

      {result && totalPages > 1 && (
        <nav aria-label="评论分页" className="flex items-center justify-between pt-5 text-sm">
          <button
            type="button"
            disabled={page <= 1 || loading}
            onClick={() => setPage(current => current - 1)}
            className="disabled:text-neutral-400"
          >
            上一页
          </button>
          <span>
            {page}
            {' / '}
            {totalPages}
          </span>
          <button
            type="button"
            disabled={page >= totalPages || loading}
            onClick={() => setPage(current => current + 1)}
            className="disabled:text-neutral-400"
          >
            下一页
          </button>
        </nav>
      )}
    </section>
  )
}
