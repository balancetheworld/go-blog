'use client'

import type { Editor } from '@tiptap/react'
import { useTranslations } from 'next-intl'
import { Braces, Code } from '@/components/tiptap-icons'
import { EditorButton } from '@/components/tiptap-ui-primitive/editor-button'

const languages = [
  { label: '', value: '' },
  { label: 'JavaScript', value: 'javascript' },
  { label: 'TypeScript', value: 'typescript' },
  { label: 'JSX', value: 'jsx' },
  { label: 'TSX', value: 'tsx' },
  { label: 'Go', value: 'go' },
  { label: 'Python', value: 'python' },
  { label: 'Java', value: 'java' },
  { label: 'C', value: 'c' },
  { label: 'C++', value: 'cpp' },
  { label: 'C#', value: 'csharp' },
  { label: 'Rust', value: 'rust' },
  { label: 'Bash', value: 'bash' },
  { label: 'SQL', value: 'sql' },
  { label: 'JSON', value: 'json' },
  { label: 'HTML', value: 'xml' },
  { label: 'CSS', value: 'css' },
  { label: 'Markdown', value: 'markdown' },
]

interface CodeControlsProps {
  editor: Editor | null
  disabled?: boolean
  inlineCode?: boolean
  codeBlock?: boolean
  language?: string
}

export function CodeControls({
  editor,
  disabled = false,
  inlineCode = false,
  codeBlock = false,
  language = '',
}: CodeControlsProps) {
  const t = useTranslations('Console.editor')

  function changeLanguage(value: string) {
    if (!editor)
      return

    if (editor.isActive('codeBlock')) {
      editor.chain().focus().updateAttributes('codeBlock', { language: value || null }).run()
      return
    }

    editor.chain().focus().setCodeBlock({ language: value }).run()
  }

  return (
    <div className="flex items-center gap-1">
      <EditorButton
        label={t('inlineCode')}
        active={inlineCode}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().toggleCode().run()}
      >
        <Code className="size-4" aria-hidden="true" />
      </EditorButton>
      <EditorButton
        label={t('codeBlock')}
        active={codeBlock}
        disabled={!editor || disabled}
        onClick={() => editor?.chain().focus().toggleCodeBlock({ language }).run()}
      >
        <Braces className="size-4" aria-hidden="true" />
      </EditorButton>
      <select
        aria-label={t('codeLanguage')}
        title={t('codeLanguage')}
        value={language}
        disabled={!editor || disabled}
        onChange={event => changeLanguage(event.target.value)}
        className="h-9 max-w-32 rounded-md border border-black/15 bg-transparent px-2 text-xs disabled:opacity-40 dark:border-white/15"
      >
        {languages.map(item => (
          <option key={item.value} value={item.value}>{item.value ? item.label : t('autoLanguage')}</option>
        ))}
      </select>
    </div>
  )
}
