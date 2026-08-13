'use client'

import type { Editor } from '@tiptap/react'
import { useEditorState } from '@tiptap/react'
import { CodeControls } from '@/components/tiptap-custom/code-controls'
import { ColorPicker } from '@/components/tiptap-custom/color-picker'
import { ImageControls } from '@/components/tiptap-custom/image-controls'
import {
  AlignCenter,
  AlignLeft,
  AlignRight,
  Bold,
  Heading1,
  Heading2,
  Highlighter,
  Italic,
  List,
  ListOrdered,
  Maximize2,
  Minimize2,
  Minus,
  Quote,
  Redo2,
  RemoveFormatting,
  Subscript,
  Superscript,
  Undo2,
} from '@/components/tiptap-icons'
import { EditorButton } from '@/components/tiptap-ui-primitive/editor-button'

interface EditorToolbarProps {
  editor: Editor | null
  disabled?: boolean
  fullscreen?: boolean
  onToggleFullscreen?: () => void
}

export function EditorToolbar({
  editor,
  disabled = false,
  fullscreen = false,
  onToggleFullscreen,
}: EditorToolbarProps) {
  const state = useEditorState({
    editor,
    selector: ({ editor }) => editor
      ? {
          bold: editor.isActive('bold'),
          inlineCode: editor.isActive('code'),
          italic: editor.isActive('italic'),
          heading1: editor.isActive('heading', { level: 1 }),
          heading2: editor.isActive('heading', { level: 2 }),
          bulletList: editor.isActive('bulletList'),
          orderedList: editor.isActive('orderedList'),
          blockquote: editor.isActive('blockquote'),
          codeBlock: editor.isActive('codeBlock'),
          codeLanguage: editor.getAttributes('codeBlock').language ?? '',
          highlight: editor.isActive('highlight'),
          subscript: editor.isActive('subscript'),
          superscript: editor.isActive('superscript'),
          alignLeft: editor.isActive({ textAlign: 'left' }),
          alignCenter: editor.isActive({ textAlign: 'center' }),
          alignRight: editor.isActive({ textAlign: 'right' }),
          selectionEmpty: editor.state.selection.empty,
          canUndo: editor.can().chain().focus().undo().run(),
          canRedo: editor.can().chain().focus().redo().run(),
        }
      : null,
  })

  return (
    <div
      role="toolbar"
      aria-label="正文格式"
      className="flex flex-wrap items-center gap-1 border-b border-black/15 p-2 dark:border-white/15"
    >
      <EditorButton
        label="加粗"
        active={state?.bold}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().toggleBold().run()}
      >
        <Bold className="size-4" aria-hidden="true" />
      </EditorButton>
      <EditorButton
        label="斜体"
        active={state?.italic}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().toggleItalic().run()}
      >
        <Italic className="size-4" aria-hidden="true" />
      </EditorButton>

      <span className="mx-1 h-6 w-px bg-black/10 dark:bg-white/10" />

      <EditorButton
        label="一级标题"
        active={state?.heading1}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().toggleHeading({ level: 1 }).run()}
      >
        <Heading1 className="size-4" aria-hidden="true" />
      </EditorButton>
      <EditorButton
        label="二级标题"
        active={state?.heading2}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().toggleHeading({ level: 2 }).run()}
      >
        <Heading2 className="size-4" aria-hidden="true" />
      </EditorButton>
      <EditorButton
        label="无序列表"
        active={state?.bulletList}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().toggleBulletList().run()}
      >
        <List className="size-4" aria-hidden="true" />
      </EditorButton>
      <EditorButton
        label="有序列表"
        active={state?.orderedList}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().toggleOrderedList().run()}
      >
        <ListOrdered className="size-4" aria-hidden="true" />
      </EditorButton>
      <EditorButton
        label="引用"
        active={state?.blockquote}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().toggleBlockquote().run()}
      >
        <Quote className="size-4" aria-hidden="true" />
      </EditorButton>
      <CodeControls
        editor={editor}
        disabled={disabled}
        inlineCode={state?.inlineCode}
        codeBlock={state?.codeBlock}
        language={state?.codeLanguage}
      />
      <EditorButton
        label="高亮"
        active={state?.highlight}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().toggleHighlight().run()}
      >
        <Highlighter className="size-4" aria-hidden="true" />
      </EditorButton>
      <EditorButton
        label="下标"
        active={state?.subscript}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().toggleSubscript().run()}
      >
        <Subscript className="size-4" aria-hidden="true" />
      </EditorButton>
      <EditorButton
        label="上标"
        active={state?.superscript}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().toggleSuperscript().run()}
      >
        <Superscript className="size-4" aria-hidden="true" />
      </EditorButton>

      <span className="mx-1 h-6 w-px bg-black/10 dark:bg-white/10" />

      <EditorButton
        label="左对齐"
        active={state?.alignLeft}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().setTextAlign('left').run()}
      >
        <AlignLeft className="size-4" aria-hidden="true" />
      </EditorButton>
      <EditorButton
        label="居中"
        active={state?.alignCenter}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().setTextAlign('center').run()}
      >
        <AlignCenter className="size-4" aria-hidden="true" />
      </EditorButton>
      <EditorButton
        label="右对齐"
        active={state?.alignRight}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().setTextAlign('right').run()}
      >
        <AlignRight className="size-4" aria-hidden="true" />
      </EditorButton>
      <EditorButton
        label="分隔线"
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().setHorizontalRule().run()}
      >
        <Minus className="size-4" aria-hidden="true" />
      </EditorButton>
      <ColorPicker editor={editor} disabled={disabled} />
      <ImageControls editor={editor} disabled={disabled} />
      <EditorButton
        label="清除格式"
        disabled={!editor || disabled || state?.selectionEmpty}
        onClick={() => editor
          ?.chain()
          .focus()
          .unsetAllMarks()
          .clearNodes()
          .unsetTextAlign()
          .run()}
      >
        <RemoveFormatting className="size-4" aria-hidden="true" />
      </EditorButton>

      <span className="mx-1 h-6 w-px bg-black/10 dark:bg-white/10" />

      <EditorButton
        label="撤销"
        disabled={!editor || disabled || !state?.canUndo}
        onClick={() => editor?.chain().focus().undo().run()}
      >
        <Undo2 className="size-4" aria-hidden="true" />
      </EditorButton>
      <EditorButton
        label="重做"
        disabled={!editor || disabled || !state?.canRedo}
        onClick={() => editor?.chain().focus().redo().run()}
      >
        <Redo2 className="size-4" aria-hidden="true" />
      </EditorButton>

      {onToggleFullscreen && (
        <EditorButton
          label={fullscreen ? '退出全屏' : '全屏编辑'}
          active={fullscreen}
          disabled={!editor || disabled}
          onClick={onToggleFullscreen}
        >
          {fullscreen
            ? <Minimize2 className="size-4" aria-hidden="true" />
            : <Maximize2 className="size-4" aria-hidden="true" />}
        </EditorButton>
      )}
    </div>
  )
}
