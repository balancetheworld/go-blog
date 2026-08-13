import type { Editor } from '@tiptap/react'
import Color from '@tiptap/extension-color'
import Highlight from '@tiptap/extension-highlight'
import Subscript from '@tiptap/extension-subscript'
import Superscript from '@tiptap/extension-superscript'
import TextAlign from '@tiptap/extension-text-align'
import { TextStyle } from '@tiptap/extension-text-style'
import Typography from '@tiptap/extension-typography'
import { useEditor } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import { useTranslations } from 'next-intl'
import { useEffect, useRef } from 'react'
import { toast } from 'sonner'
import { uploadImage } from '@/api/file'
import { AlignedImage } from '@/components/tiptap-node/aligned-image'
import { CodeBlockNode } from '@/components/tiptap-node/code-block'

interface UseTiptapEditorOptions {
  content: string
  editable?: boolean
  onChange?: (html: string, editor: Editor) => void
}

export function useTiptapEditor({
  content,
  editable = true,
  onChange,
}: UseTiptapEditorOptions) {
  const t = useTranslations('Console.editor')
  const editorRef = useRef<Editor | null>(null)
  const editor = useEditor({
    extensions: [
      StarterKit.configure({ codeBlock: false }),
      AlignedImage,
      CodeBlockNode,
      Color,
      Highlight.configure({ multicolor: true }),
      Subscript,
      Superscript,
      TextAlign.configure({ types: ['heading', 'paragraph'] }),
      TextStyle,
      Typography,
    ],
    content,
    editable,
    immediatelyRender: false,
    editorProps: {
      handlePaste(_view, event) {
        const image = Array.from(event.clipboardData?.files ?? [])
          .find(file => file.type.startsWith('image/'))

        if (!image)
          return false

        event.preventDefault()
        void uploadImage(image)
          .then((uploadedImage) => {
            editorRef.current
              ?.chain()
              .focus()
              .setImage({ src: uploadedImage.url, alt: image.name })
              .run()
          })
          .catch(() => toast.error(t('pasteUploadFailed')))

        return true
      },
    },
    onUpdate({ editor }) {
      onChange?.(editor.getHTML(), editor)
    },
  })

  useEffect(() => {
    editorRef.current = editor
    return () => {
      editorRef.current = null
    }
  }, [editor])

  useEffect(() => {
    if (!editor || editor.getHTML() === content)
      return

    editor.commands.setContent(content, { emitUpdate: false })
  }, [content, editor])

  useEffect(() => {
    if (!editor || editor.isEditable === editable)
      return

    editor.setEditable(editable, false)
  }, [editable, editor])

  return editor
}
