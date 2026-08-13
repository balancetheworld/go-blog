import { getTranslations } from 'next-intl/server'
import Link from 'next/link'

export default async function PostNotFound() {
  const t = await getTranslations('Post')

  return (
    <main className="article-detail-main">
      <section className="article-detail-shell article-empty visible">
        <h1>
          {t('unavailableTitle')}
        </h1>
        <p>{t('unavailableDescription')}</p>
        <Link href="/" className="article-back-link">{t('backHome')}</Link>
      </section>
    </main>
  )
}
