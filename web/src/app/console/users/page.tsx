import { getLocale, getTranslations } from 'next-intl/server'
import Link from 'next/link'
import { notFound } from 'next/navigation'
import { listUsers } from '@/api/user.server'
import { getCurrentUser } from '@/lib/auth/current-user'

interface UsersPageProps {
  searchParams: Promise<{
    page?: string
  }>
}

export default async function UsersPage({
  searchParams,
}: UsersPageProps) {
  const locale = await getLocale()
  const t = await getTranslations('Console.users')
  const currentUser = await getCurrentUser()
  if (!currentUser || currentUser.role !== 'admin')
    notFound()
  const params = await searchParams
  const requestedPage = Number(params.page)
  const page = Number.isInteger(requestedPage) && requestedPage > 0 ? requestedPage : 1
  const result = await listUsers({
    page,
    pageSize: 20,
  })
  const totalPages = Math.max(
    1,
    Math.ceil(result.total / result.pageSize),
  )
  return (
    <section aria-labelledby="users-title" className="space-y-6">
      <header>
        <h1 id="users-title" className="text-2xl font-semibold">
          {t('title')}
        </h1>
        <p className="mt-1 text-sm text-neutral-500">
          {t('count', { count: result.total })}
        </p>
      </header>

      <div className="overflow-x-auto border-y border-black/10 dark:border-white/10">
        <table className="w-full min-w-[720px] text-left text-sm">
          <thead className="text-neutral-500">
            <tr>
              <th className="px-3 py-3">{t('user')}</th>
              <th className="px-3 py-3">{t('username')}</th>
              <th className="px-3 py-3">{t('role')}</th>
              <th className="px-3 py-3">{t('registeredAt')}</th>
            </tr>
          </thead>

          <tbody>
            {result.items.map(user => (
              <tr
                key={user.id}
                className="border-t border-black/10 dark:border-white/10"
              >
                <td className="px-3 py-4">
                  {user.nickname || user.username}
                </td>
                <td className="px-3 py-4">{user.username}</td>
                <td className="px-3 py-4">
                  {t(user.role === 'user' ? 'member' : user.role)}
                </td>
                <td className="px-3 py-4">
                  {new Intl.DateTimeFormat(locale, { dateStyle: 'medium' }).format(new Date(user.createdAt))}
                </td>
              </tr>
            ))}

            {result.items.length === 0 && (
              <tr>
                <td
                  colSpan={4}
                  className="px-3 py-12 text-center text-neutral-500"
                >
                  {t('empty')}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <nav aria-label={t('pagination')} className="flex justify-end gap-4 text-sm">
        {page > 1
          ? <Link href={`/console/users?page=${page - 1}`}>{t('previous')}</Link>
          : <span className="text-neutral-400">{t('previous')}</span>}

        <span>
          {result.page}
          {' '}
          /
          {' '}
          {totalPages}
        </span>

        {page < totalPages
          ? <Link href={`/console/users?page=${page + 1}`}>{t('next')}</Link>
          : <span className="text-neutral-400">{t('next')}</span>}
      </nav>
    </section>
  )
}
