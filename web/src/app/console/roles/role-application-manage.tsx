'use client'
import type {
  ListRoleApplicationsResponse,
  RoleApplication,
  RoleApplicationStatus,
} from '@/models/role'
import * as Dialog from '@radix-ui/react-dialog'
import { Check, X } from 'lucide-react'
import { useLocale, useTranslations } from 'next-intl'
import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import {
  approveRoleApplication,
  listRoleApplications,
  rejectRoleApplication,
} from '@/api/role'

export function RoleApplicationManage() {
  const locale = useLocale()
  const t = useTranslations('Console.roles')
  const [result, setResult] = useState<ListRoleApplicationsResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [loadError, setLoadError] = useState('')
  const [status, setStatus] = useState<RoleApplicationStatus>('pending')
  const [page, setPage] = useState(1)
  const [approvingId, setApprovingId] = useState<number | null>(null)
  const [pendingReject, setPendingReject] = useState<RoleApplication | null>(null)
  const [rejectReason, setRejectReason] = useState('')
  const [rejectingId, setRejectingId] = useState<number | null>(null)
  const [reloadKey, setReloadKey] = useState(0)
  useEffect(() => {
    // 这是 React useEffect **防止组件卸载之后还去执行 setState 的经典模式**。
    let cancelled = false
    async function loadApplications() {
      setLoading(true)
      setLoadError('')
      try {
        const data = await listRoleApplications({
          page,
          pageSize: 20,
          status,
        })
        if (!cancelled) {
          const lastPage = Math.max(
            1,
            Math.ceil(data.total / data.pageSize),
          )

          if (page > lastPage) {
            setPage(lastPage)
            return
          }

          setResult(data)
        }
      }
      catch {
        if (!cancelled)
          setLoadError(t('applicationsLoadFailed'))
      }
      finally {
        if (!cancelled)
          setLoading(false)
      }
    }
    void loadApplications()
    return () => {
      cancelled = true
    }
  }, [page, reloadKey, status, t])

  async function handleApprove(id: number) {
    setApprovingId(id)

    try {
      await approveRoleApplication(id)
      toast.success(t('approvedMessage'))
      setReloadKey(value => value + 1)
    }
    catch {
      toast.error(t('approveFailed'))
    }
    finally {
      setApprovingId(null)
    }
  }

  async function handleReject() {
    if (!pendingReject)
      return
    const reason = rejectReason.trim()
    if (!reason) {
      toast.error(t('reasonRequired'))
      return
    }
    setRejectingId(pendingReject.id)
    try {
      await rejectRoleApplication(pendingReject.id, reason)
      toast.success(t('rejectedMessage'))
      setPendingReject(null)
      setRejectReason('')
      setReloadKey(value => value + 1)
    }
    catch {
      toast.error(t('rejectFailed'))
    }
    finally {
      setRejectingId(null)
    }
  }

  const totalPages = Math.max(
    1,
    Math.ceil(
      (result?.total ?? 0) / (result?.pageSize ?? 20),
    ),
  )

  return (
    <Dialog.Root
      open={pendingReject !== null}
      onOpenChange={(open) => {
        if (!open && rejectingId === null) {
          setPendingReject(null)
          setRejectReason('')
        }
      }}
    >
      <section aria-labelledby="role-applications-title">
        <h1
          id="role-applications-title"
          className="text-2xl font-semibold"
        >
          {t('applicationsTitle')}
        </h1>
        <p className="mt-2 text-sm text-neutral-500">
          {loading
            ? t('loading')
            : loadError || t('applicationsCount', { count: result?.total ?? 0 })}
        </p>
        <select
          aria-label={t('applicationStatus')}
          value={status}
          onChange={(event) => {
            setStatus(event.target.value as RoleApplicationStatus)
            setPage(1)
          }}
        >
          <option value="pending">{t('pending')}</option>
          <option value="approved">{t('approved')}</option>
          <option value="rejected">{t('rejected')}</option>
        </select>
        <div className="mt-6 overflow-x-auto border-y border-black/10 dark:border-white/10">
          <table className="w-full min-w-[720px] text-left text-sm">
            <thead className="text-neutral-500">
              <tr>
                <th>{t('user')}</th>
                <th>{t('requestedRole')}</th>
                <th>{t('status')}</th>
                <th>{t('applicationTime')}</th>
                <th className="px-3 py-3">{t('actions')}</th>
              </tr>
            </thead>
            <tbody>
              {result?.items.map(application => (
                <tr
                  key={application.id}
                  className="border-t border-black/10 dark:border-white/10"
                >
                  <td className="px-3 py-4">
                    {application.user.nickname || application.user.username}
                  </td>
                  <td className="px-3 py-4">
                    {application.requestedRole.name}
                  </td>
                  <td className="px-3 py-4">
                    {t(application.status)}
                  </td>
                  <td className="px-3 py-4">
                    {new Date(application.createdAt).toLocaleString(locale)}
                  </td>
                  <td className="px-3 py-4">
                    {application.status === 'pending' && (
                      <div className="flex gap-2">
                        <button
                          type="button"
                          aria-label={t('approveAction')}
                          title={t('approveAction')}
                          disabled={approvingId !== null || rejectingId !== null}
                          className="inline-flex size-9 items-center justify-center rounded-md disabled:opacity-50"
                          onClick={() => void handleApprove(application.id)}
                        >
                          <Check className="size-4" aria-hidden="true" />
                        </button>
                        <button
                          type="button"
                          aria-label={t('rejectAction')}
                          title={t('rejectAction')}
                          disabled={approvingId !== null || rejectingId !== null}
                          className="inline-flex size-9 items-center justify-center rounded-md text-red-600 disabled:opacity-50"
                          onClick={() => {
                            setPendingReject(application)
                            setRejectReason('')
                          }}
                        >
                          <X className="size-4" aria-hidden="true" />
                        </button>
                      </div>
                    )}
                  </td>
                </tr>
              ))}
              {!loading && result?.items.length === 0 && (
                <tr>
                  <td
                    colSpan={5}
                    className="px-3 py-12 text-center text-neutral-500"
                  >
                    {t('applicationsEmpty')}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
        <nav
          aria-label={t('applicationsPagination')}
          className="mt-5 flex items-center justify-end gap-4 text-sm"
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
      </section>

      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-40 bg-black/45" />
        <Dialog.Content className="fixed top-1/2 left-1/2 z-50 w-[min(92vw,420px)] -translate-x-1/2 -translate-y-1/2 rounded-md bg-white p-6 shadow-xl dark:bg-neutral-950">
          <Dialog.Title className="text-lg font-semibold">
            {t('rejectTitle')}
          </Dialog.Title>
          <Dialog.Description className="mt-2 text-sm leading-6 text-neutral-600 dark:text-neutral-400">
            {t('rejectDescription', {
              user: pendingReject?.user.nickname || pendingReject?.user.username || '',
              role: pendingReject?.requestedRole.name || '',
            })}
          </Dialog.Description>
          <textarea
            value={rejectReason}
            onChange={event => setRejectReason(event.target.value)}
            placeholder={t('reasonPlaceholder')}
            rows={4}
            className="mt-4 w-full resize-none rounded-md border border-black/15 p-3 text-sm dark:border-white/15"
          />
          <div className="mt-6 flex justify-end gap-3">
            <Dialog.Close asChild>
              <button
                type="button"
                disabled={rejectingId !== null}
                className="min-h-10 rounded-md border border-black/15 px-4 text-sm disabled:opacity-50 dark:border-white/15"
              >
                {t('cancel')}
              </button>
            </Dialog.Close>
            <button
              type="button"
              disabled={rejectingId !== null || !rejectReason.trim()}
              onClick={() => void handleReject()}
              className="min-h-10 rounded-md bg-red-600 px-4 text-sm text-white disabled:opacity-50"
            >
              {rejectingId === null ? t('confirmReject') : t('submitting')}
            </button>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}
