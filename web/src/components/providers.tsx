'use client'

import type { AbstractIntlMessages } from 'next-intl'
import type { ReactNode } from 'react'
import { NextIntlClientProvider } from 'next-intl'
import { ThemeProvider } from 'next-themes'
import { Toaster } from 'sonner'
import { AuthProvider } from '@/contexts/auth-context'

interface ProvidersProps {
  children: ReactNode
  locale: string
  messages: AbstractIntlMessages
}

export function Providers({ children, locale, messages }: ProvidersProps) {
  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      <ThemeProvider
        attribute="class"
        defaultTheme="system"
        enableSystem
        disableTransitionOnChange
      >
        <AuthProvider>
          {children}
          <Toaster
            theme="system"
            position="top-right"
            richColors
            closeButton
          />
        </AuthProvider>
      </ThemeProvider>
    </NextIntlClientProvider>
  )
}
