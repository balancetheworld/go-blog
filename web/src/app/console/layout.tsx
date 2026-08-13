import type { ReactNode } from 'react'
import { redirect } from 'next/navigation'
import { ConsoleAccessDenied } from '@/components/console/access-denied'
import { AppSidebar } from '@/components/console/app-sidebar'
import { ConsoleContent } from '@/components/console/console-content'
import { SiteHeader } from '@/components/console/site-header'
import { getCurrentUser } from '@/lib/auth/current-user'
import { isEditor } from '@/lib/permission'

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

  if (!isEditor(user.role))
    return <ConsoleAccessDenied />

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
