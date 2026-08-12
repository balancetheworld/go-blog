import type { ReactNode } from 'react'

interface ConsoleContentProps {
  children: ReactNode
}

export function ConsoleContent({
  children,
}: ConsoleContentProps) {
  return (
    <main className="min-w-0 p-4 sm:p-6 lg:p-8">
      {children}
    </main>
  )
}
