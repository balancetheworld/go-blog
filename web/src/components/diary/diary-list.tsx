import type { Diary } from '@/models/diary'
import { useTranslations } from 'next-intl'
import { DiaryCard } from './diary-card'

interface DiaryListProps {
  diaries: Diary[]
}

export function DiaryList({ diaries }: DiaryListProps) {
  const t = useTranslations('Diary')

  if (diaries.length === 0) {
    return (
      <p className="article-empty visible">
        {t('empty')}
      </p>
    )
  }

  return (
    <div className="diary-grid">
      {diaries.map(diary => (
        <DiaryCard key={diary.id} diary={diary} />
      ))}
    </div>
  )
}
