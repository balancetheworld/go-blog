import { notFound } from 'next/navigation'
import { getCurrentUser } from '@/lib/auth/current-user'

export default async function OIDCPage() {
  const user = await getCurrentUser()

  if (!user || user.role !== 'admin')
    notFound()

  return (
    <section aria-labelledby="oidc-title" className="space-y-6">
      <header>
        <h1 id="oidc-title" className="text-2xl font-semibold">
          OIDC 配置
        </h1>
        <p className="mt-1 text-sm text-neutral-500">
          OIDC 接口完成后，这里将管理身份提供方配置。
        </p>
      </header>

      <div className="border-y border-black/10 py-12 text-center text-sm text-neutral-500 dark:border-white/10">
        暂无 OIDC 配置
      </div>
    </section>
  )
}
