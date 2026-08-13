import { getTranslations } from 'next-intl/server'
import Link from 'next/link'
import { notFound } from 'next/navigation'
import { listMoments } from '@/api/moment.server'
import { getCurrentUser } from '@/lib/auth/current-user'
import { MomentManage } from './moment-manage'

interface MomentsPageProps {
  searchParams: Promise<{
    page?: string
  }>
}

export default async function MomentsPage({ searchParams }: MomentsPageProps) {
  const t = await getTranslations('Console.moments')
  const user = await getCurrentUser()
  if (!user || (user.role !== 'admin' && user.role !== 'editor'))
    notFound()

  const params = await searchParams
  const requestedPage = Number(params.page)
  const page = Number.isInteger(requestedPage) && requestedPage > 0
    ? requestedPage
    : 1
  const result = await listMoments({ page, pageSize: 20 })
  const totalPages = Math.max(1, Math.ceil(result.total / result.pageSize))

  return (
    <section aria-labelledby="moments-title" className="space-y-6">
      <MomentManage
        result={result}
        currentUserId={user.id}
        isAdmin={user.role === 'admin'}
      />

      <nav aria-label={t('pagination')} className="flex justify-end gap-4 text-sm">
        {page > 1
          ? <Link href={`/console/moments?page=${page - 1}`}>{t('previous')}</Link>
          : <span className="text-neutral-400">{t('previous')}</span>}

        <span>
          {result.page}
          {' / '}
          {totalPages}
        </span>

        {page < totalPages
          ? <Link href={`/console/moments?page=${page + 1}`}>{t('next')}</Link>
          : <span className="text-neutral-400">{t('next')}</span>}
      </nav>
    </section>
  )
}
