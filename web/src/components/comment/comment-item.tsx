'use client'

import type { Comment } from '@/models/comment'
import * as Dialog from '@radix-ui/react-dialog'
import { ChevronDown, MessageSquareReply, Trash2 } from 'lucide-react'
import { useState } from 'react'
import { toast } from 'sonner'
import { deleteComment, listCommentReplies } from '@/api/comment'
import { useAuth } from '@/contexts/auth-context'
import { CommentInput } from './comment-input'

interface CommentItemProps {
  comment: Comment
  targetAuthorID: number
  onDeleted: () => void | Promise<void>
}

const dateFormatter = new Intl.DateTimeFormat('zh-CN', {
  dateStyle: 'medium',
  timeStyle: 'short',
})

export function CommentItem({
  comment,
  targetAuthorID,
  onDeleted,
}: CommentItemProps) {
  const { currentUser, currentRole } = useAuth()
  const [replies, setReplies] = useState<Comment[] | null>(null)
  const [loadingReplies, setLoadingReplies] = useState(false)
  const [replying, setReplying] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const authorName = comment.author.nickname || comment.author.username

  const canDelete = currentUser !== null && (
    currentRole === 'admin'
    || comment.author.id === currentUser.id
    || (currentRole === 'editor' && currentUser.id === targetAuthorID)
  )

  async function loadReplies() {
    if (loadingReplies)
      return

    setLoadingReplies(true)
    try {
      setReplies(await listCommentReplies(comment.id))
    }
    catch {
      toast.error('回复加载失败')
    }
    finally {
      setLoadingReplies(false)
    }
  }

  async function handleReplyCreated() {
    setReplying(false)
    await loadReplies()
  }

  async function handleDelete() {
    if (deleting)
      return

    setDeleting(true)
    try {
      await deleteComment(comment.id)
      setDeleteOpen(false)
      await onDeleted()
      toast.success('评论已删除')
    }
    catch {
      toast.error('评论删除失败')
    }
    finally {
      setDeleting(false)
    }
  }

  return (
    <article className={comment.depth > 0 ? 'border-l border-black/10 pl-4 dark:border-white/10' : ''}>
      <div className="py-5">
        <header className="flex items-start justify-between gap-4">
          <div className="min-w-0">
            <p className="truncate font-medium">{authorName}</p>
            <div className="mt-1 flex flex-wrap items-center gap-x-2 text-xs text-neutral-500">
              <time dateTime={comment.createdAt}>
                {dateFormatter.format(new Date(comment.createdAt))}
              </time>
              {comment.replyToUser && (
                <span>
                  回复
                  {' '}
                  {comment.replyToUser.nickname || comment.replyToUser.username}
                </span>
              )}
            </div>
          </div>

          <div className="flex shrink-0 items-center gap-1">
            {currentUser && comment.depth < 2 && (
              <button
                type="button"
                aria-label={`回复 ${authorName}`}
                title="回复"
                onClick={() => setReplying(current => !current)}
                className="inline-flex size-9 items-center justify-center"
              >
                <MessageSquareReply className="size-4" aria-hidden="true" />
              </button>
            )}
            {canDelete && (
              <button
                type="button"
                aria-label="删除评论"
                title="删除评论"
                onClick={() => setDeleteOpen(true)}
                className="inline-flex size-9 items-center justify-center"
              >
                <Trash2 className="size-4" aria-hidden="true" />
              </button>
            )}
          </div>
        </header>

        <p className="mt-3 whitespace-pre-wrap break-words leading-7">
          {comment.content}
        </p>

        {replying && (
          <CommentInput
            targetType="comment"
            targetID={comment.id}
            replyToName={authorName}
            autoFocus
            onCreated={handleReplyCreated}
            onCancel={() => setReplying(false)}
          />
        )}

        {comment.replyCount > 0 && replies === null && (
          <button
            type="button"
            disabled={loadingReplies}
            onClick={() => void loadReplies()}
            className="mt-3 inline-flex min-h-9 items-center gap-1 text-sm font-medium disabled:opacity-50"
          >
            <ChevronDown className="size-4" aria-hidden="true" />
            {loadingReplies ? '加载中' : `查看 ${comment.replyCount} 条回复`}
          </button>
        )}
      </div>

      {replies && replies.length > 0 && (
        <div className="space-y-1">
          {replies.map(reply => (
            <CommentItem
              key={reply.id}
              comment={reply}
              targetAuthorID={targetAuthorID}
              onDeleted={loadReplies}
            />
          ))}
        </div>
      )}

      <Dialog.Root open={deleteOpen} onOpenChange={setDeleteOpen}>
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 z-40 bg-black/45" />
          <Dialog.Content className="fixed top-1/2 left-1/2 z-50 w-[min(92vw,420px)] -translate-x-1/2 -translate-y-1/2 rounded-md bg-white p-6 shadow-xl dark:bg-neutral-950">
            <Dialog.Title className="text-lg font-semibold">删除评论</Dialog.Title>
            <Dialog.Description className="mt-2 text-sm text-neutral-500">
              删除后无法恢复，相关回复也会被删除。
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
    </article>
  )
}
