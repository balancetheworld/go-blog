import type { LucideIcon } from 'lucide-react'
import {
  ArrowRight,
  FilePenLine,
  FileText,
  MessageSquare,
  Plus,
  Users,
} from 'lucide-react'
import { useTranslations } from 'next-intl'
import Link from 'next/link'

interface DashbardProps {
  postCount: number
  commentCount: number
  userCount: number
  draftCount: number
}
interface Statistic {
  label: 'postCount' | 'commentCount' | 'userCount' | 'draftCount'
  value: number
  href: string
  icon: LucideIcon
}

export function Dashboard({
  postCount,
  commentCount,
  userCount,
  draftCount,
}: DashbardProps) {
  const t = useTranslations('Console.dashboard')

  const statistics: Statistic[] = [
    {
      label: 'postCount',
      value: postCount,
      href: '/console/posts',
      icon: FileText,
    },
    {
      label: 'commentCount',
      value: commentCount,
      href: '/console/comments',
      icon: MessageSquare,
    },
    {
      label: 'userCount',
      value: userCount,
      href: '/console/users',
      icon: Users,
    },
    {
      label: 'draftCount',
      value: draftCount,
      href: '/console/posts/drafts',
      icon: FilePenLine,
    },
  ]

  return (
    <div className="console-dashboard">
      <section aria-labelledby="dashboard-title" className="console-dashboard-overview">
        <div className="console-page-heading">
          <div>
            <span>{t('eyebrow')}</span>
            <h1 id="dashboard-title">{t('title')}</h1>
            <p>{t('description')}</p>
          </div>
          <Link href="/console/posts/new" className="console-primary-action">
            <Plus aria-hidden="true" />
            <span>{t('newPost')}</span>
          </Link>
        </div>

        <div className="console-stat-grid">
          {statistics.map(({
            icon: Icon,
            label,
            value,
            href,
          }) => (
            <Link
              key={label}
              href={href}
              className="console-stat-card"
            >
              <div>
                <span>{t(label)}</span>
                <Icon aria-hidden="true" />
              </div>
              <p>{value}</p>
            </Link>
          ))}
        </div>
      </section>

      <div className="console-dashboard-grid">
        <section aria-labelledby="recent-posts-title" className="console-dashboard-section">
          <div className="console-section-heading">
            <h2 id="recent-posts-title">{t('recentPosts')}</h2>
            <Link href="/console/posts">
              {t('viewAll')}
              <ArrowRight aria-hidden="true" />
            </Link>
          </div>
          <p className="console-empty-state">{t('noPosts')}</p>
        </section>

        <section aria-labelledby="recent-comments-title" className="console-dashboard-section">
          <div className="console-section-heading">
            <h2 id="recent-comments-title">{t('recentComments')}</h2>
            <Link href="/console/comments">
              {t('viewAll')}
              <ArrowRight aria-hidden="true" />
            </Link>
          </div>
          <p className="console-empty-state">{t('noComments')}</p>
        </section>
      </div>
    </div>
  )
}
