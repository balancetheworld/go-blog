'use client'

import type { Editor } from '@tiptap/react'
import type { ChangeEvent } from 'react'
import { useTranslations } from 'next-intl'
import { useRef, useState } from 'react'
import { uploadImage } from '@/api/file'
import { AlignCenter, AlignLeft, AlignRight, ImagePlus, Upload } from '@/components/tiptap-icons'
import { EditorButton } from '@/components/tiptap-ui-primitive/editor-button'
import { EditorPopover } from '@/components/tiptap-ui-primitive/editor-popover'
import { isSafeImageURL } from '@/lib/tiptap-advanced-utils'

interface ImageControlsProps {
  editor: Editor | null
  disabled?: boolean
}

export function ImageControls({ editor, disabled = false }: ImageControlsProps) {
  const t = useTranslations('Console.editor')
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [url, setURL] = useState('')
  const [error, setError] = useState('')
  const [uploading, setUploading] = useState(false)
  const imageSelected = editor?.isActive('image') ?? false

  function handleInsert() {
    if (!isSafeImageURL(url)) {
      setError(t('invalidImageUrl'))
      return
    }

    editor?.chain().focus().setImage({ src: url.trim() }).run()
    setURL('')
    setError('')
  }

  function setAlign(align: 'left' | 'center' | 'right') {
    editor?.chain().focus().updateAttributes('image', { align }).run()
  }

  async function handleUpload(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0]
    event.target.value = ''

    if (!file)
      return

    if (!['image/gif', 'image/jpeg', 'image/png', 'image/webp'].includes(file.type)) {
      setError(t('imageType'))
      return
    }
    if (file.size > 10 * 1024 * 1024) {
      setError(t('imageSize'))
      return
    }

    setUploading(true)
    setError('')

    try {
      const image = await uploadImage(file)
      editor?.chain().focus().setImage({ src: image.url, alt: file.name }).run()
    }
    catch {
      setError(t('imageUploadFailed'))
    }
    finally {
      setUploading(false)
    }
  }

  return (
    <>
      <EditorPopover
        ariaLabel={t('insertImage')}
        trigger={(
          <EditorButton
            label={t('insertImage')}
            disabled={!editor || disabled}
            onClick={() => {}}
          >
            <ImagePlus className="size-4" aria-hidden="true" />
          </EditorButton>
        )}
      >
        <div className="w-72 space-y-2">
          <input
            ref={fileInputRef}
            type="file"
            accept="image/jpeg,image/png,image/gif,image/webp"
            onChange={event => void handleUpload(event)}
            className="sr-only"
          />
          <button
            type="button"
            disabled={uploading}
            onClick={() => fileInputRef.current?.click()}
            className="inline-flex min-h-9 w-full items-center justify-center gap-2 rounded-sm border border-black/15 px-3 text-sm disabled:opacity-50 dark:border-white/15"
          >
            <Upload className="size-4" aria-hidden="true" />
            {uploading ? t('uploading') : t('uploadLocal')}
          </button>
          <div className="flex items-center gap-2 text-xs text-neutral-500">
            <span className="h-px flex-1 bg-black/10 dark:bg-white/10" />
            {t('orUseUrl')}
            <span className="h-px flex-1 bg-black/10 dark:bg-white/10" />
          </div>
          <label className="block space-y-1">
            <span className="text-xs text-neutral-600 dark:text-neutral-400">{t('imageUrl')}</span>
            <input
              type="text"
              value={url}
              onChange={(event) => {
                setURL(event.target.value)
                setError('')
              }}
              placeholder="https://example.com/image.jpg"
              className="min-h-9 w-full rounded-sm border border-black/15 px-2 text-sm dark:border-white/15"
            />
          </label>
          {error && <p className="text-xs text-red-600">{error}</p>}
          <button
            type="button"
            onClick={handleInsert}
            className="min-h-9 w-full rounded-sm bg-black px-3 text-sm text-white dark:bg-white dark:text-black"
          >
            {t('insert')}
          </button>
        </div>
      </EditorPopover>
      {imageSelected && (
        <>
          <EditorButton
            label={t('imageLeft')}
            active={editor?.getAttributes('image').align === 'left'}
            disabled={disabled}
            onClick={() => setAlign('left')}
          >
            <AlignLeft className="size-4" aria-hidden="true" />
          </EditorButton>
          <EditorButton
            label={t('imageCenter')}
            active={editor?.getAttributes('image').align === 'center'}
            disabled={disabled}
            onClick={() => setAlign('center')}
          >
            <AlignCenter className="size-4" aria-hidden="true" />
          </EditorButton>
          <EditorButton
            label={t('imageRight')}
            active={editor?.getAttributes('image').align === 'right'}
            disabled={disabled}
            onClick={() => setAlign('right')}
          >
            <AlignRight className="size-4" aria-hidden="true" />
          </EditorButton>
        </>
      )}
    </>
  )
}
