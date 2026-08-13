import type { Diary } from '@/models/diary'
import { DiaryCard } from './diary-card'

interface DiaryListProps {
  diaries: Diary[]
}

export function DiaryList({ diaries }: DiaryListProps) {
  if (diaries.length === 0) {
    return (
      <p className="border-y border-black/10 py-12 text-center text-neutral-500 dark:border-white/10">
        暂无日记
      </p>
    )
  }

  return (
    <div className="border-t border-black/10 dark:border-white/10">
      {diaries.map(diary => (
        <DiaryCard key={diary.id} diary={diary} />
      ))}
    </div>
  )
}
