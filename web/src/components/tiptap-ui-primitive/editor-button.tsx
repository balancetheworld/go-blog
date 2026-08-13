import type { ButtonHTMLAttributes, ReactNode } from 'react'
import { forwardRef } from 'react'

interface EditorButtonProps extends Omit<ButtonHTMLAttributes<HTMLButtonElement>, 'children'> {
  label: string
  active?: boolean
  children: ReactNode
}

export const EditorButton = forwardRef<HTMLButtonElement, EditorButtonProps>(
  ({
    label,
    active = false,
    children,
    className = '',
    ...props
  }, ref) => {
    return (
      <button
        ref={ref}
        type="button"
        aria-label={label}
        title={label}
        aria-pressed={active}
        className={`inline-flex size-9 items-center justify-center rounded-md disabled:opacity-40 ${
          active ? 'bg-black text-white dark:bg-white dark:text-black' : ''
        } ${className}`}
        {...props}
      >
        {children}
      </button>
    )
  },
)

EditorButton.displayName = 'EditorButton'
