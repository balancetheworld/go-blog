'use client'

import type { PostListItem } from '@/models/post'
import { Eye, Heart } from 'lucide-react'
import { useLocale, useTranslations } from 'next-intl'
import Link from 'next/link'

interface PostListProps {
  posts: PostListItem[]
}

export default function PostList({ posts }: PostListProps) {
  const locale = useLocale()
  const t = useTranslations('Home.articles')
  const dateFormatter = new Intl.DateTimeFormat(locale, { dateStyle: 'medium' })

  if (posts.length === 0) {
    return (
      <div className="article-empty visible">
        {t('empty')}
      </div>
    )
  }

  return (
    <div className="article-list">
      {posts.map((post, index) => {
        const href = `/p/${post.slug || post.id}`

        return (
          <Link
            key={post.id}
            href={href}
            className="article-card"
            data-tilt
            onClick={() => {
              window.sessionStorage.setItem('article-list-return', JSON.stringify({
                href: window.location.href,
                scrollY: window.scrollY,
                targetPath: href,
                savedAt: Date.now(),
              }))
            }}
          >
            <div className="card-index">{String(index + 1).padStart(3, '0')}</div>
            <div className="card-body">
              {(post.category || post.top) && (
                <span className="card-tag">
                  {post.top ? t('top') : ''}
                  {post.category?.name || t('fallbackCategory')}
                </span>
              )}
              <h2 className="card-title">{post.title}</h2>
              {post.description && <p className="card-excerpt">{post.description}</p>}
              <div className="card-foot">
                <div className="card-meta">
                  <time dateTime={post.publishedAt ?? post.createdAt}>
                    {dateFormatter.format(new Date(post.publishedAt ?? post.createdAt))}
                  </time>
                  <span>{post.author.nickname || post.author.username}</span>
                </div>
                <div className="card-stats" aria-label={t('stats')}>
                  <span>
                    <Eye aria-hidden="true" />
                    {post.viewCount}
                    {' '}
                    {t('views')}
                  </span>
                  <span>
                    <Heart aria-hidden="true" />
                    {post.likeCount}
                    {' '}
                    {t('likes')}
                  </span>
                </div>
              </div>
            </div>
            <div className="card-arrow">→</div>
          </Link>
        )
      })}
    </div>
  )
}
