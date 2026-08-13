'use client'

import type {
  Diary,
  DiaryFolder,
  ListDiariesResponse,
} from '@/models/diary'
import * as Dialog from '@radix-ui/react-dialog'
import { Pencil, Plus, Trash2 } from 'lucide-react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { deleteDiary } from '@/api/diary'

interface DiaryManageProps {
  result: ListDiariesResponse
  folders: DiaryFolder[]
}

export function DiaryManage({ result, folders }: DiaryManageProps) {
  const router = useRouter()
  const [pendingDelete, setPendingDelete] = useState<Diary | null>(null)
  const [deleting, setDeleting] = useState(false)

  async function handleDelete() {
    if (!pendingDelete || deleting)
      return

    setDeleting(true)
    try {
      await deleteDiary(pendingDelete.id)
      setPendingDelete(null)
      toast.success('日记已删除')
      router.refresh()
    }
    catch {
      toast.error('删除日记失败')
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
        <header className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <h1 id="diaries-title" className="text-2xl font-semibold">日记管理</h1>
            <p className="mt-1 text-sm text-neutral-500">
              共
              {' '}
              {result.total}
              {' '}
              篇日记，
              {folders.length}
              {' '}
              个文件夹
            </p>
          </div>
          <Link
            href="/console/diaries/edit/new"
            className="inline-flex min-h-10 items-center gap-2 rounded-md bg-black px-4 text-sm text-white dark:bg-white dark:text-black"
          >
            <Plus className="size-4" aria-hidden="true" />
            新建日记
          </Link>
        </header>

        <div className="overflow-x-auto border-y border-black/10 dark:border-white/10">
          <table className="w-full min-w-[860px] text-left text-sm">
            <thead className="text-neutral-500">
              <tr>
                <th className="px-3 py-3">标题</th>
                <th className="px-3 py-3">文件夹</th>
                <th className="px-3 py-3">状态</th>
                <th className="px-3 py-3">可见范围</th>
                <th className="px-3 py-3">时间</th>
                <th className="px-3 py-3">操作</th>
              </tr>
            </thead>
            <tbody>
              {result.items.map(diary => (
                <tr key={diary.id} className="border-t border-black/10 dark:border-white/10">
                  <td className="max-w-sm px-3 py-4">
                    <p className="truncate font-medium">{diary.title || '无标题日记'}</p>
                    {diary.description && (
                      <p className="mt-1 truncate text-xs text-neutral-500">
                        {diary.description}
                      </p>
                    )}
                  </td>
                  <td className="px-3 py-4">{diary.folder?.name ?? '未分类'}</td>
                  <td className="px-3 py-4">
                    {diary.status === 'published' ? '已发布' : '草稿'}
                  </td>
                  <td className="px-3 py-4">
                    {diary.visibility === 'private'
                      ? '私密'
                      : diary.visibility === 'roles'
                        ? diary.visibleRoles.map(role => role.name).join('、')
                        : '公开'}
                  </td>
                  <td className="px-3 py-4">
                    {new Date(diary.updatedAt).toLocaleString('zh-CN')}
                  </td>
                  <td className="px-3 py-4">
                    <Link
                      href={`/console/diaries/edit/${diary.id}`}
                      aria-label="编辑日记"
                      title="编辑日记"
                      className="inline-flex size-9 items-center justify-center"
                    >
                      <Pencil className="size-4" aria-hidden="true" />
                    </Link>
                    <button
                      type="button"
                      onClick={() => setPendingDelete(diary)}
                      aria-label="删除日记"
                      title="删除日记"
                      className="inline-flex size-9 items-center justify-center text-red-600"
                    >
                      <Trash2 className="size-4" aria-hidden="true" />
                    </button>
                  </td>
                </tr>
              ))}
              {result.items.length === 0 && (
                <tr>
                  <td colSpan={6} className="px-3 py-12 text-center text-neutral-500">
                    暂无日记
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
          <Dialog.Title className="text-lg font-semibold">删除日记</Dialog.Title>
          <Dialog.Description className="mt-2 text-sm text-neutral-500">
            删除后无法恢复，相关评论也会被删除。
          </Dialog.Description>
          <div className="mt-6 flex justify-end gap-3">
            <Dialog.Close asChild>
              <button
                type="button"
                disabled={deleting}
                className="min-h-10 rounded-md border border-black/15 px-4 text-sm disabled:opacity-50 dark:border-white/15"
              >
                取消
              </button>
            </Dialog.Close>
            <button
              type="button"
              disabled={deleting}
              onClick={() => void handleDelete()}
              className="min-h-10 rounded-md bg-red-600 px-4 text-sm text-white disabled:opacity-50"
            >
              {deleting ? '删除中' : '删除'}
            </button>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}
