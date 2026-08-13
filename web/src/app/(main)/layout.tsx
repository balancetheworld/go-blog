import type { ReactNode } from 'react'
import { DesignShell } from '@/components/design/design-shell'
import '../design.css'

interface MainLayoutProps {
  children: ReactNode
}

export default function MainLayout({
  children,
}: MainLayoutProps) {
  return (
    <DesignShell>
      {children}
    </DesignShell>
  )
}
