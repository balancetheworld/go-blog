'use client'

import type { Editor } from '@tiptap/react'
import type { ChangeEvent } from 'react'
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
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [url, setURL] = useState('')
  const [error, setError] = useState('')
  const [uploading, setUploading] = useState(false)
  const imageSelected = editor?.isActive('image') ?? false

  function handleInsert() {
    if (!isSafeImageURL(url)) {
      setError('请输入有效的 HTTP(S) 或站内图片地址')
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
      setError('仅支持 JPEG、PNG、GIF 和 WebP 图片')
      return
    }
    if (file.size > 10 * 1024 * 1024) {
      setError('图片不能超过 10 MB')
      return
    }

    setUploading(true)
    setError('')

    try {
      const image = await uploadImage(file)
      editor?.chain().focus().setImage({ src: image.url, alt: file.name }).run()
    }
    catch {
      setError('图片上传失败')
    }
    finally {
      setUploading(false)
    }
  }

  return (
    <>
      <EditorPopover
        ariaLabel="插入图片"
        trigger={(
          <EditorButton
            label="插入图片"
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
            {uploading ? '上传中' : '上传本地图片'}
          </button>
          <div className="flex items-center gap-2 text-xs text-neutral-500">
            <span className="h-px flex-1 bg-black/10 dark:bg-white/10" />
            或使用图片地址
            <span className="h-px flex-1 bg-black/10 dark:bg-white/10" />
          </div>
          <label className="block space-y-1">
            <span className="text-xs text-neutral-600 dark:text-neutral-400">图片地址</span>
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
            插入
          </button>
        </div>
      </EditorPopover>
      {imageSelected && (
        <>
          <EditorButton
            label="图片左对齐"
            active={editor?.getAttributes('image').align === 'left'}
            disabled={disabled}
            onClick={() => setAlign('left')}
          >
            <AlignLeft className="size-4" aria-hidden="true" />
          </EditorButton>
          <EditorButton
            label="图片居中"
            active={editor?.getAttributes('image').align === 'center'}
            disabled={disabled}
            onClick={() => setAlign('center')}
          >
            <AlignCenter className="size-4" aria-hidden="true" />
          </EditorButton>
          <EditorButton
            label="图片右对齐"
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
