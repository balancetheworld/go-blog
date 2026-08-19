'use client'

import type {
  Comment,
  CommentContentTargetType,
  CommentListResponse,
} from '@/models/comment'
import { useTranslations } from 'next-intl'
import { useCallback, useEffect, useRef, useState } from 'react'
import { getCommentModeration, listComments } from '@/api/comment'
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
  const [pendingComments, setPendingComments] = useState<Comment[]>([])
  const pendingCommentsRef = useRef<Comment[]>([])
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

  useEffect(() => {
    if (pendingComments.length === 0)
      return

    let cancelled = false

    async function pollModeration() {
      const updates = await Promise.all(
        pendingCommentsRef.current.map(async (comment) => {
          if (comment.moderationStatus !== 'pending' && comment.moderationStatus !== 'manual_review')
            return comment

          try {
            const moderation = await getCommentModeration(comment.id)
            return {
              ...comment,
              moderationStatus: moderation.moderationStatus,
              moderationReason: moderation.moderationReason,
              moderatedAt: moderation.moderatedAt,
            }
          }
          catch {
            return comment
          }
        }),
      )

      if (cancelled)
        return

      const approved = updates.some(comment => comment.moderationStatus === 'approved')
      const nextPendingComments = updates.filter(
        comment => comment.moderationStatus !== 'approved',
      )
      pendingCommentsRef.current = nextPendingComments
      setPendingComments(nextPendingComments)
      if (approved)
        await loadComments()
    }

    void pollModeration()
    const timer = window.setInterval(() => {
      void pollModeration()
    }, 2000)

    return () => {
      cancelled = true
      window.clearInterval(timer)
    }
  }, [loadComments, pendingComments.length])

  async function handleCreated(comment: Comment) {
    setPendingComments((current) => {
      const next = [comment, ...current]
      pendingCommentsRef.current = next
      return next
    })
    if (page !== 1)
      setPage(1)
  }

  async function handleDeleted() {
    if (result?.items.length === 1 && page > 1) {
      setPage(current => current - 1)
      return
    }

    await loadComments()
  }

  const visibleComments = [
    ...pendingComments,
    ...(result?.items ?? []),
  ]

  const totalPages = result
    ? Math.max(1, Math.ceil(result.total / result.pageSize))
    : 1

  return (
    <section aria-labelledby="comments-title" className="article-comments">
      <div className="article-comments-head">
        <h2 id="comments-title">{t('title')}</h2>
        <span>{t('count', { count: visibleComments.length })}</span>
      </div>

      {loading && (
        <div className="comment-state">{t('loading')}</div>
      )}

      {!loading && error && (
        <p className="comment-state comment-state-error">{error}</p>
      )}

      {!loading && !error && visibleComments.length === 0 && (
        <p className="comment-state">{t('empty')}</p>
      )}

      {!loading && !error && visibleComments.length > 0 && (
        <div className="comment-list">
          {visibleComments.map(comment => (
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
