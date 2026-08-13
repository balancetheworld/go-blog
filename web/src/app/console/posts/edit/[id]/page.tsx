import { notFound } from 'next/navigation'
import {
  getPost,
  listCategories,
  listLabels,
  listVisibleRoles,
  PostServerError,
} from '@/api/post.server'
import { getCurrentUser } from '@/lib/auth/current-user'
import { PostEditor } from '../../post-editor'

interface EditPostPageProps {
  // 路由动态参数不是直接 { id: string }，而是被 Promise 包裹
  params: Promise<{
    id: string
  }>
}

export default async function EditPostPage({
  params,
}: EditPostPageProps) {
  const { id } = await params
  const postID = Number(id)

  if (!Number.isInteger(postID) || postID <= 0)
    notFound()

  try {
    const [post, categories, labels, roleOptions, currentUser] = await Promise.all([
      getPost(String(postID)),
      listCategories(),
      listLabels(),
      listVisibleRoles(),
      getCurrentUser(),
    ])

    if (!currentUser)
      notFound()

    if (currentUser.role !== 'admin') {
      notFound()
    }
    return (
      <PostEditor
        post={post}
        currentUserID={currentUser.id}
        categories={categories}
        labels={labels}
        roleOptions={roleOptions}
      />
    )
  }
  catch (error) {
    if (
      error instanceof PostServerError
      && [401, 403, 404].includes(error.status)
    ) {
      notFound()
    }

    throw error
  }
}
