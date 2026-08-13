'use client'

import type { Editor } from '@tiptap/react'
import { useTranslations } from 'next-intl'
import { Palette, RemoveFormatting } from '@/components/tiptap-icons'
import { EditorButton } from '@/components/tiptap-ui-primitive/editor-button'
import { EditorPopover } from '@/components/tiptap-ui-primitive/editor-popover'

const colors = ['#171717', '#dc2626', '#ea580c', '#ca8a04', '#16a34a', '#2563eb', '#9333ea']

interface ColorPickerProps {
  editor: Editor | null
  disabled?: boolean
}

export function ColorPicker({ editor, disabled = false }: ColorPickerProps) {
  const t = useTranslations('Console.editor')

  return (
    <EditorPopover
      ariaLabel={t('textColor')}
      trigger={(
        <EditorButton
          label={t('textColor')}
          disabled={!editor || disabled}
          onClick={() => {}}
        >
          <Palette className="size-4" aria-hidden="true" />
        </EditorButton>
      )}
    >
      <div className="flex items-center gap-2">
        {colors.map(color => (
          <button
            key={color}
            type="button"
            aria-label={t('setTextColor', { color })}
            title={color}
            onClick={() => editor?.chain().focus().setColor(color).run()}
            className="size-7 rounded-sm border border-black/15 dark:border-white/15"
            style={{ backgroundColor: color }}
          />
        ))}
        <button
          type="button"
          aria-label={t('clearTextColor')}
          title={t('clearTextColor')}
          onClick={() => editor?.chain().focus().unsetColor().run()}
          className="inline-flex size-7 items-center justify-center rounded-sm border border-black/15 dark:border-white/15"
        >
          <RemoveFormatting className="size-4" aria-hidden="true" />
        </button>
      </div>
    </EditorPopover>
  )
}
