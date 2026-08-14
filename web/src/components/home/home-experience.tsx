'use client'

import type { Diary } from '@/models/diary'
import type { Moment } from '@/models/moment'
import type { PostListItem } from '@/models/post'
import { useLocale, useTranslations } from 'next-intl'
import Link from 'next/link'
import { Guestbook } from './guestbook'

interface HomeExperienceProps {
  posts: PostListItem[]
  diaries: Diary[]
  moments: Moment[]
}

interface TimelineEntry {
  id: string
  type: 'article' | 'diary' | 'moment'
  label: string
  title: string
  href: string
  date: string
}

export function HomeExperience({ posts, diaries, moments = [] }: HomeExperienceProps) {
  const locale = useLocale()
  const momentsT = useTranslations('Home.moments')
  const timelineT = useTranslations('Home.timeline')
  const guestbookT = useTranslations('Home.guestbook')
  const diaryT = useTranslations('Diary')
  const momentDateFormatter = new Intl.DateTimeFormat(locale, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
  const timelineDateFormatter = new Intl.DateTimeFormat(locale, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
  const timelineEntries: TimelineEntry[] = [
    ...posts.map(post => ({
      id: `post-${post.id}`,
      type: 'article' as const,
      label: timelineT('article'),
      title: post.title,
      href: `/p/${post.slug || post.id}`,
      date: post.publishedAt ?? post.createdAt,
    })),
    ...diaries.map(diary => ({
      id: `diary-${diary.id}`,
      type: 'diary' as const,
      label: timelineT('diary'),
      title: diary.title || diaryT('untitled'),
      href: `/diary/${diary.slug || diary.id}`,
      date: diary.publishedAt ?? diary.createdAt,
    })),
    ...moments.map(moment => ({
      id: `moment-${moment.id}`,
      type: 'moment' as const,
      label: timelineT('moment'),
      title: moment.content,
      href: '/#moments',
      date: moment.createdAt,
    })),
  ]
    .sort((left, right) => new Date(right.date).getTime() - new Date(left.date).getTime())
    .slice(0, 6)

  return (
    <>
      <section className="section moments" id="moments">
        <div className="section-header">
          <h2 className="section-title">
            <span className="title-num">03</span>
            <span className="title-text">{momentsT('title')}</span>
          </h2>
          <p className="section-sub">{momentsT('subtitle')}</p>
        </div>
        <div className="moments-feed">
          {moments.map(moment => (
            <div key={moment.id} className="moment-card">
              <div className="moment-time">
                <span className="time-dot" />
                <time dateTime={moment.createdAt}>
                  {momentDateFormatter.format(new Date(moment.createdAt))}
                </time>
              </div>
              <p className="moment-text">{moment.content}</p>
            </div>
          ))}
          {moments.length === 0 && (
            <p className="moment-text">{momentsT('empty')}</p>
          )}
        </div>
      </section>

      <section className="section timeline" id="timeline">
        <div className="section-header">
          <h2 className="section-title">
            <span className="title-num">04</span>
            <span className="title-text">{timelineT('title')}</span>
          </h2>
        </div>
        <div className="timeline-list">
          {timelineEntries.map(entry => (
            <article key={entry.id} className="timeline-item" data-type={entry.type}>
              <time className="timeline-date" dateTime={entry.date}>
                {timelineDateFormatter.format(new Date(entry.date)).replaceAll('/', '.')}
              </time>
              <Link href={entry.href} className="timeline-entry">
                <span className="timeline-label">{entry.label}</span>
                <h3 className="timeline-title">{entry.title}</h3>
              </Link>
            </article>
          ))}
        </div>
      </section>

      <section className="section guestbook" id="guestbook">
        <div className="section-header">
          <h2 className="section-title">
            <span className="title-num">05</span>
            <span className="title-text">{guestbookT('title')}</span>
          </h2>
          <p className="section-sub">{guestbookT('subtitle')}</p>
        </div>
        <Guestbook />
      </section>
    </>
  )
}
