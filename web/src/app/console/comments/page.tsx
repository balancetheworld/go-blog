import { getLocale, getTranslations } from 'next-intl/server'
import Link from 'next/link'
import { notFound } from 'next/navigation'
import { listAdminComments } from '@/api/comment.server'
import { getCurrentUser } from '@/lib/auth/current-user'
import { CommentDeleteButton } from './comment-delete-button'

interface CommentsPageProps {
  searchParams: Promise<{
    page?: string
    keyword?: string
    targetType?: 'post' | 'page' | 'diary' | 'guestbook'
  }>
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
  const locale = await getLocale()
  const t = await getTranslations('Console.adminComments')
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
        <h1 id="comments-title" className="text-2xl font-semibold">{t('title')}</h1>
        <p className="mt-1 text-sm text-neutral-500">
          {t('count', { count: result.total })}
        </p>
      </header>

      <form className="flex max-w-2xl flex-wrap gap-2">
        <input
          name="keyword"
          defaultValue={keyword}
          placeholder={t('search')}
          className="min-h-10 min-w-48 flex-1 rounded-md border border-black/15 px-3 dark:border-white/15"
        />
        <select
          name="targetType"
          defaultValue={targetType}
          className="min-h-10 rounded-md border border-black/15 px-3 dark:border-white/15"
        >
          <option value="">{t('allTypes')}</option>
          <option value="post">{t('post')}</option>
          <option value="page">{t('page')}</option>
          <option value="diary">{t('diary')}</option>
          <option value="guestbook">{t('guestbook')}</option>
        </select>
        <button
          type="submit"
          className="min-h-10 rounded-md bg-black px-4 text-sm text-white dark:bg-white dark:text-black"
        >
          {t('query')}
        </button>
      </form>

      <div className="overflow-x-auto border-y border-black/10 dark:border-white/10">
        <table className="w-full min-w-[1040px] text-left text-sm">
          <thead className="text-neutral-500">
            <tr>
              <th className="px-3 py-3">{t('content')}</th>
              <th className="px-3 py-3">{t('user')}</th>
              <th className="px-3 py-3">{t('target')}</th>
              <th className="px-3 py-3">{t('depth')}</th>
              <th className="px-3 py-3">{t('time')}</th>
              <th className="px-3 py-3">{t('moderationStatus')}</th>
              <th className="px-3 py-3">{t('moderationResult')}</th>
              <th className="px-3 py-3">{t('actions')}</th>
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
                      {t('reply')}
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
                          {t(comment.targetType)}
                          {' #'}
                          {comment.targetId}
                        </span>
                      )
                    : comment.targetType === 'guestbook'
                      ? (
                          <Link href="/#guestbook">
                            {t(comment.targetType)}
                          </Link>
                        )
                      : (
                          <Link href={`/p/${comment.targetId}`}>
                            {t(comment.targetType)}
                            {' #'}
                            {comment.targetId}
                          </Link>
                        )}
                </td>
                <td className="px-3 py-4">{comment.depth}</td>
                <td className="px-3 py-4">
                  {new Date(comment.createdAt).toLocaleString(locale)}
                </td>
                <td className="px-3 py-4">
                  {t(comment.moderationStatus)}
                </td>
                <td className="max-w-sm px-3 py-4">
                  <p className="whitespace-pre-wrap break-words">
                    {comment.moderationReason || '-'}
                  </p>
                  <p className="mt-1 text-xs text-neutral-500">
                    {t('categories')}
                    {': '}
                    {comment.moderationCategories || '-'}
                    {' · '}
                    {t('confidence')}
                    {': '}
                    {comment.moderationConfidence === undefined
                      ? '-'
                      : `${Math.round(comment.moderationConfidence * 100)}%`}
                  </p>
                </td>
                <td className="px-3 py-4">
                  <CommentDeleteButton commentID={comment.id} />
                </td>
              </tr>
            ))}

            {result.items.length === 0 && (
              <tr>
                <td colSpan={8} className="px-3 py-12 text-center text-neutral-500">
                  {t('empty')}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <nav className="flex justify-end gap-4 text-sm" aria-label={t('pagination')}>
        {page > 1
          ? <Link href={createPageHref(page - 1, keyword, targetType)}>{t('previous')}</Link>
          : <span className="text-neutral-400">{t('previous')}</span>}

        <span>
          {result.page}
          {' / '}
          {totalPages}
        </span>

        {page < totalPages
          ? <Link href={createPageHref(page + 1, keyword, targetType)}>{t('next')}</Link>
          : <span className="text-neutral-400">{t('next')}</span>}
      </nav>
    </section>
  )
}
