import Link from 'next/link'
import { notFound, redirect } from 'next/navigation'
import { listPosts } from '@/api/post.server'
import { getCurrentUser } from '@/lib/auth/current-user'
import { PostManage } from '../post-manage'

interface ConsolePostDraftsPageProps {
  searchParams: Promise<{
    page?: string
    keyword?: string
  }>
}

function createPageHref(page: number, keyword: string) {
  const params = new URLSearchParams()
  params.set('page', String(page))
  if (keyword)
    params.set('keyword', keyword)

  return `/console/posts/drafts?${params.toString()}`
}

export default async function ConsolePostDraftsPage({ searchParams }: ConsolePostDraftsPageProps) {
  const currentUser = await getCurrentUser()
  if (!currentUser)
    redirect('/auth/login?next=/console/posts/drafts')
  if (currentUser.role !== 'admin' && currentUser.role !== 'editor')
    notFound()

  const params = await searchParams
  const keyword = params.keyword?.trim() ?? ''
  const requestedPage = Number(params.page)
  const page = Number.isInteger(requestedPage) && requestedPage > 0 ? requestedPage : 1
  const result = await listPosts({
    page,
    pageSize: 20,
    keyword,
    status: 'draft',
    sort: 'latest',
  })
  const totalPages = Math.max(1, Math.ceil(result.total / result.pageSize))

  return (
    <section aria-labelledby="posts-title" className="space-y-6">
      <PostManage
        key={keyword}
        result={result}
        initialKeyword={keyword}
        title="草稿箱"
        basePath="/console/posts/drafts"
      />

      <nav aria-label="草稿分页" className="flex items-center justify-end gap-4 text-sm">
        {page > 1
          ? <Link href={createPageHref(page - 1, keyword)}>上一页</Link>
          : <span className="text-neutral-400">上一页</span>}

        <span>
          {result.page}
          {' / '}
          {totalPages}
        </span>

        {page < totalPages
          ? <Link href={createPageHref(page + 1, keyword)}>下一页</Link>
          : <span className="text-neutral-400">下一页</span>}
      </nav>
    </section>
  )
}
