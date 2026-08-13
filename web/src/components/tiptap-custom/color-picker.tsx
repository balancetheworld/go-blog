'use client'

import type { Editor } from '@tiptap/react'
import { Palette, RemoveFormatting } from '@/components/tiptap-icons'
import { EditorButton } from '@/components/tiptap-ui-primitive/editor-button'
import { EditorPopover } from '@/components/tiptap-ui-primitive/editor-popover'

const colors = ['#171717', '#dc2626', '#ea580c', '#ca8a04', '#16a34a', '#2563eb', '#9333ea']

interface ColorPickerProps {
  editor: Editor | null
  disabled?: boolean
}

export function ColorPicker({ editor, disabled = false }: ColorPickerProps) {
  return (
    <EditorPopover
      ariaLabel="文字颜色"
      trigger={(
        <EditorButton
          label="文字颜色"
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
            aria-label={`设置文字颜色 ${color}`}
            title={color}
            onClick={() => editor?.chain().focus().setColor(color).run()}
            className="size-7 rounded-sm border border-black/15 dark:border-white/15"
            style={{ backgroundColor: color }}
          />
        ))}
        <button
          type="button"
          aria-label="清除文字颜色"
          title="清除文字颜色"
          onClick={() => editor?.chain().focus().unsetColor().run()}
          className="inline-flex size-7 items-center justify-center rounded-sm border border-black/15 dark:border-white/15"
        >
          <RemoveFormatting className="size-4" aria-hidden="true" />
        </button>
      </div>
    </EditorPopover>
  )
}
