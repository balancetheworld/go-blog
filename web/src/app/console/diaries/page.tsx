import Link from 'next/link'
import { notFound } from 'next/navigation'
import { listDiaries, listDiaryFolders } from '@/api/diary.server'
import { getCurrentUser } from '@/lib/auth/current-user'
import { DiaryFolderManage } from './diary-folder-manage'
import { DiaryManage } from './diary-manage'

interface DiariesPageProps {
  searchParams: Promise<{
    page?: string
  }>
}

export default async function DiariesPage({
  searchParams,
}: DiariesPageProps) {
  const user = await getCurrentUser()
  if (!user || (user.role !== 'admin' && user.role !== 'editor'))
    notFound()

  const params = await searchParams
  const requestedPage = Number(params.page)
  const page = Number.isInteger(requestedPage) && requestedPage > 0
    ? requestedPage
    : 1
  const [result, folders] = await Promise.all([
    listDiaries({ page, pageSize: 20, status: 'all' }),
    listDiaryFolders(true),
  ])
  const totalPages = Math.max(1, Math.ceil(result.total / result.pageSize))

  return (
    <section aria-labelledby="diaries-title" className="space-y-6">
      <DiaryManage result={result} folders={folders} />
      <DiaryFolderManage folders={folders} />

      <nav aria-label="日记分页" className="flex justify-end gap-4 text-sm">
        {page > 1
          ? <Link href={`/console/diaries?page=${page - 1}`}>上一页</Link>
          : <span className="text-neutral-400">上一页</span>}

        <span>
          {result.page}
          {' '}
          /
          {' '}
          {totalPages}
        </span>

        {page < totalPages
          ? <Link href={`/console/diaries?page=${page + 1}`}>下一页</Link>
          : <span className="text-neutral-400">下一页</span>}
      </nav>
    </section>
  )
}
