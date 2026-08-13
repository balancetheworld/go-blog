'use client'

import type { ReactNode } from 'react'
import * as Popover from '@radix-ui/react-popover'

interface EditorPopoverProps {
  trigger: ReactNode
  children: ReactNode
  ariaLabel: string
}

export function EditorPopover({
  trigger,
  children,
  ariaLabel,
}: EditorPopoverProps) {
  return (
    <Popover.Root>
      <Popover.Trigger asChild>
        {trigger}
      </Popover.Trigger>
      <Popover.Portal>
        <Popover.Content
          aria-label={ariaLabel}
          sideOffset={6}
          align="start"
          className="z-50 rounded-md border border-black/15 bg-white p-3 shadow-lg dark:border-white/15 dark:bg-neutral-950"
        >
          {children}
        </Popover.Content>
      </Popover.Portal>
    </Popover.Root>
  )
}
