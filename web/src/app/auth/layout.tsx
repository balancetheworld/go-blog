import type { ReactNode } from 'react'
import { DesignShell } from '@/components/design/design-shell'
import '../design.css'

interface AuthLayoutProps {
  children: ReactNode
}

export default function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <DesignShell auth footer={false}>
      {children}
    </DesignShell>
  )
}
