'use client'

import type { NodeViewProps } from '@tiptap/react'
import type { PointerEvent as ReactPointerEvent } from 'react'
import { NodeViewWrapper } from '@tiptap/react'
import { useEffect, useRef } from 'react'

export function ResizableImageView({
  editor,
  node,
  selected,
  updateAttributes,
}: NodeViewProps) {
  const imageRef = useRef<HTMLImageElement>(null)
  const cleanupRef = useRef<(() => void) | null>(null)
  const align = node.attrs.align ?? 'center'
  const width = typeof node.attrs.width === 'number' ? node.attrs.width : null

  useEffect(() => () => cleanupRef.current?.(), [])

  function startResize(event: ReactPointerEvent<HTMLButtonElement>, direction: -1 | 1) {
    if (!editor.isEditable || !imageRef.current)
      return

    event.preventDefault()
    event.stopPropagation()

    cleanupRef.current?.()

    const startX = event.clientX
    const startWidth = imageRef.current.getBoundingClientRect().width
    const maxWidth = Math.max(120, editor.view.dom.clientWidth - 32)

    function handlePointerMove(pointerEvent: PointerEvent) {
      const nextWidth = Math.round(Math.min(maxWidth, Math.max(120, startWidth + (pointerEvent.clientX - startX) * direction)))
      updateAttributes({ width: nextWidth })
    }

    function cleanup() {
      window.removeEventListener('pointermove', handlePointerMove)
      window.removeEventListener('pointerup', cleanup)
      window.removeEventListener('pointercancel', cleanup)
      cleanupRef.current = null
    }

    cleanupRef.current = cleanup
    window.addEventListener('pointermove', handlePointerMove)
    window.addEventListener('pointerup', cleanup)
    window.addEventListener('pointercancel', cleanup)
  }

  return (
    <NodeViewWrapper
      className="relative max-w-full leading-none"
      style={{
        marginLeft: align === 'right' || align === 'center' ? 'auto' : undefined,
        marginRight: align === 'left' || align === 'center' ? 'auto' : undefined,
        width: width ? `${width}px` : 'fit-content',
      }}
    >
      <img
        ref={imageRef}
        src={node.attrs.src}
        alt={node.attrs.alt ?? ''}
        title={node.attrs.title ?? undefined}
        draggable={false}
        className={`block h-auto max-w-full ${selected ? 'outline-2 outline-black/50 dark:outline-white/50' : ''}`}
        style={{ width: width ? '100%' : undefined }}
      />
      {selected && editor.isEditable && (
        <>
          <button
            type="button"
            aria-label="从左侧缩放图片"
            title="拖动缩放图片"
            contentEditable={false}
            onPointerDown={event => startResize(event, -1)}
            className="absolute -bottom-1.5 -left-1.5 size-3 cursor-nwse-resize rounded-full border border-white bg-black dark:border-black dark:bg-white"
          />
          <button
            type="button"
            aria-label="从右侧缩放图片"
            title="拖动缩放图片"
            contentEditable={false}
            onPointerDown={event => startResize(event, 1)}
            className="absolute -bottom-1.5 -right-1.5 size-3 cursor-nwse-resize rounded-full border border-white bg-black dark:border-black dark:bg-white"
          />
        </>
      )}
    </NodeViewWrapper>
  )
}
