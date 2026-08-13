'use client'

import type { FormEvent } from 'react'
import type { Label } from '@/models/post'
import * as Dialog from '@radix-ui/react-dialog'
import { Pencil, Plus, Trash2, X } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { createLabel, deleteLabel, updateLabel } from '@/api/post'

interface LabelManageProps {
  labels: Label[]
}

export function LabelManage({ labels }: LabelManageProps) {
  const router = useRouter()
  const [name, setName] = useState('')
  const [slug, setSlug] = useState('')
  const [saving, setSaving] = useState(false)
  const [editingLabel, setEditingLabel] = useState<Label | null>(null)
  const [pendingDelete, setPendingDelete] = useState<Label | null>(null)
  const [deleting, setDeleting] = useState(false)

  function openEdit(label: Label) {
    setEditingLabel(label)
    setName(label.name)
    setSlug(label.slug)
  }

  function cancelEdit() {
    setEditingLabel(null)
    setName('')
    setSlug('')
  }

  async function handleSave(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const labelName = name.trim()
    const labelSlug = slug.trim()

    if (!labelName) {
      toast.error('请输入标签名称')
      return
    }
    if (!labelSlug) {
      toast.error('请输入标签 Slug')
      return
    }

    setSaving(true)

    try {
      const request = {
        name: labelName,
        slug: labelSlug,
      }

      if (editingLabel) {
        await updateLabel(editingLabel.id, request)
        toast.success('标签已更新')
      }
      else {
        await createLabel(request)
        toast.success('标签已创建')
      }

      cancelEdit()
      router.refresh()
    }
    catch {
      toast.error(
        editingLabel
          ? '更新标签失败，请检查名称或 Slug 是否重复'
          : '创建标签失败，请检查名称或 Slug 是否重复',
      )
    }
    finally {
      setSaving(false)
    }
  }

  async function handleDelete() {
    if (!pendingDelete)
      return

    setDeleting(true)

    try {
      await deleteLabel(pendingDelete.id)

      if (editingLabel?.id === pendingDelete.id)
        cancelEdit()

      setPendingDelete(null)
      toast.success('标签已删除')
      router.refresh()
    }
    catch {
      toast.error('删除标签失败')
    }
    finally {
      setDeleting(false)
    }
  }

  const deleteDialog = (
    <Dialog.Root
      open={pendingDelete !== null}
      onOpenChange={(open) => {
        if (!open && !deleting)
          setPendingDelete(null)
      }}
    >
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-40 bg-black/45" />
        <Dialog.Content className="fixed top-1/2 left-1/2 z-50 w-[min(92vw,420px)] -translate-x-1/2 -translate-y-1/2 rounded-md bg-white p-6 shadow-xl dark:bg-neutral-950">
          <Dialog.Title className="text-lg font-semibold">
            删除标签
          </Dialog.Title>
          <Dialog.Description className="mt-2 text-sm leading-6 text-neutral-600 dark:text-neutral-400">
            确定删除标签“
            {pendingDelete?.name}
            ”吗？删除后无法恢复。
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

  return (
    <section aria-labelledby="labels-title" className="space-y-6">
      <header>
        <h1 id="labels-title" className="text-2xl font-semibold">
          标签管理
        </h1>
        <p className="mt-1 text-sm text-neutral-500">
          共
          {' '}
          {labels.length}
          {' '}
          个标签
        </p>
      </header>

      <form
        onSubmit={handleSave}
        className="grid gap-4 border-y border-black/10 py-5 dark:border-white/10 md:grid-cols-2"
      >
        <label className="space-y-2 text-sm">
          <span>名称</span>
          <input
            value={name}
            onChange={event => setName(event.target.value)}
            maxLength={64}
            className="min-h-10 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
          />
        </label>

        <label className="space-y-2 text-sm">
          <span>Slug</span>
          <input
            value={slug}
            onChange={event => setSlug(event.target.value)}
            maxLength={64}
            className="min-h-10 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
          />
        </label>

        <div className="flex items-center gap-2 md:col-span-2">
          <button
            type="submit"
            disabled={saving}
            className="inline-flex min-h-10 items-center gap-2 rounded-md bg-black px-4 text-sm text-white disabled:opacity-50 dark:bg-white dark:text-black"
          >
            {editingLabel
              ? <Pencil className="size-4" aria-hidden="true" />
              : <Plus className="size-4" aria-hidden="true" />}
            {saving
              ? '保存中'
              : editingLabel
                ? '保存修改'
                : '创建标签'}
          </button>
          {editingLabel && (
            <button
              type="button"
              onClick={cancelEdit}
              disabled={saving}
              aria-label="取消编辑"
              title="取消编辑"
              className="inline-flex size-10 items-center justify-center rounded-md border border-black/15 disabled:opacity-50 dark:border-white/15"
            >
              <X className="size-4" aria-hidden="true" />
            </button>
          )}
        </div>
      </form>

      <div className="overflow-x-auto border-y border-black/10 dark:border-white/10">
        <table className="w-full min-w-[520px] text-left text-sm">
          <thead className="text-neutral-500">
            <tr>
              <th className="px-3 py-3">名称</th>
              <th className="px-3 py-3">Slug</th>
              <th className="px-3 py-3">操作</th>
            </tr>
          </thead>
          <tbody>
            {labels.map(label => (
              <tr
                key={label.id}
                className="border-t border-black/10 dark:border-white/10"
              >
                <td className="px-3 py-4 font-medium">{label.name}</td>
                <td className="px-3 py-4">{label.slug}</td>
                <td className="px-3 py-4">
                  <button
                    type="button"
                    onClick={() => openEdit(label)}
                    aria-label={`编辑标签 ${label.name}`}
                    title="编辑标签"
                    className="inline-flex size-9 items-center justify-center"
                  >
                    <Pencil className="size-4" aria-hidden="true" />
                  </button>
                  <button
                    type="button"
                    onClick={() => setPendingDelete(label)}
                    disabled={deleting}
                    aria-label={`删除标签 ${label.name}`}
                    title="删除标签"
                    className="inline-flex size-9 items-center justify-center text-red-600 disabled:opacity-50"
                  >
                    <Trash2 className="size-4" aria-hidden="true" />
                  </button>
                </td>
              </tr>
            ))}

            {labels.length === 0 && (
              <tr>
                <td
                  colSpan={3}
                  className="px-3 py-12 text-center text-neutral-500"
                >
                  暂无标签
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {deleteDialog}
    </section>
  )
}
