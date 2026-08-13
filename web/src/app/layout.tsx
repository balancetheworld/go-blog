import type { Metadata } from 'next'
import type { ReactNode } from 'react'
import { getLocale, getMessages } from 'next-intl/server'
import { Providers } from '@/components/providers'
import 'highlight.js/styles/github-dark.css'
import './globals.css'

export const metadata: Metadata = {
  title: '我的博客',
  description: '记录技术、生活与思考',
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
