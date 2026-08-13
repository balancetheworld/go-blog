import { notFound } from 'next/navigation'
import { listCategories } from '@/api/post.server'
import { getCurrentUser } from '@/lib/auth/current-user'
import { CategoryManage } from './category-manage'

export default async function CategoriesPage() {
  const user = await getCurrentUser()

  if (
    !user
    || (user.role !== 'admin' && user.role !== 'editor')
  ) {
    notFound()
  }

  const categories = await listCategories()

  return <CategoryManage categories={categories} />
}
