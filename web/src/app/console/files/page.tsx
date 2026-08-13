import { getTranslations } from 'next-intl/server'

export default async function FilesPage() {
  const t = await getTranslations('Console.placeholders')

  return (
    <section aria-labelledby="files-title" className="space-y-6">
      <header>
        <h1 id="files-title" className="text-2xl font-semibold">
          {t('filesTitle')}
        </h1>
        <p className="mt-1 text-sm text-neutral-500">
          {t('filesDescription')}
        </p>
      </header>

      <div className="border-y border-black/10 py-12 text-center text-sm text-neutral-500 dark:border-white/10">
        {t('filesEmpty')}
      </div>
    </section>
  )
}
