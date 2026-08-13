'use client'

import type {
  CommentContentTargetType,
  CommentListResponse,
} from '@/models/comment'
import { useTranslations } from 'next-intl'
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
  const t = useTranslations('Comments')
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
      setError(t('loadFailed'))
    }
    finally {
      setLoading(false)
    }
  }, [page, t, targetID, targetType])

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
    <section aria-labelledby="comments-title" className="article-comments">
      <div className="article-comments-head">
        <h2 id="comments-title">{t('title')}</h2>
        <span>{t('count', { count: result?.total ?? 0 })}</span>
      </div>

      {loading && (
        <div className="comment-state">{t('loading')}</div>
      )}

      {!loading && error && (
        <p className="comment-state comment-state-error">{error}</p>
      )}

      {!loading && !error && result?.items.length === 0 && (
        <p className="comment-state">{t('empty')}</p>
      )}

      {!loading && !error && result && result.items.length > 0 && (
        <div className="comment-list">
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
        <nav aria-label={t('pagination')} className="comment-pagination">
          <button
            type="button"
            disabled={page <= 1 || loading}
            onClick={() => setPage(current => current - 1)}
          >
            {t('previous')}
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
          >
            {t('next')}
          </button>
        </nav>
      )}

      <CommentInput
        targetType={targetType}
        targetID={targetID}
        onCreated={handleCreated}
      />
    </section>
  )
}
