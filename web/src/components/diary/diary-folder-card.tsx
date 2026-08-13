import type { DiaryFolder } from '@/models/diary'
import { useTranslations } from 'next-intl'
import Link from 'next/link'
import { DiaryFolderMedia } from '@/components/diary/diary-folder-media'
import { getDiaryFolderMedia, getDiaryFolderStyle, isDiaryFolderSlug } from '@/lib/diary-folders'

interface DiaryFolderCardProps {
  folder: DiaryFolder
  active?: boolean
}

export function DiaryFolderCard({
  folder,
  active = false,
}: DiaryFolderCardProps) {
  const t = useTranslations('Diary')
  const folderStyle = getDiaryFolderStyle(folder.slug)
  const folderMedia = getDiaryFolderMedia(folder.slug)
  const folderName = isDiaryFolderSlug(folder.slug)
    ? t(`folders.${folder.slug}`)
    : folder.name

  return (
    <Link
      href={`/diary?folder=${folder.id}`}
      aria-current={active ? 'page' : undefined}
      className={`diary-folder ${folderStyle}${active ? ' is-active' : ''}`}
    >
      <span className="folder-paper paper-one" />
      <span className="folder-paper paper-two" />
      <span className="folder-back" />
      <span className="folder-front">
        <DiaryFolderMedia media={folderMedia} />
      </span>
      <span className="folder-tab">{folderName}</span>
      <span className="folder-meta">{folder.description || t('view')}</span>
    </Link>
  )
}
