import type { Diary } from '@/models/diary'
import { Eye, MessageCircle } from 'lucide-react'
import Link from 'next/link'

interface DiaryCardProps {
  diary: Diary
}

const dateFormatter = new Intl.DateTimeFormat('zh-CN', {
  dateStyle: 'medium',
})

export function DiaryCard({ diary }: DiaryCardProps) {
  const href = `/diary/${diary.slug || diary.id}`

  return (
    <article className="border-b border-black/10 py-6 dark:border-white/10">
      <div className="flex flex-wrap items-center gap-2 text-xs text-neutral-500">
        {diary.folder && <span>{diary.folder.name}</span>}
        <time dateTime={diary.publishedAt ?? diary.createdAt}>
          {dateFormatter.format(new Date(diary.publishedAt ?? diary.createdAt))}
        </time>
      </div>
      <h2 className="mt-2 text-xl font-semibold">
        <Link href={href}>{diary.title || '无标题日记'}</Link>
      </h2>
      {diary.description && (
        <p className="mt-2 line-clamp-2 leading-7 text-neutral-600 dark:text-neutral-400">
          {diary.description}
        </p>
      )}
      <div className="mt-3 flex items-center gap-4 text-xs text-neutral-500">
        <span className="inline-flex items-center gap-1">
          <Eye className="size-4" aria-hidden="true" />
          {diary.viewCount}
        </span>
        <span className="inline-flex items-center gap-1">
          <MessageCircle className="size-4" aria-hidden="true" />
          {diary.commentCount}
        </span>
      </div>
    </article>
  )
}
