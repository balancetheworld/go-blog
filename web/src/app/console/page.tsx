import { listAdminComments } from '@/api/comment.server'
import { listPosts } from '@/api/post.server'
import { listUsers } from '@/api/user.server'
import { Dashboard } from './dashboard'

export default async function ConsolePage() {
  const [posts, comments, users, drafts] = await Promise.all([
    listPosts({
      page: 1,
      pageSize: 1,
      status: 'all',
    }),
    listAdminComments({
      page: 1,
      pageSize: 1,
    }),
    listUsers({
      page: 1,
      pageSize: 1,
    }),
    listPosts({
      page: 1,
      pageSize: 1,
      status: 'draft',
    }),
  ])
  return (
    <Dashboard
      postCount={posts.total}
      commentCount={comments.total}
      userCount={users.total}
      draftCount={drafts.total}
    />
  )
}
