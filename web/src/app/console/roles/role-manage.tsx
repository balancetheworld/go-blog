'use client'

import type { FormEvent } from 'react'
import type { ListRolesResponse, Role } from '@/models/role'
import * as Dialog from '@radix-ui/react-dialog'
import { Pencil, Plus, Search, Trash2 } from 'lucide-react'
import { useTranslations } from 'next-intl'
import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import {
  createRole,
  deleteRole,
  listRoles,
  updateRole,
} from '@/api/role'

interface RoleFormState {
  code: string
  name: string
  description: string
  isRequestable: boolean
  enabled: boolean
}

const emptyForm: RoleFormState = {
  code: '',
  name: '',
  description: '',
  isRequestable: false,
  enabled: true,
}

export function RoleManage() {
  const t = useTranslations('Console.roles')
  const [result, setResult] = useState<ListRolesResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [loadError, setLoadError] = useState('')
  const [keyword, setKeyword] = useState('')
  const [appliedKeyword, setAppliedKeyword] = useState('')
  const [page, setPage] = useState(1)
  const [reloadKey, setReloadKey] = useState(0)
  const [editorOpen, setEditorOpen] = useState(false)
  const [editingRole, setEditingRole] = useState<Role | null>(null)
  const [form, setForm] = useState<RoleFormState>(emptyForm)
  const [saving, setSaving] = useState(false)
  const [pendingDelete, setPendingDelete] = useState<Role | null>(null)
  const [deleting, setDeleting] = useState(false)

  useEffect(() => {
    let cancelled = false

    async function loadRoles() {
      setLoading(true)
      setLoadError('')

      try {
        const data = await listRoles({
          page,
          pageSize: 20,
          keyword: appliedKeyword,
        })

        if (!cancelled) {
          const lastPage = Math.max(1, Math.ceil(data.total / data.pageSize))

          if (page > lastPage) {
            setPage(lastPage)
            return
          }

          setResult(data)
        }
      }
      catch {
        if (!cancelled)
          setLoadError(t('loadFailed'))
      }
      finally {
        if (!cancelled)
          setLoading(false)
      }
    }

    void loadRoles()

    return () => {
      cancelled = true
    }
  }, [appliedKeyword, page, reloadKey, t])

  function openCreate() {
    setEditingRole(null)
    setForm(emptyForm)
    setEditorOpen(true)
  }

  function openEdit(role: Role) {
    setEditingRole(role)
    setForm({
      code: role.code,
      name: role.name,
      description: role.description,
      isRequestable: role.isRequestable,
      enabled: role.enabled,
    })
    setEditorOpen(true)
  }

  function handleSearch(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setPage(1)
    setAppliedKeyword(keyword.trim())
  }

  async function handleSave(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const code = form.code.trim().toLowerCase()
    const name = form.name.trim()
    if (!editingRole && !/^[a-z][a-z0-9_-]*$/.test(code)) {
      toast.error(t('codeInvalid'))
      return
    }
    if (!name) {
      toast.error(t('nameRequired'))
      return
    }

    setSaving(true)

    try {
      if (editingRole) {
        await updateRole(editingRole.id, {
          name,
          description: form.description.trim(),
          isRequestable: form.isRequestable,
          enabled: form.enabled,
        })
        toast.success(t('updated'))
      }
      else {
        await createRole({
          code,
          name,
          description: form.description.trim(),
          isRequestable: form.isRequestable,
          enabled: form.enabled,
        })
        toast.success(t('created'))
      }

      setEditorOpen(false)
      setReloadKey(value => value + 1)
    }
    catch {
      toast.error(editingRole ? t('updateFailed') : t('createFailed'))
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
      await deleteRole(pendingDelete.id)
      toast.success(t('deleted'))
      setPendingDelete(null)
      setReloadKey(value => value + 1)
    }
    catch {
      toast.error(t('deleteFailed'))
    }
    finally {
      setDeleting(false)
    }
  }

  const totalPages = Math.max(
    1,
    Math.ceil((result?.total ?? 0) / (result?.pageSize ?? 20)),
  )

  return (
    <section aria-labelledby="roles-title" className="space-y-6">
      <header className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 id="roles-title" className="text-2xl font-semibold">
            {t('title')}
          </h1>
          <p className="mt-1 text-sm text-neutral-500">
            {loading
              ? t('loading')
              : loadError || t('count', { count: result?.total ?? 0 })}
          </p>
        </div>
        <button
          type="button"
          onClick={openCreate}
          className="inline-flex min-h-10 items-center gap-2 rounded-md bg-black px-4 text-sm text-white dark:bg-white dark:text-black"
        >
          <Plus className="size-4" aria-hidden="true" />
          {t('new')}
        </button>
      </header>

      <form onSubmit={handleSearch} className="flex max-w-lg gap-2">
        <input
          aria-label={t('search')}
          value={keyword}
          onChange={event => setKeyword(event.target.value)}
          placeholder={t('searchPlaceholder')}
          className="min-h-10 min-w-0 flex-1 rounded-md border border-black/15 px-3 dark:border-white/15"
        />
        <button
          type="submit"
          aria-label={t('searchAction')}
          title={t('searchAction')}
          className="inline-flex size-10 items-center justify-center rounded-md border border-black/15 dark:border-white/15"
        >
          <Search className="size-4" aria-hidden="true" />
        </button>
      </form>

      <div className="overflow-x-auto border-y border-black/10 dark:border-white/10">
        <table className="w-full min-w-[760px] text-left text-sm">
          <thead className="text-neutral-500">
            <tr>
              <th className="px-3 py-3">{t('name')}</th>
              <th className="px-3 py-3">{t('code')}</th>
              <th className="px-3 py-3">{t('type')}</th>
              <th className="px-3 py-3">{t('requestable')}</th>
              <th className="px-3 py-3">{t('status')}</th>
              <th className="px-3 py-3">{t('actions')}</th>
            </tr>
          </thead>
          <tbody>
            {result?.items.map(role => (
              <tr key={role.id} className="border-t border-black/10 dark:border-white/10">
                <td className="px-3 py-4">
                  <div className="font-medium">{role.name}</div>
                  {role.description && (
                    <div className="mt-1 max-w-sm text-neutral-500">
                      {role.description}
                    </div>
                  )}
                </td>
                <td className="px-3 py-4">{role.code}</td>
                <td className="px-3 py-4">
                  {role.isSystem ? t('system') : t('custom')}
                  {role.isDefault && ` / ${t('default')}`}
                </td>
                <td className="px-3 py-4">{role.isRequestable ? t('yes') : t('no')}</td>
                <td className="px-3 py-4">{role.enabled ? t('enabled') : t('disabled')}</td>
                <td className="px-3 py-4">
                  {!role.isSystem && (
                    <div className="flex gap-2">
                      <button
                        type="button"
                        aria-label={t('editAction', { name: role.name })}
                        title={t('editTitle')}
                        className="inline-flex size-9 items-center justify-center"
                        onClick={() => openEdit(role)}
                      >
                        <Pencil className="size-4" aria-hidden="true" />
                      </button>
                      <button
                        type="button"
                        aria-label={t('deleteAction', { name: role.name })}
                        title={t('deleteTitle')}
                        className="inline-flex size-9 items-center justify-center text-red-600"
                        onClick={() => setPendingDelete(role)}
                      >
                        <Trash2 className="size-4" aria-hidden="true" />
                      </button>
                    </div>
                  )}
                </td>
              </tr>
            ))}
            {!loading && result?.items.length === 0 && (
              <tr>
                <td colSpan={6} className="px-3 py-12 text-center text-neutral-500">
                  {t('empty')}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <nav
        aria-label={t('pagination')}
        className="flex items-center justify-end gap-4 text-sm"
      >
        <button
          type="button"
          disabled={loading || page <= 1}
          onClick={() => setPage(value => value - 1)}
          className="disabled:text-neutral-400"
        >
          {t('previous')}
        </button>
        <span>
          {page}
          {' '}
          /
          {' '}
          {totalPages}
        </span>
        <button
          type="button"
          disabled={loading || page >= totalPages}
          onClick={() => setPage(value => value + 1)}
          className="disabled:text-neutral-400"
        >
          {t('next')}
        </button>
      </nav>

      <Dialog.Root
        open={editorOpen}
        onOpenChange={(open) => {
          if (!saving)
            setEditorOpen(open)
        }}
      >
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 z-40 bg-black/45" />
          <Dialog.Content className="fixed top-1/2 left-1/2 z-50 max-h-[90vh] w-[min(92vw,520px)] -translate-x-1/2 -translate-y-1/2 overflow-y-auto rounded-md bg-white p-6 shadow-xl dark:bg-neutral-950">
            <Dialog.Title className="text-lg font-semibold">
              {editingRole ? t('editTitle') : t('new')}
            </Dialog.Title>
            <Dialog.Description className="mt-2 text-sm text-neutral-500">
              {editingRole ? t('editDescription') : t('createDescription')}
            </Dialog.Description>
            <form onSubmit={handleSave} className="mt-6 space-y-4">
              <label className="block text-sm">
                <span>{t('roleCode')}</span>
                <input
                  value={form.code}
                  disabled={editingRole !== null}
                  required
                  maxLength={64}
                  onChange={event => setForm(value => ({ ...value, code: event.target.value }))}
                  className="mt-2 min-h-10 w-full rounded-md border border-black/15 px-3 disabled:bg-black/5 dark:border-white/15 dark:disabled:bg-white/5"
                />
              </label>
              <label className="block text-sm">
                <span>{t('roleName')}</span>
                <input
                  value={form.name}
                  required
                  maxLength={64}
                  onChange={event => setForm(value => ({ ...value, name: event.target.value }))}
                  className="mt-2 min-h-10 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
                />
              </label>
              <label className="block text-sm">
                <span>{t('roleDescription')}</span>
                <textarea
                  value={form.description}
                  maxLength={512}
                  rows={4}
                  onChange={event => setForm(value => ({ ...value, description: event.target.value }))}
                  className="mt-2 w-full resize-none rounded-md border border-black/15 p-3 dark:border-white/15"
                />
              </label>
              <label className="flex items-center gap-3 text-sm">
                <input
                  type="checkbox"
                  checked={form.isRequestable}
                  onChange={event => setForm(value => ({ ...value, isRequestable: event.target.checked }))}
                />
                {t('allowRequest')}
              </label>
              <label className="flex items-center gap-3 text-sm">
                <input
                  type="checkbox"
                  checked={form.enabled}
                  onChange={event => setForm(value => ({ ...value, enabled: event.target.checked }))}
                />
                {t('enableRole')}
              </label>
              <div className="flex justify-end gap-3 pt-2">
                <Dialog.Close asChild>
                  <button
                    type="button"
                    disabled={saving}
                    className="min-h-10 rounded-md border border-black/15 px-4 text-sm disabled:opacity-50 dark:border-white/15"
                  >
                    {t('cancel')}
                  </button>
                </Dialog.Close>
                <button
                  type="submit"
                  disabled={saving}
                  className="min-h-10 rounded-md bg-black px-4 text-sm text-white disabled:opacity-50 dark:bg-white dark:text-black"
                >
                  {saving ? t('saving') : t('save')}
                </button>
              </div>
            </form>
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>

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
              {t('deleteTitle')}
            </Dialog.Title>
            <Dialog.Description className="mt-2 text-sm leading-6 text-neutral-600 dark:text-neutral-400">
              {t('deleteDescription', { name: pendingDelete?.name ?? '' })}
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
    </section>
  )
}
