import { getTranslations } from 'next-intl/server'
import Link from 'next/link'
import {
  DiaryServerError,
  listDiaries,
  listDiaryFolders,
} from '@/api/diary.server'
import { DiaryFolderCard, DiaryList } from '@/components/diary'
import { DiaryFolderMedia } from '@/components/diary/diary-folder-media'
import { allDiaryFolderMedia } from '@/lib/diary-folders'

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
  const t = await getTranslations('Diary')
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
      <main>
        <section aria-labelledby="diary-title" className="section diary" id="diary">
          <div className="section-header">
            <h1 id="diary-title" className="section-title">
              <span className="title-num">02</span>
              <span className="title-text">{t('title')}</span>
            </h1>
            <p className="section-sub">{t('subtitle')}</p>
          </div>

          <div className="diary-archive design-diary-live">
            <div className="diary-topline">
              <div className="diary-months">
                <span className="diary-month active">{t('allRecords')}</span>
              </div>
              <form action="/diary" className="diary-search-form">
                <label className="article-search">
                  <svg className="search-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                    <circle cx="11" cy="11" r="8" />
                    <line x1="21" y1="21" x2="16.65" y2="16.65" />
                  </svg>
                  <input type="search" name="keyword" defaultValue={keyword} placeholder={t('searchLabel')} aria-label={t('searchLabel')} />
                </label>
                {folderID && <input type="hidden" name="folder" value={folderID} />}
                <button type="submit" className="filter-chip active">{t('search')}</button>
              </form>
            </div>

            <div className={`diary-folder-view${folderID ? ' has-active' : ''}`}>
              <Link href="/diary" className={`diary-folder diary-folder-all${folderID ? '' : ' is-active'}`}>
                <span className="folder-paper paper-one" />
                <span className="folder-paper paper-two" />
                <span className="folder-back" />
                <span className="folder-front">
                  <DiaryFolderMedia media={allDiaryFolderMedia} />
                </span>
                <span className="folder-tab">{t('all')}</span>
                <span className="folder-meta">{t('allDescription')}</span>
              </Link>
              {folders.map(folder => (
                <DiaryFolderCard key={folder.id} folder={folder} active={folder.id === folderID} />
              ))}
            </div>

            <div className="diary-atmosphere" aria-hidden="true">
              <span className="diary-float diary-float-photo" />
              <span className="diary-float diary-float-ticket" />
              <span className="diary-float diary-float-note" />
              <span className="diary-float diary-float-dot" />
              <span className="diary-orbit">
                <i />
                <i />
                <i />
              </span>
            </div>

            <div className="diary-list-view">
              <div className="diary-list-head">
                <h2 className="diary-list-title">{folderID ? t('folderDiaries') : t('allDiaries')}</h2>
                {(keyword || folderID) && <Link href="/diary" className="diary-back">{t('clear')}</Link>}
              </div>
              <DiaryList diaries={result.items} />
              <nav aria-label={t('pagination')} className="design-pagination">
                {page > 1
                  ? <Link href={createDiaryHref(page - 1, keyword, folderID)}>{t('previous')}</Link>
                  : <span aria-disabled="true">{t('previous')}</span>}
                <span>
                  {result.page}
                  {' / '}
                  {totalPages}
                </span>
                {page < totalPages
                  ? <Link href={createDiaryHref(page + 1, keyword, folderID)}>{t('next')}</Link>
                  : <span aria-disabled="true">{t('next')}</span>}
              </nav>
            </div>
          </div>
        </section>
      </main>
    )
  }
  catch (error) {
    const message = error instanceof DiaryServerError
      ? error.message
      : t('unavailable')

    return (
      <main>
        <section className="section diary">
          <div className="article-empty visible">
            <h1>{t('loadFailed')}</h1>
            <p>{message}</p>
          </div>
        </section>
      </main>
    )
  }
}
