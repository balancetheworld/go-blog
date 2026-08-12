import { listCategories, listLabels, listVisibleRoles } from '@/api/post.server'
import { PostEditor } from '../post-editor'

export default async function NewPostPage() {
  const [categories, labels, roleOptions] = await Promise.all([
    listCategories(),
    listLabels(),
    listVisibleRoles(),
  ])

  return (
    <PostEditor
      categories={categories}
      labels={labels}
      roleOptions={roleOptions}
    />
  )
}
