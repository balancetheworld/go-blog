import type { ReactNode } from 'react'

interface ConsoleContentProps {
  children: ReactNode
}

export function ConsoleContent({
  children,
}: ConsoleContentProps) {
  return (
    <main className="console-main">
      <div className="console-content">
        {children}
      </div>
    </main>
  )
}
