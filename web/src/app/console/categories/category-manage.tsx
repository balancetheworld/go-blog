'use client'

import type { FormEvent } from 'react'
import type { Category } from '@/models/post'
import * as Dialog from '@radix-ui/react-dialog'
import { Pencil, Plus, Trash2, X } from 'lucide-react'
import { useTranslations } from 'next-intl'
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
  const t = useTranslations('Console.taxonomy')
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
      toast.error(t('categoryNameRequired'))
      return
    }

    if (!categorySlug) {
      toast.error(t('categorySlugRequired'))
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
        toast.success(t('categoryUpdated'))
      }
      else {
        await createCategory(request)
        toast.success(t('categoryCreated'))
      }

      cancelEdit()
      router.refresh()
    }
    catch {
      toast.error(
        editingCategory
          ? t('categoryUpdateFailed')
          : t('categoryCreateFailed'),
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
      toast.success(t('categoryDeleted'))
      router.refresh()
    }
    catch {
      toast.error(t('categoryDeleteFailed'))
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
            {t('deleteCategory')}
          </Dialog.Title>
          <Dialog.Description className="mt-2 text-sm leading-6 text-neutral-600 dark:text-neutral-400">
            {t('deleteCategoryDescription', { name: pendingDelete?.name ?? '' })}
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

  return (
    <section aria-labelledby="categories-title" className="space-y-6">
      <header>
        <h1 id="categories-title" className="text-2xl font-semibold">
          {t('categoriesTitle')}
        </h1>
        <p className="mt-1 text-sm text-neutral-500">
          {t('categoriesCount', { count: categories.length })}
        </p>
      </header>

      <form
        onSubmit={handleSave}
        className="grid gap-4 border-y border-black/10 py-5 dark:border-white/10 md:grid-cols-2"
      >
        <label className="space-y-2 text-sm">
          <span>{t('name')}</span>
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
          <span>{t('description')}</span>
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
              ? t('saving')
              : editingCategory
                ? t('saveChanges')
                : t('createCategory')}
          </button>
          {editingCategory && (
            <button
              type="button"
              onClick={cancelEdit}
              disabled={saving}
              aria-label={t('cancelEdit')}
              title={t('cancelEdit')}
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
              <th className="px-3 py-3">{t('name')}</th>
              <th className="px-3 py-3">Slug</th>
              <th className="px-3 py-3">{t('description')}</th>
              <th className="px-3 py-3">{t('actions')}</th>
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
                  {category.description || t('noDescription')}
                </td>
                <td className="px-3 py-4">
                  <button
                    type="button"
                    onClick={() => openEdit(category)}
                    aria-label={t('editCategory', { name: category.name })}
                    title={t('editCategoryTitle')}
                    className="inline-flex size-9 items-center justify-center"
                  >
                    <Pencil className="size-4" aria-hidden="true" />
                  </button>
                  <button
                    type="button"
                    onClick={() => setPendingDelete(category)}
                    disabled={deleting}
                    aria-label={t('deleteCategoryAction', { name: category.name })}
                    title={t('deleteCategory')}
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
                  {t('emptyCategories')}
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
