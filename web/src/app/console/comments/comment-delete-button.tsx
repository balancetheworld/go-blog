'use client'

import * as Dialog from '@radix-ui/react-dialog'
import { Trash2 } from 'lucide-react'
import { useTranslations } from 'next-intl'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { deleteAdminComment } from '@/api/comment'

interface CommentDeleteButtonProps {
  commentID: number
}

export function CommentDeleteButton({
  commentID,
}: CommentDeleteButtonProps) {
  const t = useTranslations('Console.adminComments')
  const router = useRouter()
  const [open, setOpen] = useState(false)
  const [deleting, setDeleting] = useState(false)
  async function handleDelete() {
    if (deleting)
      return

    setDeleting(true)

    try {
      await deleteAdminComment(commentID)
      toast.success(t('deleted'))
      setOpen(false)
      router.refresh()
    }
    catch {
      toast.error(t('deleteFailed'))
    }
    finally {
      setDeleting(false)
    }
  }
  return (
    <Dialog.Root open={open} onOpenChange={setOpen}>
      <Dialog.Trigger asChild>
        <button
          type="button"
          aria-label={t('deleteAction')}
          title={t('deleteAction')}
          className="inline-flex size-9 items-center justify-center"
        >
          <Trash2 className="size-4" aria-hidden="true" />
        </button>
      </Dialog.Trigger>

      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-40 bg-black/45" />
        <Dialog.Content className="fixed top-1/2 left-1/2 z-50 w-[min(92vw,420px)] -translate-x-1/2 -translate-y-1/2 rounded-md bg-white p-6 shadow-xl dark:bg-neutral-950">
          <Dialog.Title className="text-lg font-semibold">
            {t('deleteTitle')}
          </Dialog.Title>
          <Dialog.Description className="mt-2 text-sm text-neutral-500">
            {t('deleteDescription')}
          </Dialog.Description>

          <div className="mt-6 flex justify-end gap-3">
            <Dialog.Close asChild>
              <button
                type="button"
                disabled={deleting}
                className="min-h-10 rounded-md border border-black/15 px-4 text-sm disabled:opacity-50 dark:border-white/15"
              >
                {t('cancel')}
              </button>
            </Dialog.Close>
            <button
              type="button"
              disabled={deleting}
              onClick={() => void handleDelete()}
              className="min-h-10 rounded-md bg-red-600 px-4 text-sm text-white disabled:opacity-50"
            >
              {deleting ? t('deleting') : t('delete')}
            </button>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}
