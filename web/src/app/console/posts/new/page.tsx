import { notFound } from 'next/navigation'
import { listCategories, listLabels, listVisibleRoles } from '@/api/post.server'
import { getCurrentUser } from '@/lib/auth/current-user'
import { PostEditor } from '../post-editor'

export default async function NewPostPage() {
  const [categories, labels, roleOptions, currentUser] = await Promise.all([
    listCategories(),
    listLabels(),
    listVisibleRoles(),
    getCurrentUser(),
  ])

  if (!currentUser || currentUser.role !== 'admin')
    notFound()

  return (
    <PostEditor
      currentUserID={currentUser.id}
      categories={categories}
      labels={labels}
      roleOptions={roleOptions}
    />
  )
}
