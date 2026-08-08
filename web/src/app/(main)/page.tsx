import type { PostSort } from '@/models/post'
import Link from 'next/link'
import {
  listCategories,
  listLabels,
  listPosts,
  PostServerError,
} from '@/api/post.server'
import PostList from '@/components/home/post-list'

type QueryValue = string | string[] | undefined

interface HomePageProps {
  searchParams: Promise<{
    page?: QueryValue
    keyword?: QueryValue
    category?: QueryValue
    label?: QueryValue
    sort?: QueryValue
  }>
}

interface HomeQuery {
  page: number
  keyword: string
  categoryId?: number
  labelId?: number
  sort: PostSort
}

function firstValue(value: QueryValue): string {
  return Array.isArray(value) ? value[0] ?? '' : value ?? ''
}

function positiveInteger(value: QueryValue): number | undefined {
  const parsed = Number(firstValue(value))
  return Number.isInteger(parsed) && parsed > 0 ? parsed : undefined
}

function createHomeHref(query: HomeQuery): string {
  const params = new URLSearchParams()
  if (query.page > 1)
    params.set('page', String(query.page))
  if (query.keyword)
    params.set('keyword', query.keyword)
  if (query.categoryId)
    params.set('category', String(query.categoryId))
  if (query.labelId)
    params.set('label', String(query.labelId))
  if (query.sort !== 'latest')
    params.set('sort', query.sort)

  const value = params.toString()
  return value ? `/?${value}` : '/'
}

export default async function HomePage({ searchParams }: HomePageProps) {
  const params = await searchParams
  const requestedSort = firstValue(params.sort)
  const query: HomeQuery = {
    page: positiveInteger(params.page) ?? 1,
    keyword: firstValue(params.keyword).trim(),
    categoryId: positiveInteger(params.category),
    labelId: positiveInteger(params.label),
    sort: requestedSort === 'oldest' || requestedSort === 'hot'
      ? requestedSort
      : 'latest',
  }

  try {
    const [result, categories, labels] = await Promise.all([
      listPosts({
        page: query.page,
        pageSize: 10,
        keyword: query.keyword,
        categoryId: query.categoryId,
        labelId: query.labelId,
        status: 'published',
        sort: query.sort,
      }),
      listCategories(),
      listLabels(),
    ])
    const totalPages = Math.max(1, Math.ceil(result.total / result.pageSize))

    return (
      <section aria-labelledby="latest-posts-title" className="space-y-8">
        <header className="border-b border-black/10 pb-6 dark:border-white/10">
          <h1 id="latest-posts-title" className="text-3xl font-semibold">最新文章</h1>
          <p className="mt-2 text-neutral-600 dark:text-neutral-400">
            浏览最近发布的公开内容
          </p>
        </header>

        <form action="/" className="grid gap-3 border-b border-black/10 pb-6 sm:grid-cols-2 lg:grid-cols-[minmax(0,1fr)_180px_180px_150px_auto] dark:border-white/10">
          <input
            type="search"
            name="keyword"
            defaultValue={query.keyword}
            placeholder="搜索标题或摘要"
            aria-label="搜索文章"
            className="min-h-10 min-w-0 rounded-md border border-black/15 px-3 dark:border-white/15"
          />
          <select
            name="category"
            defaultValue={query.categoryId ?? ''}
            aria-label="按分类筛选"
            className="min-h-10 rounded-md border border-black/15 px-3 dark:border-white/15"
          >
            <option value="">全部分类</option>
            {categories.map(category => (
              <option key={category.id} value={category.id}>{category.name}</option>
            ))}
          </select>
          <select
            name="label"
            defaultValue={query.labelId ?? ''}
            aria-label="按标签筛选"
            className="min-h-10 rounded-md border border-black/15 px-3 dark:border-white/15"
          >
            <option value="">全部标签</option>
            {labels.map(label => (
              <option key={label.id} value={label.id}>{label.name}</option>
            ))}
          </select>
          <select
            name="sort"
            defaultValue={query.sort}
            aria-label="文章排序"
            className="min-h-10 rounded-md border border-black/15 px-3 dark:border-white/15"
          >
            <option value="latest">最新发布</option>
            <option value="oldest">最早发布</option>
            <option value="hot">热度最高</option>
          </select>
          <button
            type="submit"
            className="min-h-10 rounded-md bg-black px-4 text-sm text-white dark:bg-white dark:text-black"
          >
            筛选
          </button>
        </form>

        {(query.keyword || query.categoryId || query.labelId || query.sort !== 'latest') && (
          <div className="flex items-center justify-between gap-4 text-sm">
            <span className="text-neutral-500">
              找到
              {' '}
              {result.total}
              {' '}
              篇文章
            </span>
            <Link href="/" className="font-medium">清除筛选</Link>
          </div>
        )}

        <PostList posts={result.items} />

        <nav aria-label="文章分页" className="flex min-h-10 items-center justify-between border-t border-black/10 pt-6 text-sm dark:border-white/10">
          {query.page > 1
            ? <Link href={createHomeHref({ ...query, page: query.page - 1 })}>上一页</Link>
            : <span className="text-neutral-400">上一页</span>}
          <span>
            {result.page}
            {' '}
            /
            {' '}
            {totalPages}
          </span>
          {query.page < totalPages
            ? <Link href={createHomeHref({ ...query, page: query.page + 1 })}>下一页</Link>
            : <span className="text-neutral-400">下一页</span>}
        </nav>
      </section>
    )
  }
  catch (error) {
    const message = error instanceof PostServerError
      ? error.message
      : '文章列表暂时无法加载'

    return (
      <section className="border-y border-black/10 py-12 dark:border-white/10">
        <h1 className="text-2xl font-semibold">文章加载失败</h1>
        <p className="mt-3 text-neutral-600 dark:text-neutral-400">{message}</p>
      </section>
    )
  }
}
