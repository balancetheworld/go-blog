import { notFound } from 'next/navigation'
import { getCurrentUser } from '@/lib/auth/current-user'

export default async function StoragesPage() {
  const user = await getCurrentUser()

  if (!user || user.role !== 'admin')
    notFound()

  return (
    <section aria-labelledby="storages-title" className="space-y-6">
      <header>
        <h1 id="storages-title" className="text-2xl font-semibold">
          存储管理
        </h1>
        <p className="mt-1 text-sm text-neutral-500">
          存储接口完成后，这里将展示存储配置和连接状态。
        </p>
      </header>

      <div className="border-y border-black/10 py-12 text-center text-sm text-neutral-500 dark:border-white/10">
        暂无存储配置
      </div>
    </section>
  )
}
