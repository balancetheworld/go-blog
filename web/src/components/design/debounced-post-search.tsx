'use client'

import { Search } from 'lucide-react'
import { useTranslations } from 'next-intl'
import { usePathname, useRouter, useSearchParams } from 'next/navigation'
import { useEffect, useState, useTransition } from 'react'

interface DebouncedPostSearchProps {
  initialKeyword: string
}

export function DebouncedPostSearch({
  initialKeyword,
}: DebouncedPostSearchProps) {
  const t = useTranslations('Home.articles')
  const pathname = usePathname()
  const router = useRouter()
  const searchParams = useSearchParams()
  const [keyword, setKeyword] = useState(initialKeyword)
  const [pending, startTransition] = useTransition()

  useEffect(() => {
    setKeyword(initialKeyword)
  }, [initialKeyword])

  useEffect(() => {
    function updateSearch(value: string) {
      const normalized = value.trim()
      if (normalized === (searchParams.get('keyword') ?? ''))
        return

      const params = new URLSearchParams(searchParams.toString())
      params.delete('page')
      if (normalized)
        params.set('keyword', normalized)
      else
        params.delete('keyword')

      const query = params.toString()
      startTransition(() => {
        router.replace(query ? `${pathname}?${query}` : pathname, { scroll: false })
      })
    }

    const timer = window.setTimeout(updateSearch, 350, keyword)
    return () => window.clearTimeout(timer)
  }, [keyword, pathname, router, searchParams])

  return (
    <label className="article-search" data-pending={pending || undefined}>
      <Search className="search-icon" aria-hidden="true" />
      <input
        type="search"
        value={keyword}
        placeholder={t('searchPlaceholder')}
        aria-label={t('search')}
        onChange={event => setKeyword(event.target.value)}
      />
      {keyword && <input type="hidden" name="keyword" value={keyword} />}
    </label>
  )
}
