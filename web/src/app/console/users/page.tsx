import Link from 'next/link'
import { notFound } from 'next/navigation'
import { listUsers } from '@/api/user.server'
import { getCurrentUser } from '@/lib/auth/current-user'

interface UsersPageProps {
  searchParams: Promise<{
    page?: string
  }>
}

const dateFormatter = new Intl.DateTimeFormat('zh-CN', {
  dateStyle: 'medium',
})

const roleNames = {
  user: '普通用户',
  editor: '编辑者',
  admin: '管理员',
}

export default async function UsersPage({
  searchParams,
}: UsersPageProps) {
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
          用户管理
        </h1>
        <p className="mt-1 text-sm text-neutral-500">
          共
          {' '}
          {result.total}
          {' '}
          位用户
        </p>
      </header>

      <div className="overflow-x-auto border-y border-black/10 dark:border-white/10">
        <table className="w-full min-w-[720px] text-left text-sm">
          <thead className="text-neutral-500">
            <tr>
              <th className="px-3 py-3">用户</th>
              <th className="px-3 py-3">用户名</th>
              <th className="px-3 py-3">权限</th>
              <th className="px-3 py-3">注册时间</th>
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
                  {roleNames[user.role]}
                </td>
                <td className="px-3 py-4">
                  {dateFormatter.format(new Date(user.createdAt))}
                </td>
              </tr>
            ))}

            {result.items.length === 0 && (
              <tr>
                <td
                  colSpan={4}
                  className="px-3 py-12 text-center text-neutral-500"
                >
                  暂无用户
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <nav aria-label="用户分页" className="flex justify-end gap-4 text-sm">
        {page > 1
          ? <Link href={`/console/users?page=${page - 1}`}>上一页</Link>
          : <span className="text-neutral-400">上一页</span>}

        <span>
          {result.page}
          {' '}
          /
          {' '}
          {totalPages}
        </span>

        {page < totalPages
          ? <Link href={`/console/users?page=${page + 1}`}>下一页</Link>
          : <span className="text-neutral-400">下一页</span>}
      </nav>
    </section>
  )
}
