import { notFound, redirect } from 'next/navigation'
import { getCurrentUser } from '@/lib/auth/current-user'
import { RoleApplicationManage } from '../role-application-manage'

export default async function ConsoleRoleApplicationsPage() {
  const currentUser = await getCurrentUser()

  if (!currentUser)
    redirect('/auth/login?next=/console/roles/applications')
  if (currentUser.role !== 'admin')
    notFound()

  return <RoleApplicationManage />
}
