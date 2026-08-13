import type { Editor } from '@tiptap/react'

export function getEditorHTML(editor: Editor | null): string {
  return editor?.getHTML() ?? ''
}

export function isEditorEmpty(editor: Editor | null): boolean {
  return !editor || editor.isEmpty
}
