'use client'

import type { FormEvent } from 'react'
import type { Comment } from '@/models/comment'
import { gsap } from 'gsap'
import { useLocale, useTranslations } from 'next-intl'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useCallback, useEffect, useRef, useState } from 'react'
import { toast } from 'sonner'
import { createComment, listCommentReplies, listComments } from '@/api/comment'
import { useAuth } from '@/contexts/auth-context'

const pageSize = 4
export function Guestbook() {
  const locale = useLocale()
  const t = useTranslations('Home.guestbook')
  const dateFormatter = new Intl.DateTimeFormat(locale, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
  const pathname = usePathname()
  const { currentUser, isLoading: authLoading } = useAuth()
  const [messages, setMessages] = useState<Comment[]>([])
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [content, setContent] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [flippedMessage, setFlippedMessage] = useState<number | null>(null)
  const [replies, setReplies] = useState<Record<number, Comment[]>>({})
  const pageRef = useRef<HTMLDivElement>(null)
  const pageCount = Math.max(1, Math.ceil(total / pageSize))

  const loadMessages = useCallback(async (nextPage: number, direction = 0) => {
    setLoading(true)
    const pageElement = pageRef.current
    try {
      if (direction !== 0 && pageElement) {
        await gsap.to(pageElement, {
          x: direction > 0 ? -72 : 72,
          opacity: 0,
          duration: 0.24,
          ease: 'power2.in',
        })
      }

      const result = await listComments({
        targetType: 'guestbook',
        targetId: 1,
        page: nextPage,
        pageSize,
      })
      setMessages(result.items)
      setTotal(result.total)
      setPage(result.page)
      setFlippedMessage(null)

      if (direction !== 0 && pageElement) {
        await new Promise<void>(resolve => requestAnimationFrame(() => resolve()))
        gsap.fromTo(pageElement, {
          x: direction > 0 ? 72 : -72,
          opacity: 0,
        }, {
          x: 0,
          opacity: 1,
          duration: 0.46,
          ease: 'expo.out',
        })
      }
    }
    catch {
      if (pageElement) {
        gsap.to(pageElement, {
          x: 0,
          opacity: 1,
          duration: 0.2,
        })
      }
      toast.error(t('loadFailed'))
    }
    finally {
      setLoading(false)
    }
  }, [t])

  useEffect(() => {
    void loadMessages(1)
    const pageElement = pageRef.current
    return () => {
      if (pageElement)
        gsap.killTweensOf(pageElement)
    }
  }, [loadMessages])

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const value = content.trim()
    if (!value || submitting)
      return

    setSubmitting(true)
    try {
      await createComment({
        targetType: 'guestbook',
        targetId: 1,
        content: value,
      })
      setContent('')
      await loadMessages(1)
      toast.success(t('published'))
    }
    catch {
      toast.error(t('publishFailed'))
    }
    finally {
      setSubmitting(false)
    }
  }

  async function toggleReply(message: Comment) {
    if (flippedMessage === message.id) {
      setFlippedMessage(null)
      return
    }

    setFlippedMessage(message.id)
    if (replies[message.id] || message.replyCount === 0)
      return

    try {
      const items = await listCommentReplies(message.id)
      setReplies(current => ({ ...current, [message.id]: items }))
    }
    catch {
      toast.error(t('replyFailed'))
    }
  }

  return (
    <div className="guestbook-carousel">
      <div className="guestbook-compose">
        {!authLoading && currentUser && (
          <form className="comment-form" onSubmit={handleSubmit}>
            <textarea
              aria-label={t('write')}
              placeholder={t('placeholder')}
              rows={3}
              maxLength={2000}
              required
              value={content}
              onChange={event => setContent(event.target.value)}
            />
            <button type="submit" disabled={submitting || !content.trim()}>
              {submitting ? t('publishing') : t('publish')}
            </button>
          </form>
        )}
        {!authLoading && !currentUser && (
          <div className="comment-login-prompt">
            <span>{t('loginPrompt')}</span>
            <Link href={`/auth/login?next=${encodeURIComponent(`${pathname}#guestbook`)}`}>
              {t('goLogin')}
            </Link>
          </div>
        )}
      </div>

      <div className="guestbook-viewport">
        <div ref={pageRef} className="guestbook-page">
          {messages.map((message) => {
            const flipped = flippedMessage === message.id
            const messageReplies = replies[message.id] ?? []
            const reply = messageReplies[0]
            const authorName = message.author.nickname || message.author.username

            return (
              <article key={message.id} className={`message-card${flipped ? ' is-flipped' : ''}`}>
                <div className="message-card-inner">
                  <div className="message-face message-front">
                    <div className="message-meta">
                      <span className="message-name">{authorName}</span>
                      <time dateTime={message.createdAt}>
                        {dateFormatter.format(new Date(message.createdAt)).replaceAll('/', '.')}
                      </time>
                    </div>
                    <p className="message-text" title={message.content}>{message.content}</p>
                  </div>
                  <div className="message-face message-reply">
                    <span className="reply-label">{t('reply')}</span>
                    <p className="message-text" title={reply?.content}>
                      {reply ? reply.content : t('noReply')}
                    </p>
                  </div>
                </div>
                <button
                  className="guestbook-toggle"
                  type="button"
                  aria-label={flipped ? t('backMessage') : t('viewReply')}
                  aria-expanded={flipped}
                  onClick={() => void toggleReply(message)}
                >
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                    <polyline points="6 9 12 15 18 9" />
                  </svg>
                </button>
              </article>
            )
          })}
        </div>
      </div>

      {!loading && messages.length === 0 && (
        <p className="guestbook-empty">{t('empty')}</p>
      )}

      <div className="guestbook-controls" aria-label={t('pagination')}>
        <button
          className="guestbook-page-btn magnetic"
          type="button"
          disabled={loading || page <= 1}
          aria-label={t('previous')}
          onClick={() => void loadMessages(page - 1, -1)}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
            <polyline points="15 18 9 12 15 6" />
          </svg>
        </button>
        <span className="guestbook-page-count" aria-live="polite">
          {page}
          {' / '}
          {pageCount}
        </span>
        <button
          className="guestbook-page-btn magnetic"
          type="button"
          disabled={loading || page >= pageCount}
          aria-label={t('next')}
          onClick={() => void loadMessages(page + 1, 1)}
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
            <polyline points="9 18 15 12 9 6" />
          </svg>
        </button>
      </div>
    </div>
  )
}
