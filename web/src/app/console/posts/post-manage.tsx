'use client'

import type { FormEvent } from 'react'
import type { PostListResponse } from '@/models/post'
import * as Dialog from '@radix-ui/react-dialog'
import { Pencil, Plus, Search, Trash2 } from 'lucide-react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { deletePost } from '@/api/post'

interface PostManageProps {
  result: PostListResponse
  initialKeyword: string
}

// 创建一个**中文、中等长度样式**的全局日期格式化工具，传入 Date 就能直接输出 `2026年8月7日` 格式的中文日期
const dateFormatter = new Intl.DateTimeFormat('zh-CN', {
  dateStyle: 'medium',
})

export function PostManage({
  result,
  initialKeyword,
}: PostManageProps) {
  const router = useRouter()
  const [keyword, setKeyword] = useState(initialKeyword)
  const [deletingID, setDeletingID] = useState<number | null>(null)
  const [pendingDelete, setPendingDelete] = useState<{
    id: number
    title: string
  } | null>(null)

  function handleSearch(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const value = keyword.trim()
    // `URLSearchParams` 是浏览器 / Node.js 原生内置 API，专门处理 URL 问号后面的查询参数（query string）
    const params = new URLSearchParams()

    if (value) {
      params.set('keyword', value)
    }

    // 添加一个查询参数键值对
    params.set('page', '1')
    router.push(`/console/posts?${params.toString()}`)
  }

  async function handleDelete() {
    if (!pendingDelete)
      return

    setDeletingID(pendingDelete.id)

    try {
      await deletePost(pendingDelete.id)
      toast.success('文章已删除')
      setPendingDelete(null)
      router.refresh()
    }
    catch {
      toast.error('删除文章失败')
    }
    finally {
      setDeletingID(null)
    }
  }

  return (
    <Dialog.Root
      open={pendingDelete !== null}
      onOpenChange={(open) => {
        if (!open && deletingID === null)
          setPendingDelete(null)
      }}
    >
      <div className="space-y-6">
        <header className="flex flex-wrap items-center justify-between gap-4">
          <div>
            <h1 id="posts-title" className="text-2xl font-semibold">文章管理</h1>
            <p className="mt-1 text-sm text-neutral-500">
              共
              {' '}
              {result.total}
              {' '}
              篇文章
            </p>
          </div>

          <Link
            href="/console/posts/new"
            className="inline-flex min-h-10 items-center gap-2 rounded-md bg-black px-4 text-sm text-white dark:bg-white dark:text-black"
          >
            <Plus className="size-4" aria-hidden="true" />
            新建文章
          </Link>
        </header>

        <form onSubmit={handleSearch} className="flex max-w-lg gap-2">
          <input
            aria-label="搜索文章"
            value={keyword}
            onChange={event => setKeyword(event.target.value)}
            placeholder="搜索文章"
            className="min-h-10 min-w-0 flex-1 rounded-md border border-black/15 px-3 dark:border-white/15"
          />
          <button
            type="submit"
            className="inline-flex size-10 items-center justify-center rounded-md border border-black/15 dark:border-white/15"
            aria-label="搜索"
            title="搜索"
          >
            <Search className="size-4" aria-hidden="true" />
          </button>
        </form>

        <div className="overflow-x-auto border-y border-black/10 dark:border-white/10">
          <table className="w-full min-w-[760px] text-left text-sm">
            <thead className="text-neutral-500">
              <tr>
                <th className="px-3 py-3">标题</th>
                <th className="px-3 py-3">状态</th>
                <th className="px-3 py-3">分类</th>
                <th className="px-3 py-3">可见范围</th>
                <th className="px-3 py-3">发布时间</th>
                <th className="px-3 py-3">操作</th>
              </tr>
            </thead>
            <tbody>
              {result.items.map(post => (
                <tr key={post.id} className="border-t border-black/10 dark:border-white/10">
                  <td className="px-3 py-4 font-medium">{post.title}</td>
                  <td className="px-3 py-4">
                    {post.status === 'published' ? '已发布' : '草稿'}
                  </td>
                  <td className="px-3 py-4">{post.category?.name ?? '未分类'}</td>
                  <td className="px-3 py-4">{post.isPrivate ? '私密' : '公开'}</td>
                  <td className="px-3 py-4">
                    {post.publishedAt
                      ? dateFormatter.format(new Date(post.publishedAt))
                      : '未发布'}
                  </td>
                  <td className="flex gap-2 px-3 py-4">
                    <Link
                      href={`/console/posts/edit/${post.id}`}
                      aria-label="编辑文章"
                      title="编辑文章"
                      className="inline-flex size-9 items-center justify-center"
                    >
                      <Pencil className="size-4" aria-hidden="true" />
                    </Link>
                    <button
                      type="button"
                      aria-label="删除文章"
                      title="删除文章"
                      disabled={deletingID !== null}
                      className="inline-flex size-9 items-center justify-center disabled:opacity-50"
                      onClick={() => setPendingDelete({ id: post.id, title: post.title })}
                    >
                      <Trash2 className="size-4" aria-hidden="true" />
                    </button>
                  </td>
                </tr>
              ))}
              {result.items.length === 0 && (
                <tr>
                  <td colSpan={6} className="px-3 py-12 text-center text-neutral-500">
                    暂无文章
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
          <Dialog.Title className="text-lg font-semibold">
            删除文章
          </Dialog.Title>
          <Dialog.Description className="mt-2 text-sm leading-6 text-neutral-600 dark:text-neutral-400">
            确定删除文章“
            {pendingDelete?.title}
            ”吗？删除后无法恢复。
          </Dialog.Description>
          <div className="mt-6 flex justify-end gap-3">
            <Dialog.Close asChild>
              <button
                type="button"
                disabled={deletingID !== null}
                className="min-h-10 rounded-md border border-black/15 px-4 text-sm disabled:opacity-50 dark:border-white/15"
              >
                取消
              </button>
            </Dialog.Close>
            <button
              type="button"
              disabled={deletingID !== null}
              onClick={() => void handleDelete()}
              className="min-h-10 rounded-md bg-red-600 px-4 text-sm text-white disabled:opacity-50"
            >
              {deletingID === null ? '删除' : '删除中'}
            </button>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}
