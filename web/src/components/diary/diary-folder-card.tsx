import type { DiaryFolder } from '@/models/diary'
import Link from 'next/link'

interface DiaryFolderCardProps {
  folder: DiaryFolder
  active?: boolean
}

export function DiaryFolderCard({
  folder,
  active = false,
}: DiaryFolderCardProps) {
  return (
    <Link
      href={`/diary?folder=${folder.id}`}
      aria-current={active ? 'page' : undefined}
      className={`block rounded-md border p-4 ${active ? 'border-black dark:border-white' : 'border-black/10 dark:border-white/10'}`}
    >
      <p className="font-medium">{folder.name}</p>
      {folder.description && (
        <p className="mt-1 line-clamp-2 text-sm leading-6 text-neutral-500">
          {folder.description}
        </p>
      )}
    </Link>
  )
}
