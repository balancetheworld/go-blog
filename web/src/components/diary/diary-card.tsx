import type { Diary } from '@/models/diary'
import { useLocale, useTranslations } from 'next-intl'
import Link from 'next/link'

interface DiaryCardProps {
  diary: Diary
}

export function DiaryCard({ diary }: DiaryCardProps) {
  const locale = useLocale()
  const t = useTranslations('Diary')
  const dateFormatter = new Intl.DateTimeFormat(locale, { dateStyle: 'medium' })
  const href = `/diary/${diary.slug || diary.id}`

  return (
    <Link href={href} className="diary-card diary-card-daily" data-tilt>
      <div className="diary-card-body">
        <span className="diary-badge">{diary.folder?.name || t('fallback')}</span>
        <time className="diary-date" dateTime={diary.publishedAt ?? diary.createdAt}>
          {dateFormatter.format(new Date(diary.publishedAt ?? diary.createdAt))}
        </time>
        <h2 className="diary-title">{diary.title || t('untitled')}</h2>
        {diary.description && <p className="diary-excerpt">{diary.description}</p>}
        <span className="diary-card-stats">
          {diary.viewCount}
          {' '}
          {t('views')}
          {' · '}
          {' '}
          {diary.commentCount}
          {' '}
          {t('comments')}
        </span>
      </div>
    </Link>
  )
}
