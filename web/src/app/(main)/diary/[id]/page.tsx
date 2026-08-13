import type { Metadata } from 'next'
import { getLocale, getTranslations } from 'next-intl/server'
import Image from 'next/image'
import Link from 'next/link'
import { notFound } from 'next/navigation'
import { DiaryServerError, getDiary } from '@/api/diary.server'
import { CommentSection } from '@/components/comment'
import { sanitizePostHTML } from '@/lib/post-html'

interface DiaryDetailPageProps {
  params: Promise<{
    id: string
  }>
}

export async function generateMetadata({
  params,
}: DiaryDetailPageProps): Promise<Metadata> {
  const { id } = await params
  const t = await getTranslations('Diary')

  try {
    const diary = await getDiary(id)
    return {
      title: diary.title || t('fallback'),
      description: diary.description,
    }
  }
  catch {
    return {}
  }
}

export default async function DiaryDetailPage({
  params,
}: DiaryDetailPageProps) {
  const { id } = await params
  const locale = await getLocale()
  const t = await getTranslations('Diary')
  const dateFormatter = new Intl.DateTimeFormat(locale, { dateStyle: 'long' })

  let diary
  try {
    diary = await getDiary(id)
  }
  catch (error) {
    if (
      error instanceof DiaryServerError
      && [401, 403, 404].includes(error.status)
    ) {
      notFound()
    }

    throw error
  }

  const authorName = diary.author.nickname || diary.author.username
  const content = sanitizePostHTML(diary.content)

  return (
    <main className="article-detail-main">
      <article className="article-detail-shell">
        <Link href="/diary" className="article-back-link">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
            <path d="M19 12H5" />
            <polyline points="12 19 5 12 12 5" />
          </svg>
          {t('back')}
        </Link>

        <header className="article-detail-hero">
          <div className="article-detail-copy">
            <div className="article-detail-kicker">{diary.folder?.name || t('title')}</div>
            <h1>{diary.title || t('untitled')}</h1>
            {diary.description && <p>{diary.description}</p>}
            <div className="article-detail-meta">
              <time dateTime={diary.publishedAt ?? diary.createdAt}>
                {dateFormatter.format(new Date(diary.publishedAt ?? diary.createdAt))}
              </time>
              <span>{authorName}</span>
              {diary.visibility === 'private' && <span>{t('private')}</span>}
              {diary.visibility === 'roles' && (
                <span>{t('rolesVisible', { roles: diary.visibleRoles.map(role => role.name).join('、') })}</span>
              )}
            </div>
          </div>
          {diary.cover && (
            <div className="article-detail-cover">
              <Image
                src={diary.cover}
                alt={diary.title || t('cover')}
                fill
                priority
                sizes="(max-width: 768px) 100vw, 768px"
                className="object-cover"
              />
            </div>
          )}
        </header>

        <div
          className="article-prose"
          dangerouslySetInnerHTML={{ __html: content }}
        />

        <footer className="article-detail-footer">
          <div className="article-detail-actions">
            <div className="article-detail-stat">
              <span>{t('viewCount')}</span>
              <strong>{diary.viewCount}</strong>
            </div>
            <div className="article-detail-stat">
              <span>{t('likes')}</span>
              <strong>{diary.likeCount}</strong>
            </div>
            <div className="article-detail-stat">
              <span>{t('comments')}</span>
              <strong>{diary.commentCount}</strong>
            </div>
          </div>
          <CommentSection
            targetType="diary"
            targetID={diary.id}
            targetAuthorID={diary.author.id}
          />
        </footer>
      </article>
    </main>
  )
}
