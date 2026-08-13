import { notFound } from 'next/navigation'
import {
  DiaryServerError,
  getDiary,
  listDiaryFolders,
} from '@/api/diary.server'
import { listVisibleRoles } from '@/api/post.server'
import { getCurrentUser } from '@/lib/auth/current-user'
import { DiaryEditor } from '../../diary-editor'

interface EditDiaryPageProps {
  params: Promise<{
    id: string
  }>
}

export default async function EditDiaryPage({
  params,
}: EditDiaryPageProps) {
  const { id } = await params
  const currentUser = await getCurrentUser()
  if (!currentUser)
    notFound()

  const [folders, roleOptions] = await Promise.all([
    listDiaryFolders(true),
    listVisibleRoles(),
  ])

  if (id === 'new') {
    return (
      <DiaryEditor
        folders={folders}
        roleOptions={roleOptions}
      />
    )
  }

  const diaryID = Number(id)
  if (!Number.isInteger(diaryID) || diaryID <= 0)
    notFound()

  try {
    const diary = await getDiary(diaryID)
    if (
      currentUser.role === 'editor'
      && currentUser.id !== diary.author.id
    ) {
      notFound()
    }

    return (
      <DiaryEditor
        diary={diary}
        folders={folders}
        roleOptions={roleOptions}
      />
    )
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
}
