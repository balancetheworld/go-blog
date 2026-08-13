import { getTranslations } from 'next-intl/server'
import { notFound } from 'next/navigation'
import { getCurrentUser } from '@/lib/auth/current-user'

export default async function StoragesPage() {
  const user = await getCurrentUser()
  const t = await getTranslations('Console.placeholders')

  if (!user || user.role !== 'admin')
    notFound()

  return (
    <section aria-labelledby="storages-title" className="space-y-6">
      <header>
        <h1 id="storages-title" className="text-2xl font-semibold">
          {t('storagesTitle')}
        </h1>
        <p className="mt-1 text-sm text-neutral-500">
          {t('storagesDescription')}
        </p>
      </header>

      <div className="border-y border-black/10 py-12 text-center text-sm text-neutral-500 dark:border-white/10">
        {t('storagesEmpty')}
      </div>
    </section>
  )
}
