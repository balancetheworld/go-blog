'use client'

import type { FormEvent } from 'react'
import type { Category } from '@/models/post'
import * as Dialog from '@radix-ui/react-dialog'
import { Pencil, Plus, Trash2, X } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import {
  createCategory,
  deleteCategory,
  updateCategory,
} from '@/api/post'

interface CategoryManageProps {
  categories: Category[]
}

export function CategoryManage({
  categories,
}: CategoryManageProps) {
  const router = useRouter()
  const [name, setName] = useState('')
  const [slug, setSlug] = useState('')
  const [description, setDescription] = useState('')
  const [saving, setSaving] = useState(false)
  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [pendingDelete, setPendingDelete] = useState<Category | null>(null)
  const [deleting, setDeleting] = useState(false)

  function openEdit(category: Category) {
    setEditingCategory(category)
    setName(category.name)
    setSlug(category.slug)
    setDescription(category.description)
  }

  function cancelEdit() {
    setEditingCategory(null)
    setName('')
    setSlug('')
    setDescription('')
  }

  async function handleSave(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const categoryName = name.trim()
    const categorySlug = slug.trim()
    if (!categoryName) {
      toast.error('请输入分类名称')
      return
    }

    if (!categorySlug) {
      toast.error('请输入分类 Slug')
      return
    }

    setSaving(true)
    try {
      const request = {
        name: categoryName,
        slug: categorySlug,
        description: description.trim(),
      }

      if (editingCategory) {
        await updateCategory(editingCategory.id, request)
        toast.success('分类已更新')
      }
      else {
        await createCategory(request)
        toast.success('分类已创建')
      }

      cancelEdit()
      router.refresh()
    }
    catch {
      toast.error(
        editingCategory
          ? '更新分类失败，请检查名称或 Slug 是否重复'
          : '创建分类失败，请检查名称或 Slug 是否重复',
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
      await deleteCategory(pendingDelete.id)

      if (editingCategory?.id === pendingDelete.id)
        cancelEdit()

      setPendingDelete(null)
      toast.success('分类已删除')
      router.refresh()
    }
    catch {
      toast.error('删除分类失败')
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
            删除分类
          </Dialog.Title>
          <Dialog.Description className="mt-2 text-sm leading-6 text-neutral-600 dark:text-neutral-400">
            确定删除分类“
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
    <section aria-labelledby="categories-title" className="space-y-6">
      <header>
        <h1 id="categories-title" className="text-2xl font-semibold">
          分类管理
        </h1>
        <p className="mt-1 text-sm text-neutral-500">
          共
          {' '}
          {categories.length}
          {' '}
          个分类
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

        <label className="space-y-2 text-sm md:col-span-2">
          <span>描述</span>
          <textarea
            value={description}
            onChange={event => setDescription(event.target.value)}
            maxLength={512}
            rows={3}
            className="w-full rounded-md border border-black/15 px-3 py-2 dark:border-white/15"
          />
        </label>

        <div className="flex items-center gap-2 md:col-span-2">
          <button
            type="submit"
            disabled={saving}
            className="inline-flex min-h-10 items-center gap-2 rounded-md bg-black px-4 text-sm text-white disabled:opacity-50 dark:bg-white dark:text-black"
          >
            {editingCategory
              ? <Pencil className="size-4" aria-hidden="true" />
              : <Plus className="size-4" aria-hidden="true" />}
            {saving
              ? '保存中'
              : editingCategory
                ? '保存修改'
                : '创建分类'}
          </button>
          {editingCategory && (
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
        <table className="w-full min-w-[640px] text-left text-sm">
          <thead className="text-neutral-500">
            <tr>
              <th className="px-3 py-3">名称</th>
              <th className="px-3 py-3">Slug</th>
              <th className="px-3 py-3">描述</th>
              <th className="px-3 py-3">操作</th>
            </tr>
          </thead>
          <tbody>
            {categories.map(category => (
              <tr
                key={category.id}
                className="border-t border-black/10 dark:border-white/10"
              >
                <td className="px-3 py-4 font-medium">
                  {category.name}
                </td>
                <td className="px-3 py-4">{category.slug}</td>
                <td className="px-3 py-4 text-neutral-500">
                  {category.description || '暂无描述'}
                </td>
                <td className="px-3 py-4">
                  <button
                    type="button"
                    onClick={() => openEdit(category)}
                    aria-label={`编辑分类 ${category.name}`}
                    title="编辑分类"
                    className="inline-flex size-9 items-center justify-center"
                  >
                    <Pencil className="size-4" aria-hidden="true" />
                  </button>
                  <button
                    type="button"
                    onClick={() => setPendingDelete(category)}
                    disabled={deleting}
                    aria-label={`删除分类 ${category.name}`}
                    title="删除分类"
                    className="inline-flex size-9 items-center justify-center text-red-600 disabled:opacity-50"
                  >
                    <Trash2 className="size-4" aria-hidden="true" />
                  </button>
                </td>
              </tr>
            ))}

            {categories.length === 0 && (
              <tr>
                <td
                  colSpan={4}
                  className="px-3 py-12 text-center text-neutral-500"
                >
                  暂无分类
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
