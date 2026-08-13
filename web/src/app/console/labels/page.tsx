import { notFound } from 'next/navigation'
import { listLabels } from '@/api/post.server'
import { getCurrentUser } from '@/lib/auth/current-user'
import { LabelManage } from './label-manage'

export default async function LabelsPage() {
  const user = await getCurrentUser()

  if (
    !user
    || (user.role !== 'admin' && user.role !== 'editor')
  ) {
    notFound()
  }

  const labels = await listLabels()

  return <LabelManage labels={labels} />
}
