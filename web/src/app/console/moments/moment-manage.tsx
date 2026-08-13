'use client'

import type { ListMomentsResponse, Moment } from '@/models/moment'
import * as Dialog from '@radix-ui/react-dialog'
import { Send, Trash2 } from 'lucide-react'
import { useLocale, useTranslations } from 'next-intl'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { createMoment, deleteMoment } from '@/api/moment'

interface MomentManageProps {
  result: ListMomentsResponse
  currentUserId: number
  isAdmin: boolean
}

export function MomentManage({
  result,
  currentUserId,
  isAdmin,
}: MomentManageProps) {
  const locale = useLocale()
  const t = useTranslations('Console.moments')
  const router = useRouter()
  const [content, setContent] = useState('')
  const [creating, setCreating] = useState(false)
  const [pendingDelete, setPendingDelete] = useState<Moment | null>(null)
  const [deleting, setDeleting] = useState(false)

  async function handleCreate(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const value = content.trim()
    if (!value) {
      toast.error(t('contentRequired'))
      return
    }

    setCreating(true)
    try {
      await createMoment({ content: value })
      setContent('')
      toast.success(t('created'))
      router.refresh()
    }
    catch {
      toast.error(t('createFailed'))
    }
    finally {
      setCreating(false)
    }
  }

  async function handleDelete() {
    if (!pendingDelete || deleting)
      return

    setDeleting(true)
    try {
      await deleteMoment(pendingDelete.id)
      setPendingDelete(null)
      toast.success(t('deleted'))
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
    <Dialog.Root
      open={pendingDelete !== null}
      onOpenChange={(open) => {
        if (!open && !deleting)
          setPendingDelete(null)
      }}
    >
      <div className="space-y-6">
        <header>
          <h1 id="moments-title" className="text-2xl font-semibold">{t('title')}</h1>
          <p className="mt-1 text-sm text-neutral-500">
            {t('count', { count: result.total })}
          </p>
        </header>

        <form onSubmit={handleCreate} className="space-y-3 border-y py-5">
          <label className="block space-y-2">
            <span className="text-sm">{t('content')}</span>
            <textarea
              value={content}
              onChange={event => setContent(event.target.value)}
              maxLength={2000}
              rows={5}
              placeholder={t('placeholder')}
              className="w-full resize-y rounded-md border border-black/15 px-3 py-2 dark:border-white/15"
            />
          </label>
          <div className="flex items-center justify-between gap-4">
            <span className="text-xs text-neutral-500">
              {t('characterCount', { count: content.length })}
            </span>
            <button
              type="submit"
              disabled={creating || !content.trim()}
              className="inline-flex min-h-10 items-center gap-2 rounded-md bg-black px-4 text-sm text-white disabled:opacity-50 dark:bg-white dark:text-black"
            >
              <Send className="size-4" aria-hidden="true" />
              {creating ? t('creating') : t('create')}
            </button>
          </div>
        </form>

        <div className="overflow-x-auto border-y border-black/10 dark:border-white/10">
          <table className="w-full min-w-[720px] text-left text-sm">
            <thead className="text-neutral-500">
              <tr>
                <th className="px-3 py-3">{t('content')}</th>
                <th className="px-3 py-3">{t('author')}</th>
                <th className="px-3 py-3">{t('time')}</th>
                <th className="px-3 py-3">{t('actions')}</th>
              </tr>
            </thead>
            <tbody>
              {result.items.map(moment => (
                <tr key={moment.id} className="border-t border-black/10 dark:border-white/10">
                  <td className="max-w-xl px-3 py-4">
                    <p className="line-clamp-4 whitespace-pre-wrap break-words">{moment.content}</p>
                  </td>
                  <td className="px-3 py-4">
                    {moment.author.nickname || moment.author.username}
                  </td>
                  <td className="px-3 py-4">
                    {new Date(moment.createdAt).toLocaleString(locale)}
                  </td>
                  <td className="px-3 py-4">
                    {(isAdmin || moment.author.id === currentUserId) && (
                      <button
                        type="button"
                        onClick={() => setPendingDelete(moment)}
                        aria-label={t('deleteAction')}
                        title={t('deleteAction')}
                        className="inline-flex size-9 items-center justify-center text-red-600"
                      >
                        <Trash2 className="size-4" aria-hidden="true" />
                      </button>
                    )}
                  </td>
                </tr>
              ))}
              {result.items.length === 0 && (
                <tr>
                  <td colSpan={4} className="px-3 py-12 text-center text-neutral-500">
                    {t('empty')}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-40 bg-black/45" />
        <Dialog.Content className="fixed top-1/2 left-1/2 z-50 w-[min(92vw,420px)] -translate-x-1/2 -translate-y-1/2 rounded-md bg-white p-6 shadow-xl dark:bg-neutral-950">
          <Dialog.Title className="text-lg font-semibold">{t('deleteTitle')}</Dialog.Title>
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
