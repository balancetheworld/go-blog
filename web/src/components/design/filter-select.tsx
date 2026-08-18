'use client'

import * as Select from '@radix-ui/react-select'
import { Check, ChevronDown } from 'lucide-react'
import { useRef, useState } from 'react'

interface FilterSelectOption {
  label: string
  value: string
}

interface FilterSelectProps {
  name: string
  defaultValue: string
  ariaLabel: string
  options: FilterSelectOption[]
  submitOnChange?: boolean
  disabled?: boolean
}

export function FilterSelect({
  name,
  defaultValue,
  ariaLabel,
  options,
  submitOnChange = false,
  disabled = false,
}: FilterSelectProps) {
  const [value, setValue] = useState(defaultValue)
  const triggerRef = useRef<HTMLButtonElement>(null)

  function handleValueChange(nextValue: string) {
    setValue(nextValue)
    if (submitOnChange) {
      requestAnimationFrame(() => {
        triggerRef.current?.closest('form')?.requestSubmit()
      })
    }
  }

  return (
    <>
      {value !== 'all' && <input type="hidden" name={name} value={value} />}
      <Select.Root value={value} onValueChange={handleValueChange}>
        <Select.Trigger ref={triggerRef} className="filter-select-trigger" aria-label={ariaLabel} disabled={disabled}>
          <Select.Value />
          <Select.Icon>
            <ChevronDown aria-hidden="true" />
          </Select.Icon>
        </Select.Trigger>
        <Select.Portal>
          <Select.Content className="filter-select-content" position="popper" sideOffset={6}>
            <Select.Viewport className="filter-select-viewport">
              {options.map(option => (
                <Select.Item key={option.value} value={option.value} className="filter-select-item">
                  <Select.ItemText>{option.label}</Select.ItemText>
                  <Select.ItemIndicator className="filter-select-indicator">
                    <Check aria-hidden="true" />
                  </Select.ItemIndicator>
                </Select.Item>
              ))}
            </Select.Viewport>
          </Select.Content>
        </Select.Portal>
      </Select.Root>
    </>
  )
}
