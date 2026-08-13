import Link from 'next/link'
import {
  DiaryServerError,
  listDiaries,
  listDiaryFolders,
} from '@/api/diary.server'
import { DiaryFolderCard, DiaryList } from '@/components/diary'

type QueryValue = string | string[] | undefined

interface DiaryPageProps {
  searchParams: Promise<{
    page?: QueryValue
    keyword?: QueryValue
    folder?: QueryValue
  }>
}

function firstValue(value: QueryValue): string {
  return Array.isArray(value) ? value[0] ?? '' : value ?? ''
}

function positiveInteger(value: QueryValue): number | undefined {
  const parsed = Number(firstValue(value))
  return Number.isInteger(parsed) && parsed > 0 ? parsed : undefined
}

function createDiaryHref(
  page: number,
  keyword: string,
  folderID?: number,
): string {
  const params = new URLSearchParams()
  if (page > 1)
    params.set('page', String(page))
  if (keyword)
    params.set('keyword', keyword)
  if (folderID)
    params.set('folder', String(folderID))

  const query = params.toString()
  return `/diary${query ? `?${query}` : ''}`
}

export default async function DiaryPage({ searchParams }: DiaryPageProps) {
  const params = await searchParams
  const page = positiveInteger(params.page) ?? 1
  const folderID = positiveInteger(params.folder)
  const keyword = firstValue(params.keyword).trim()

  try {
    const [result, folders] = await Promise.all([
      listDiaries({
        page,
        pageSize: 12,
        status: 'published',
        keyword,
        folderId: folderID,
      }),
      listDiaryFolders(),
    ])
    const totalPages = Math.max(1, Math.ceil(result.total / result.pageSize))

    return (
      <section aria-labelledby="diary-title" className="space-y-8">
        <header className="border-b border-black/10 pb-6 dark:border-white/10">
          <h1 id="diary-title" className="text-3xl font-semibold">日记</h1>
          <p className="mt-2 text-neutral-600 dark:text-neutral-400">
            记录日常片段与阶段思考
          </p>
        </header>

        {folders.length > 0 && (
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {folders.map(folder => (
              <DiaryFolderCard
                key={folder.id}
                folder={folder}
                active={folder.id === folderID}
              />
            ))}
          </div>
        )}

        <form action="/diary" className="flex flex-wrap gap-3">
          <input
            type="search"
            name="keyword"
            defaultValue={keyword}
            placeholder="搜索日记"
            aria-label="搜索日记"
            className="min-h-10 min-w-56 flex-1 rounded-md border border-black/15 px-3 dark:border-white/15"
          />
          {folderID && <input type="hidden" name="folder" value={folderID} />}
          <button
            type="submit"
            className="min-h-10 rounded-md bg-black px-4 text-sm text-white dark:bg-white dark:text-black"
          >
            搜索
          </button>
          {(keyword || folderID) && (
            <Link href="/diary" className="inline-flex min-h-10 items-center px-2 text-sm">
              清除筛选
            </Link>
          )}
        </form>

        <DiaryList diaries={result.items} />

        <nav aria-label="日记分页" className="flex min-h-10 items-center justify-between text-sm">
          {page > 1
            ? <Link href={createDiaryHref(page - 1, keyword, folderID)}>上一页</Link>
            : <span className="text-neutral-400">上一页</span>}
          <span>
            {result.page}
            {' / '}
            {totalPages}
          </span>
          {page < totalPages
            ? <Link href={createDiaryHref(page + 1, keyword, folderID)}>下一页</Link>
            : <span className="text-neutral-400">下一页</span>}
        </nav>
      </section>
    )
  }
  catch (error) {
    const message = error instanceof DiaryServerError
      ? error.message
      : '日记列表暂时无法加载'

    return (
      <section className="border-y border-black/10 py-12 dark:border-white/10">
        <h1 className="text-2xl font-semibold">日记加载失败</h1>
        <p className="mt-3 text-neutral-600 dark:text-neutral-400">{message}</p>
      </section>
    )
  }
}
