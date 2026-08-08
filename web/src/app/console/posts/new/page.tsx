import { listCategories, listLabels } from '@/api/post.server'
import { PostEditor } from '../post-editor'

export default async function NewPostPage() {
  const [categories, labels] = await Promise.all([
    listCategories(),
    listLabels(),
  ])

  return (
    <PostEditor
      categories={categories}
      labels={labels}
    />
  )
}
