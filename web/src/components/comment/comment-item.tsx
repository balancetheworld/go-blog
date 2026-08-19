'use client'

import type { Comment } from '@/models/comment'
import * as Dialog from '@radix-ui/react-dialog'
import { ChevronDown, MessageSquareReply, Trash2 } from 'lucide-react'
import { useLocale, useTranslations } from 'next-intl'
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

export function CommentItem({
  comment,
  targetAuthorID,
  onDeleted,
}: CommentItemProps) {
  const locale = useLocale()
  const t = useTranslations('Comments')
  const dateFormatter = new Intl.DateTimeFormat(locale, {
    dateStyle: 'medium',
    timeStyle: 'short',
  })
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
      toast.error(t('repliesFailed'))
    }
    finally {
      setLoadingReplies(false)
    }
  }

  async function handleReplyCreated(_comment: Comment) {
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
      toast.success(t('deleted'))
    }
    catch {
      toast.error(t('deleteFailed'))
    }
    finally {
      setDeleting(false)
    }
  }

  return (
    <article className={`comment-item${comment.depth > 0 ? ' comment-item-reply' : ''}`}>
      <div className="comment-avatar">{authorName.charAt(0).toUpperCase()}</div>
      <div className="comment-body">
        <header className="comment-meta">
          <div>
            <strong>{authorName}</strong>
            <div className="comment-reply-context">
              <time dateTime={comment.createdAt}>
                {dateFormatter.format(new Date(comment.createdAt))}
              </time>
              {comment.replyToUser && (
                <span>{t('replyTo', { name: comment.replyToUser.nickname || comment.replyToUser.username })}</span>
              )}
              {comment.moderationStatus === 'pending' && (
                <span>{t('moderationPending')}</span>
              )}
              {comment.moderationStatus === 'manual_review' && (
                <span>{t('moderationManualReview')}</span>
              )}
              {comment.moderationStatus === 'rejected' && (
                <span>{t('moderationRejected')}</span>
              )}
            </div>
          </div>

          <div className="comment-actions">
            {currentUser && comment.depth < 2 && (
              <button
                type="button"
                aria-label={t('replyTo', { name: authorName })}
                title={t('reply')}
                onClick={() => setReplying(current => !current)}
                className="comment-icon-action"
              >
                <MessageSquareReply className="size-4" aria-hidden="true" />
              </button>
            )}
            {canDelete && (
              <button
                type="button"
                aria-label={t('delete')}
                title={t('delete')}
                onClick={() => setDeleteOpen(true)}
                className="comment-icon-action"
              >
                <Trash2 className="size-4" aria-hidden="true" />
              </button>
            )}
          </div>
        </header>

        <p className="comment-content">
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
            className="comment-replies-action"
          >
            <ChevronDown className="size-4" aria-hidden="true" />
            {loadingReplies ? t('loadingReplies') : t('viewReplies', { count: comment.replyCount })}
          </button>
        )}
        {replies && replies.length > 0 && (
          <div className="comment-replies">
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
      </div>

      <Dialog.Root open={deleteOpen} onOpenChange={setDeleteOpen}>
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
                {deleting ? t('deleting') : t('deleteConfirm')}
              </button>
            </div>
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>
    </article>
  )
}
