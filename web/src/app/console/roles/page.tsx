import { notFound, redirect } from 'next/navigation'
import { getCurrentUser } from '@/lib/auth/current-user'
import { RoleManage } from './role-manage'

export default async function ConsoleRolesPage() {
  const currentUser = await getCurrentUser()

  if (!currentUser)
    redirect('/auth/login?next=/console/roles')
  if (currentUser.role !== 'admin')
    notFound()
  return <RoleManage />
}
