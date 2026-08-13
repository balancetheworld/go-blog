'use client'

import type { CSSProperties } from 'react'
import { EditorContent } from '@tiptap/react'
import { useEffect, useState } from 'react'
import { EditorToolbar } from '@/components/tiptap-ui/editor-toolbar'
import { useTiptapEditor } from '@/hooks/use-tiptap-editor'

interface SimpleEditorProps {
  value?: string
  defaultValue?: string
  onChange?: (html: string) => void
  disabled?: boolean
  showToolbar?: boolean
  minHeight?: number
  ariaLabel?: string
}

export function SimpleEditor({
  value,
  defaultValue = '',
  onChange,
  disabled = false,
  showToolbar = true,
  minHeight = 320,
  ariaLabel = '正文编辑器',
}: SimpleEditorProps) {
  const [internalValue, setInternalValue] = useState(defaultValue)
  const [fullscreen, setFullscreen] = useState(false)
  const controlled = value !== undefined
  const content = controlled ? value : internalValue
  const editor = useTiptapEditor({
    content,
    editable: !disabled,
    onChange(html) {
      if (!controlled)
        setInternalValue(html)

      onChange?.(html)
    },
  })
  const style = {
    '--tiptap-min-height': `${minHeight}px`,
  } as CSSProperties

  useEffect(() => {
    if (!fullscreen)
      return

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape')
        setFullscreen(false)
    }

    document.addEventListener('keydown', handleKeyDown)
    return () => document.removeEventListener('keydown', handleKeyDown)
  }, [fullscreen])

  return (
    <div
      aria-label={ariaLabel}
      className={fullscreen
        ? 'fixed inset-0 z-50 flex flex-col overflow-hidden bg-white dark:bg-neutral-950'
        : 'overflow-hidden rounded-md border border-black/15 dark:border-white/15'}
      style={style}
    >
      {showToolbar && (
        <EditorToolbar
          editor={editor}
          disabled={disabled}
          fullscreen={fullscreen}
          onToggleFullscreen={() => setFullscreen(current => !current)}
        />
      )}
      <EditorContent
        editor={editor}
        className={`[&_.ProseMirror]:min-h-[var(--tiptap-min-height)] [&_.ProseMirror]:p-4 [&_.ProseMirror]:outline-none [&_.ProseMirror>_*+*]:mt-3 [&_.ProseMirror_blockquote]:border-l-4 [&_.ProseMirror_blockquote]:border-black/20 [&_.ProseMirror_blockquote]:pl-4 [&_.ProseMirror_h1]:text-3xl [&_.ProseMirror_h1]:font-semibold [&_.ProseMirror_h2]:text-2xl [&_.ProseMirror_h2]:font-semibold [&_.ProseMirror_img]:max-w-full [&_.ProseMirror_ol]:list-decimal [&_.ProseMirror_ol]:pl-6 [&_.ProseMirror_pre]:overflow-x-auto [&_.ProseMirror_pre]:rounded-md [&_.ProseMirror_pre]:bg-neutral-900 [&_.ProseMirror_pre]:p-4 [&_.ProseMirror_pre]:text-neutral-100 [&_.ProseMirror_ul]:list-disc [&_.ProseMirror_ul]:pl-6 [&_.ProseMirror_.ProseMirror-selectednode]:outline-2 [&_.ProseMirror_.ProseMirror-selectednode]:outline-black/40 ${fullscreen ? 'flex-1 overflow-y-auto [&_.ProseMirror]:mx-auto [&_.ProseMirror]:w-full [&_.ProseMirror]:max-w-5xl' : ''}`}
      />
    </div>
  )
}
