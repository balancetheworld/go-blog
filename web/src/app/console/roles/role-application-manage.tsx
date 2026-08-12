'use client'
import type {
  ListRoleApplicationsResponse,
  RoleApplication,
  RoleApplicationStatus,
} from '@/models/role'
import * as Dialog from '@radix-ui/react-dialog'
import { Check, X } from 'lucide-react'
import { useEffect, useState } from 'react'
import { toast } from 'sonner'
import {
  approveRoleApplication,
  listRoleApplications,
  rejectRoleApplication,
} from '@/api/role'

export function RoleApplicationManage() {
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
          setLoadError('身份申请加载失败')
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
  }, [page, status, reloadKey])

  async function handleApprove(id: number) {
    setApprovingId(id)

    try {
      await approveRoleApplication(id)
      toast.success('身份申请已通过')
      setReloadKey(value => value + 1)
    }
    catch {
      toast.error('审核身份申请失败')
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
      toast.error('请填写拒绝原因')
      return
    }
    setRejectingId(pendingReject.id)
    try {
      await rejectRoleApplication(pendingReject.id, reason)
      toast.success('身份申请已拒绝')
      setPendingReject(null)
      setRejectReason('')
      setReloadKey(value => value + 1)
    }
    catch {
      toast.error('拒绝身份申请失败')
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
          身份审核
        </h1>
        <p className="mt-2 text-sm text-neutral-500">
          {loading
            ? '加载中...'
            : loadError || `共 ${result?.total ?? 0} 条申请`}
        </p>
        <select
          aria-label="申请状态"
          value={status}
          onChange={(event) => {
            setStatus(event.target.value as RoleApplicationStatus)
            setPage(1)
          }}
        >
          <option value="pending">待审核</option>
          <option value="approved">已通过</option>
          <option value="rejected">已拒绝</option>
        </select>
        <div className="mt-6 overflow-x-auto border-y border-black/10 dark:border-white/10">
          <table className="w-full min-w-[720px] text-left text-sm">
            <thead className="text-neutral-500">
              <tr>
                <th>用户</th>
                <th>申请身份</th>
                <th>状态</th>
                <th>申请时间</th>
                <th className="px-3 py-3">操作</th>
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
                    {application.status}
                  </td>
                  <td className="px-3 py-4">
                    {new Date(application.createdAt).toLocaleString('zh-CN')}
                  </td>
                  <td className="px-3 py-4">
                    {application.status === 'pending' && (
                      <div className="flex gap-2">
                        <button
                          type="button"
                          aria-label="通过申请"
                          title="通过申请"
                          disabled={approvingId !== null || rejectingId !== null}
                          className="inline-flex size-9 items-center justify-center rounded-md disabled:opacity-50"
                          onClick={() => void handleApprove(application.id)}
                        >
                          <Check className="size-4" aria-hidden="true" />
                        </button>
                        <button
                          type="button"
                          aria-label="拒绝申请"
                          title="拒绝申请"
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
                    暂无身份申请
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
        <nav
          aria-label="身份申请分页"
          className="mt-5 flex items-center justify-end gap-4 text-sm"
        >
          <button
            type="button"
            disabled={loading || page <= 1}
            onClick={() => setPage(value => value - 1)}
            className="disabled:text-neutral-400"
          >
            上一页
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
            下一页
          </button>
        </nav>
      </section>

      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-40 bg-black/45" />
        <Dialog.Content className="fixed top-1/2 left-1/2 z-50 w-[min(92vw,420px)] -translate-x-1/2 -translate-y-1/2 rounded-md bg-white p-6 shadow-xl dark:bg-neutral-950">
          <Dialog.Title className="text-lg font-semibold">
            拒绝身份申请
          </Dialog.Title>
          <Dialog.Description className="mt-2 text-sm leading-6 text-neutral-600 dark:text-neutral-400">
            拒绝
            {pendingReject?.user.nickname || pendingReject?.user.username}
            申请的
            {pendingReject?.requestedRole.name}
            身份
          </Dialog.Description>
          <textarea
            value={rejectReason}
            onChange={event => setRejectReason(event.target.value)}
            placeholder="请输入拒绝原因"
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
                取消
              </button>
            </Dialog.Close>
            <button
              type="button"
              disabled={rejectingId !== null || !rejectReason.trim()}
              onClick={() => void handleReject()}
              className="min-h-10 rounded-md bg-red-600 px-4 text-sm text-white disabled:opacity-50"
            >
              {rejectingId === null ? '确认拒绝' : '提交中'}
            </button>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}
