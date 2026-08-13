import type { Metadata } from 'next'
import type { ReactNode } from 'react'
import { getLocale, getMessages, getTranslations } from 'next-intl/server'
import { Providers } from '@/components/providers'
import 'highlight.js/styles/github-dark.css'
import './globals.css'

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations('Metadata')

  return {
    title: t('title'),
    description: t('description'),
  }
}

interface RootLayoutProps {
  children: ReactNode
}

export default async function RootLayout({
  children,
}: RootLayoutProps) {
  const locale = await getLocale()
  const messages = await getMessages()

  return (
    <html lang={locale} suppressHydrationWarning>
      <body>
        <Providers
          locale={locale}
          messages={messages}
        >
          {children}
        </Providers>
      </body>
    </html>
  )
}
