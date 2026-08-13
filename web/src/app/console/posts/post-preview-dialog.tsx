'use client'

import type {
  Category,
  Label,
  PostVisibility,
} from '@/models/post'
import type { RoleOption } from '@/models/role'
import * as Dialog from '@radix-ui/react-dialog'
import { X } from 'lucide-react'
import Image from 'next/image'
import { SimpleEditor } from '@/components/tiptap-editor/simple-editor'

export interface PostPreview {
  title: string
  description: string
  cover: string
  category?: Category
  labels: Label[]
  visibility: PostVisibility
  visibleRoles: RoleOption[]
  content: string
}

interface PostPreviewDialogProps {
  preview: PostPreview | null
  onClose: () => void
}

export function PostPreviewDialog({
  preview,
  onClose,
}: PostPreviewDialogProps) {
  return (
    <Dialog.Root
      open={preview !== null}
      onOpenChange={(open) => {
        if (!open)
          onClose()
      }}
    >
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-40 bg-black/45" />
        <Dialog.Content className="fixed inset-x-4 top-4 bottom-4 z-50 overflow-y-auto rounded-md bg-white shadow-xl sm:inset-x-8 lg:inset-x-[max(2rem,calc((100vw-960px)/2))] dark:bg-neutral-950">
          <div className="sticky top-0 z-10 flex min-h-14 items-center justify-between border-b border-black/10 bg-white px-5 dark:border-white/10 dark:bg-neutral-950">
            <Dialog.Title className="font-medium">文章预览</Dialog.Title>
            <Dialog.Close asChild>
              <button
                type="button"
                aria-label="关闭预览"
                title="关闭预览"
                className="inline-flex size-9 items-center justify-center"
              >
                <X className="size-4" aria-hidden="true" />
              </button>
            </Dialog.Close>
          </div>

          {preview && (
            <article className="mx-auto w-full max-w-3xl px-5 py-8 sm:px-8">
              <header className="border-b border-black/10 pb-8 dark:border-white/10">
                <div className="flex flex-wrap items-center gap-3 text-sm text-neutral-500">
                  {preview.category && <span>{preview.category.name}</span>}
                  {preview.visibility === 'private' && <span>仅作者和管理员</span>}
                  {preview.visibility === 'roles' && (
                    <span>{preview.visibleRoles.map(role => role.name).join('、') || '未选择身份'}</span>
                  )}
                </div>

                <h1 className="mt-4 text-3xl font-semibold leading-tight">
                  {preview.title}
                </h1>

                {preview.description && (
                  <p className="mt-4 leading-7 text-neutral-600 dark:text-neutral-400">
                    {preview.description}
                  </p>
                )}

                {preview.cover && (
                  <div className="relative mt-8 aspect-video overflow-hidden rounded-md">
                    <Image
                      src={preview.cover}
                      alt={preview.title}
                      fill
                      sizes="(max-width: 768px) 100vw, 768px"
                      className="object-cover"
                    />
                  </div>
                )}
              </header>

              <div className="py-8">
                <SimpleEditor
                  value={preview.content}
                  disabled
                  showToolbar={false}
                  minHeight={120}
                  ariaLabel="文章预览正文"
                />
              </div>

              {preview.labels.length > 0 && (
                <footer className="flex flex-wrap gap-2 border-t border-black/10 pt-6 dark:border-white/10">
                  {preview.labels.map(label => (
                    <span
                      key={label.id}
                      className="rounded-sm border border-black/10 px-2 py-1 text-sm dark:border-white/10"
                    >
                      {label.name}
                    </span>
                  ))}
                </footer>
              )}
            </article>
          )}
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}
