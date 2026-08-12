import type { ReactNode } from 'react'
import { notFound, redirect } from 'next/navigation'
import { AppSidebar } from '@/components/console/app-sidebar'
import { ConsoleContent } from '@/components/console/console-content'
import { SiteHeader } from '@/components/console/site-header'
import { getCurrentUser } from '@/lib/auth/current-user'

interface ConsoleLayoutProps {
  children: ReactNode
}

export default async function ConsoleLayout({
  children,
}: ConsoleLayoutProps) {
  const user = await getCurrentUser()

  if (!user) {
    redirect('/auth/login?next=/console')
  }

  if (user.role !== 'admin' && user.role !== 'editor') {
    notFound()
  }

  return (
    <div className="grid min-h-screen md:grid-cols-[220px_minmax(0,1fr)]">
      <AppSidebar role={user.role} />

      <div className="min-w-0">
        <SiteHeader user={user} />

        <ConsoleContent>
          {children}
        </ConsoleContent>
      </div>
    </div>
  )
}
