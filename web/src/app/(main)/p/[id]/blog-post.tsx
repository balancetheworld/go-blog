'use client'

import type { ReactNode } from 'react'
import type { Post } from '@/models/post'
import { Heart } from 'lucide-react'
import { useLocale, useTranslations } from 'next-intl'
import { useRouter } from 'next/navigation'
import { useEffect, useRef, useState } from 'react'
import { toast } from 'sonner'
import { togglePostLike } from '@/api/post'
import { ArticleLikeRive } from '@/components/design/article-like-rive'
import { useAuth } from '@/contexts/auth-context'
import { sanitizePostHTML } from '@/lib/post-html'

interface BlogPostProps {
  post: Post
  children?: ReactNode
}

export function BlogPost({ post, children }: BlogPostProps) {
  const locale = useLocale()
  const t = useTranslations('Post')
  const dateFormatter = new Intl.DateTimeFormat(locale, { dateStyle: 'long' })
  const router = useRouter()
  const { currentUser } = useAuth()
  const [liked, setLiked] = useState(post.liked)
  const [likeCount, setLikeCount] = useState(post.likeCount)
  const [liking, setLiking] = useState(false)
  const [likeAnimationKey, setLikeAnimationKey] = useState(0)
  const [likeAnimating, setLikeAnimating] = useState(false)
  const likeAnimationTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const likeRequestRef = useRef(false)
  const authorName
    = post.author.nickname || post.author.username
  const content = sanitizePostHTML(post.content)

  useEffect(() => {
    return () => {
      if (likeAnimationTimerRef.current)
        clearTimeout(likeAnimationTimerRef.current)
    }
  }, [])

  async function handleLike() {
    if (!currentUser) {
      router.push('/auth/login')
      return
    }
    if (likeRequestRef.current)
      return

    likeRequestRef.current = true
    setLiking(true)
    if (liked) {
      setLikeAnimating(false)
      if (likeAnimationTimerRef.current) {
        clearTimeout(likeAnimationTimerRef.current)
        likeAnimationTimerRef.current = null
      }
    }
    try {
      const result = await togglePostLike(post.id)
      setLiked(result.liked)
      setLikeCount(result.likeCount)
      if (result.liked) {
        setLikeAnimationKey(value => value + 1)
        setLikeAnimating(true)
        if (likeAnimationTimerRef.current)
          clearTimeout(likeAnimationTimerRef.current)
        likeAnimationTimerRef.current = setTimeout(setLikeAnimating, 1400, false)
      }
    }
    catch {
      toast.error(t('likeFailed'))
    }
    finally {
      likeRequestRef.current = false
      setLiking(false)
    }
  }

  function handleBack() {
    const value = window.sessionStorage.getItem('article-list-return')
    if (value) {
      try {
        const saved = JSON.parse(value) as {
          href?: string
          scrollY?: number
          targetPath?: string
          savedAt?: number
        }
        const valid = saved.targetPath === window.location.pathname
          && typeof saved.scrollY === 'number'
          && typeof saved.savedAt === 'number'
          && Date.now() - saved.savedAt < 60 * 60 * 1000

        if (valid) {
          window.sessionStorage.setItem('article-list-restore-scroll', String(saved.scrollY))
          window.sessionStorage.removeItem('article-list-return')
          router.back()
          return
        }
      }
      catch {
        window.sessionStorage.removeItem('article-list-return')
      }
    }

    router.push('/#articles')
  }

  return (
    <article className="article-detail-shell">
      <button type="button" className="article-back-link" onClick={handleBack}>
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
          <path d="M19 12H5" />
          <polyline points="12 19 5 12 12 5" />
        </svg>
        {t('back')}
      </button>

      <div className="article-detail-content">
        <header className="article-detail-hero">
          <div className="article-detail-copy">
            <h1>{post.title}</h1>
            <div className="article-detail-meta">
              <span>{authorName}</span>
              <time dateTime={post.publishedAt ?? post.createdAt}>
                {dateFormatter.format(new Date(post.publishedAt ?? post.createdAt))}
              </time>
            </div>
          </div>
        </header>

        <div
          className="article-prose"
          dangerouslySetInnerHTML={{ __html: content }}
        />
      </div>

      <footer className="article-detail-footer">
        <div className="article-detail-actions">
          <div className="article-detail-stat">
            <span>{t('views')}</span>
            <strong>{post.viewCount}</strong>
          </div>
          <button
            type="button"
            className={`article-like-btn${liked ? ' is-liked' : ''}${likeAnimating ? ' is-popping' : ''}`}
            aria-pressed={liked}
            aria-disabled={liking}
            aria-label={liked ? t('unlikeAction') : t('likeAction')}
            onClick={() => void handleLike()}
          >
            <Heart size={17} fill={liked ? 'currentColor' : 'none'} aria-hidden="true" />
            <span>{t('like')}</span>
            <strong>{likeCount}</strong>
            <ArticleLikeRive playKey={likeAnimationKey} />
          </button>
          <div className="article-detail-stat">
            <span>{t('comments')}</span>
            <strong>{post.commentCount}</strong>
          </div>
          {post.labels.map(label => (
            <div key={label.id} className="article-detail-stat"><span>{label.name}</span></div>
          ))}
        </div>
        {children}
      </footer>
    </article>
  )
}
