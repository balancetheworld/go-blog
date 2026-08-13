import { notFound } from 'next/navigation'
import { getCurrentUser } from '@/lib/auth/current-user'

export default async function GlobalPage() {
  const user = await getCurrentUser()

  if (!user || user.role !== 'admin')
    notFound()

  return (
    <section aria-labelledby="global-title" className="space-y-6">
      <header>
        <h1 id="global-title" className="text-2xl font-semibold">
          全局配置
        </h1>
        <p className="mt-1 text-sm text-neutral-500">
          配置接口完成后，这里将管理站点名称、外观和基础设置。
        </p>
      </header>

      <div className="border-y border-black/10 py-12 text-center text-sm text-neutral-500 dark:border-white/10">
        暂无全局配置
      </div>
    </section>
  )
}
