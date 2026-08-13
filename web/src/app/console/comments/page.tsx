import Link from 'next/link'
import { notFound } from 'next/navigation'
import { listAdminComments } from '@/api/comment.server'
import { getCurrentUser } from '@/lib/auth/current-user'
import { CommentDeleteButton } from './comment-delete-button'

interface CommentsPageProps {
  searchParams: Promise<{
    page?: string
    keyword?: string
    targetType?: 'post' | 'page' | 'diary'
  }>
}

const targetNames = {
  post: '文章',
  page: '页面',
  diary: '日记',
}

function createPageHref(
  page: number,
  keyword: string,
  targetType: string,
) {
  const params = new URLSearchParams()
  params.set('page', String(page))
  if (keyword)
    params.set('keyword', keyword)
  if (targetType)
    params.set('targetType', targetType)

  return `/console/comments?${params.toString()}`
}

export default async function CommentsPage({
  searchParams,
}: CommentsPageProps) {
  const user = await getCurrentUser()
  if (!user || user.role !== 'admin')
    notFound()

  const params = await searchParams
  const requestedPage = Number(params.page)
  const page = Number.isInteger(requestedPage) && requestedPage > 0 ? requestedPage : 1
  const keyword = params.keyword?.trim() ?? ''
  const targetType = params.targetType ?? ''
  const result = await listAdminComments({
    page,
    pageSize: 20,
    keyword,
    targetType: targetType || undefined,
  })
  const totalPages = Math.max(1, Math.ceil(result.total / result.pageSize))

  return (
    <section aria-labelledby="comments-title" className="space-y-6">
      <header>
        <h1 id="comments-title" className="text-2xl font-semibold">评论管理</h1>
        <p className="mt-1 text-sm text-neutral-500">
          共
          {' '}
          {result.total}
          {' '}
          条评论
        </p>
      </header>

      <form className="flex max-w-2xl flex-wrap gap-2">
        <input
          name="keyword"
          defaultValue={keyword}
          placeholder="搜索评论内容"
          className="min-h-10 min-w-48 flex-1 rounded-md border border-black/15 px-3 dark:border-white/15"
        />
        <select
          name="targetType"
          defaultValue={targetType}
          className="min-h-10 rounded-md border border-black/15 px-3 dark:border-white/15"
        >
          <option value="">全部类型</option>
          <option value="post">文章</option>
          <option value="page">页面</option>
          <option value="diary">日记</option>
        </select>
        <button
          type="submit"
          className="min-h-10 rounded-md bg-black px-4 text-sm text-white dark:bg-white dark:text-black"
        >
          查询
        </button>
      </form>

      <div className="overflow-x-auto border-y border-black/10 dark:border-white/10">
        <table className="w-full min-w-[880px] text-left text-sm">
          <thead className="text-neutral-500">
            <tr>
              <th className="px-3 py-3">评论内容</th>
              <th className="px-3 py-3">评论用户</th>
              <th className="px-3 py-3">目标</th>
              <th className="px-3 py-3">回复层级</th>
              <th className="px-3 py-3">时间</th>
              <th className="px-3 py-3">操作</th>
            </tr>
          </thead>
          <tbody>
            {result.items.map(comment => (
              <tr key={comment.id} className="border-t border-black/10 dark:border-white/10">
                <td className="max-w-md px-3 py-4">
                  <p className="line-clamp-3 whitespace-pre-wrap break-words">
                    {comment.content}
                  </p>
                  {comment.replyToUser && (
                    <p className="mt-1 text-xs text-neutral-500">
                      回复
                      {' '}
                      {comment.replyToUser.nickname || comment.replyToUser.username}
                    </p>
                  )}
                </td>
                <td className="px-3 py-4">
                  {comment.author.nickname || comment.author.username}
                </td>
                <td className="px-3 py-4">
                  {comment.targetType === 'diary'
                    ? (
                        <span>
                          {targetNames[comment.targetType]}
                          {' #'}
                          {comment.targetId}
                        </span>
                      )
                    : (
                        <Link href={`/p/${comment.targetId}`}>
                          {targetNames[comment.targetType]}
                          {' #'}
                          {comment.targetId}
                        </Link>
                      )}
                </td>
                <td className="px-3 py-4">{comment.depth}</td>
                <td className="px-3 py-4">
                  {new Date(comment.createdAt).toLocaleString('zh-CN')}
                </td>
                <td className="px-3 py-4">
                  <CommentDeleteButton commentID={comment.id} />
                </td>
              </tr>
            ))}

            {result.items.length === 0 && (
              <tr>
                <td colSpan={6} className="px-3 py-12 text-center text-neutral-500">
                  暂无评论
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <nav className="flex justify-end gap-4 text-sm" aria-label="评论分页">
        {page > 1
          ? <Link href={createPageHref(page - 1, keyword, targetType)}>上一页</Link>
          : <span className="text-neutral-400">上一页</span>}

        <span>
          {result.page}
          {' / '}
          {totalPages}
        </span>

        {page < totalPages
          ? <Link href={createPageHref(page + 1, keyword, targetType)}>下一页</Link>
          : <span className="text-neutral-400">下一页</span>}
      </nav>
    </section>
  )
}
