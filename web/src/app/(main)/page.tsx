import type { PostSort } from '@/models/post'
import { getTranslations } from 'next-intl/server'
import { headers } from 'next/headers'
import Link from 'next/link'
import { listDiaries, listDiaryFolders } from '@/api/diary.server'
import { listMoments } from '@/api/moment.server'
import {
  listCategories,
  listPosts,
  PostServerError,
} from '@/api/post.server'
import { DebouncedPostSearch } from '@/components/design/debounced-post-search'
import { FilterSelect } from '@/components/design/filter-select'
import { HomeDiaryArchive } from '@/components/home/home-diary-archive'
import { HomeExperience } from '@/components/home/home-experience'
import PostList from '@/components/home/post-list'

type QueryValue = string | string[] | undefined

interface HomePageProps {
  searchParams: Promise<{
    page?: QueryValue
    keyword?: QueryValue
    category?: QueryValue
    sort?: QueryValue
  }>
}

interface HomeQuery {
  page: number
  keyword: string
  categoryId?: number
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
  if (query.sort !== 'latest')
    params.set('sort', query.sort)

  const value = params.toString()
  return value ? `/?${value}` : '/'
}

export default async function HomePage({ searchParams }: HomePageProps) {
  const articleT = await getTranslations('Home.articles')
  const diaryT = await getTranslations('Home.diaries')
  const params = await searchParams
  const requestHeaders = await headers()
  const userAgent = requestHeaders.get('user-agent') ?? ''
  const pageSize = /Android|iPhone|iPad|iPod|Mobile/i.test(userAgent) ? 5 : 10
  const requestedSort = firstValue(params.sort)
  const query: HomeQuery = {
    page: positiveInteger(params.page) ?? 1,
    keyword: firstValue(params.keyword).trim(),
    categoryId: positiveInteger(params.category),
    sort: requestedSort === 'oldest' || requestedSort === 'hot'
      ? requestedSort
      : 'latest',
  }

  try {
    const [result, categories, diaries, diaryFolders, moments] = await Promise.all([
      listPosts({
        page: query.page,
        pageSize,
        keyword: query.keyword,
        categoryId: query.categoryId,
        status: 'published',
        sort: query.sort,
      }),
      listCategories(),
      listDiaries({
        page: 1,
        pageSize: 6,
        status: 'published',
      }).catch(() => null),
      listDiaryFolders().catch(() => []),
      listMoments({ page: 1, pageSize: 6 }).catch(() => null),
    ])
    const totalPages = Math.max(1, Math.ceil(result.total / result.pageSize))

    return (
      <main>
        <section aria-labelledby="latest-posts-title" className="section articles" id="articles">
          <div className="section-header">
            <h1 id="latest-posts-title" className="section-title">
              <span className="title-num">01</span>
              <span className="title-text">{articleT('title')}</span>
            </h1>
            <p className="section-sub">{articleT('subtitle')}</p>
          </div>

          <form action="/" className="article-toolbar">
            <div className="article-filter" aria-label={articleT('category')}>
              <Link
                href={createHomeHref({ ...query, page: 1, categoryId: undefined })}
                className={`filter-chip${query.categoryId ? '' : ' active'}`}
              >
                {articleT('all')}
              </Link>
              {categories.map(category => (
                <Link
                  key={category.id}
                  href={createHomeHref({ ...query, page: 1, categoryId: category.id })}
                  className={`filter-chip${query.categoryId === category.id ? ' active' : ''}`}
                >
                  {category.name}
                </Link>
              ))}
            </div>
            <div className="article-category-select">
              <FilterSelect
                name="category"
                defaultValue={query.categoryId ? String(query.categoryId) : 'all'}
                ariaLabel={articleT('category')}
                submitOnChange
                options={[
                  { label: articleT('allCategories'), value: 'all' },
                  ...categories.map(category => ({
                    label: category.name,
                    value: String(category.id),
                  })),
                ]}
              />
            </div>
            <DebouncedPostSearch initialKeyword={query.keyword} />
            <FilterSelect
              name="sort"
              defaultValue={query.sort}
              ariaLabel={articleT('sort')}
              submitOnChange
              options={[
                { label: articleT('latest'), value: 'latest' },
                { label: articleT('oldest'), value: 'oldest' },
                { label: articleT('hot'), value: 'hot' },
              ]}
            />
          </form>

          {(query.keyword || query.categoryId || query.sort !== 'latest') && (
            <div className="design-result-summary">
              <span>{articleT('result', { count: result.total })}</span>
              <Link href="/">{articleT('clear')}</Link>
            </div>
          )}

          <PostList posts={result.items} />

          <nav aria-label={articleT('pagination')} className="design-pagination">
            {query.page > 1
              ? <Link href={createHomeHref({ ...query, page: query.page - 1 })}>{articleT('previous')}</Link>
              : <span aria-disabled="true">{articleT('previous')}</span>}
            <span>
              {result.page}
              {' / '}
              {totalPages}
            </span>
            {query.page < totalPages
              ? <Link href={createHomeHref({ ...query, page: query.page + 1 })}>{articleT('next')}</Link>
              : <span aria-disabled="true">{articleT('next')}</span>}
          </nav>
        </section>

        {diaries && (
          <section aria-labelledby="home-diary-title" className="section diary" id="diary">
            <div className="section-header">
              <h2 id="home-diary-title" className="section-title">
                <span className="title-num">02</span>
                <span className="title-text">{diaryT('title')}</span>
              </h2>
              <p className="section-sub">{diaryT('subtitle')}</p>
            </div>
            <HomeDiaryArchive diaries={diaries.items} folders={diaryFolders} />
          </section>
        )}
        <HomeExperience
          posts={result.items}
          diaries={diaries?.items ?? []}
          moments={moments?.items ?? []}
        />
      </main>
    )
  }
  catch (error) {
    const message = error instanceof PostServerError
      ? error.message
      : articleT('unavailable')

    return (
      <main>
        <section className="section articles">
          <div className="article-empty visible">
            <h1>{articleT('loadFailed')}</h1>
            <p>{message}</p>
          </div>
        </section>
      </main>
    )
  }
}
