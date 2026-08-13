'use client'

import type { ReactNode } from 'react'
import * as Tooltip from '@radix-ui/react-tooltip'

interface EditorTooltipProps {
  label: string
  children: ReactNode
}

export function EditorTooltip({ label, children }: EditorTooltipProps) {
  return (
    <Tooltip.Provider delayDuration={400}>
      <Tooltip.Root>
        <Tooltip.Trigger asChild>
          {children}
        </Tooltip.Trigger>
        <Tooltip.Portal>
          <Tooltip.Content
            sideOffset={5}
            className="z-50 rounded-sm bg-black px-2 py-1 text-xs text-white dark:bg-white dark:text-black"
          >
            {label}
          </Tooltip.Content>
        </Tooltip.Portal>
      </Tooltip.Root>
    </Tooltip.Provider>
  )
}
