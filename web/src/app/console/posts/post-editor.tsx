'use client'

import type { FormEvent } from 'react'
import type { PostPreview } from './post-preview-dialog'
import type {
  Category,
  CreatePostReq,
  Label,
  Post,
  PostVisibility,
} from '@/models/post'
import type { RoleOption } from '@/models/role'
import { Eye, Save, Send } from 'lucide-react'
import { useLocale, useTranslations } from 'next-intl'
import { useRouter } from 'next/navigation'
import { useEffect, useRef, useState } from 'react'
import { toast } from 'sonner'
import { createPost, updatePost } from '@/api/post'
import { SimpleEditor } from '@/components/tiptap-editor/simple-editor'
import { isHTMLContentEmpty } from '@/lib/tiptap-advanced-utils'
import { PostPreviewDialog } from './post-preview-dialog'

interface PostEditorProps {
  post?: Post
  currentUserID: number
  categories: Category[]
  labels: Label[]
  roleOptions: RoleOption[]
}

interface PostEditorDraft {
  title: string
  description: string
  cover: string
  type: string
  categoryId: string
  labelIds: number[]
  visibility: PostVisibility
  visibleRoleIds: number[]
  top: boolean
  content: string
}

interface PostEditorSnapshot {
  version: 1
  savedAt: string
  draft: PostEditorDraft
}

function createInitialDraft(post?: Post): PostEditorDraft {
  return {
    title: post?.title ?? '',
    description: post?.description ?? '',
    cover: post?.cover ?? '',
    type: post?.type ?? 'article',
    categoryId: post?.categoryId ? String(post.categoryId) : '',
    labelIds: post?.labels.map(label => label.id) ?? [],
    visibility: post?.visibility ?? (post?.isPrivate ? 'private' : 'public'),
    visibleRoleIds: post?.visibleRoles.map(role => role.id) ?? [],
    top: post?.top ?? false,
    content: post?.draftContent ?? post?.content ?? '',
  }
}

function readSnapshot(value: string): PostEditorSnapshot | null {
  try {
    const snapshot = JSON.parse(value) as Partial<PostEditorSnapshot>
    if (
      snapshot.version !== 1
      || typeof snapshot.savedAt !== 'string'
      || !snapshot.draft
      || typeof snapshot.draft.title !== 'string'
      || typeof snapshot.draft.content !== 'string'
      || !Array.isArray(snapshot.draft.labelIds)
      || !Array.isArray(snapshot.draft.visibleRoleIds)
    ) {
      return null
    }

    return snapshot as PostEditorSnapshot
  }
  catch {
    return null
  }
}

function formatAutoSaveTime(value: string, locale: string): string {
  return new Intl.DateTimeFormat(locale, {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).format(new Date(value))
}

export function PostEditor({
  post,
  currentUserID,
  categories,
  labels,
  roleOptions,
}: PostEditorProps) {
  const locale = useLocale()
  const t = useTranslations('Console.posts')
  const router = useRouter()
  const [saving, setSaving] = useState(false)
  const [preview, setPreview] = useState<PostPreview | null>(null)
  const initialDraftRef = useRef(createInitialDraft(post))
  const [draft, setDraft] = useState<PostEditorDraft>(
    initialDraftRef.current,
  )
  const [autoSaveMessage, setAutoSaveMessage] = useState('')
  const [autoSaveReady, setAutoSaveReady] = useState(false)
  const draftRef = useRef(draft)
  const roleOptionsRef = useRef(roleOptions)
  const savedDraftRef = useRef(JSON.stringify(initialDraftRef.current))
  const isCreating = post === undefined
  const storageKey = `my-blog:post-editor:${currentUserID}:${post?.id ?? 'new'}`

  draftRef.current = draft
  roleOptionsRef.current = roleOptions

  useEffect(() => {
    const storedValue = localStorage.getItem(storageKey)
    if (!storedValue) {
      setAutoSaveReady(true)
      return
    }

    const snapshot = readSnapshot(storedValue)
    if (!snapshot) {
      localStorage.removeItem(storageKey)
      setAutoSaveReady(true)
      return
    }

    const serverUpdatedAt = post?.updatedAt
      ? new Date(post.updatedAt).getTime()
      : 0
    const localSavedAt = new Date(snapshot.savedAt).getTime()
    if (
      !Number.isFinite(localSavedAt)
      || localSavedAt <= serverUpdatedAt
      || JSON.stringify(snapshot.draft) === savedDraftRef.current
    ) {
      localStorage.removeItem(storageKey)
      setAutoSaveReady(true)
      return
    }

    const allowedRoleIDs = new Set(roleOptionsRef.current.map(role => role.id))
    setDraft({
      ...snapshot.draft,
      visibleRoleIds: snapshot.draft.visibleRoleIds.filter(roleID => allowedRoleIDs.has(roleID)),
    })
    setAutoSaveMessage(t('restored', { time: formatAutoSaveTime(snapshot.savedAt, locale) }))
    setAutoSaveReady(true)
    toast.info(t('restoredToast'))
  }, [locale, post?.updatedAt, storageKey, t])

  useEffect(() => {
    if (!autoSaveReady)
      return

    const serializedDraft = JSON.stringify(draft)
    if (serializedDraft === savedDraftRef.current) {
      localStorage.removeItem(storageKey)
      return
    }

    setAutoSaveMessage(t('autoSaving'))
    const timer = window.setTimeout(() => {
      try {
        const savedAt = new Date().toISOString()
        const snapshot: PostEditorSnapshot = {
          version: 1,
          savedAt,
          draft,
        }
        localStorage.setItem(storageKey, JSON.stringify(snapshot))
        setAutoSaveMessage(t('autoSaved', { time: formatAutoSaveTime(savedAt, locale) }))
      }
      catch {
        setAutoSaveMessage(t('autoSaveFailed'))
      }
    }, 1000)

    return () => window.clearTimeout(timer)
  }, [autoSaveReady, draft, locale, storageKey, t])

  useEffect(() => {
    function saveBeforeLeave() {
      const currentDraft = draftRef.current
      if (JSON.stringify(currentDraft) === savedDraftRef.current)
        return

      const snapshot: PostEditorSnapshot = {
        version: 1,
        savedAt: new Date().toISOString(),
        draft: currentDraft,
      }

      try {
        localStorage.setItem(storageKey, JSON.stringify(snapshot))
      }
      catch {
      }
    }

    function handleVisibilityChange() {
      if (document.visibilityState === 'hidden')
        saveBeforeLeave()
    }

    window.addEventListener('beforeunload', saveBeforeLeave)
    document.addEventListener('visibilitychange', handleVisibilityChange)
    return () => {
      saveBeforeLeave()
      window.removeEventListener('beforeunload', saveBeforeLeave)
      document.removeEventListener('visibilitychange', handleVisibilityChange)
    }
  }, [storageKey])

  function toggleLabel(labelID: number) {
    setDraft(current => ({
      ...current,
      labelIds: current.labelIds.includes(labelID)
        ? current.labelIds.filter(id => id !== labelID)
        : [...current.labelIds, labelID],
    }))
  }

  function toggleVisibleRole(roleID: number) {
    setDraft(current => ({
      ...current,
      visibleRoleIds: current.visibleRoleIds.includes(roleID)
        ? current.visibleRoleIds.filter(id => id !== roleID)
        : [...current.visibleRoleIds, roleID],
    }))
  }

  function handlePreview() {
    const categoryID = Number(draft.categoryId)

    setPreview({
      title: draft.title.trim() || t('untitled'),
      description: draft.description.trim(),
      cover: draft.cover.trim(),
      category: categories.find(category => category.id === categoryID),
      labels: labels.filter(label => draft.labelIds.includes(label.id)),
      visibility: draft.visibility,
      visibleRoles: roleOptions.filter(role => draft.visibleRoleIds.includes(role.id)),
      content: draft.content,
    })
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (saving)
      return

    const submitter = (event.nativeEvent as SubmitEvent)
      .submitter as HTMLButtonElement | null
    const publish = submitter?.value === 'publish'
    const title = draft.title.trim()

    if (!title) {
      toast.error(t('titleRequired'))
      return
    }

    if (publish && isHTMLContentEmpty(draft.content)) {
      toast.error(t('contentRequired'))
      return
    }

    if (draft.visibility === 'roles' && draft.visibleRoleIds.length === 0) {
      toast.error(t('roleRequired'))
      return
    }

    const request: CreatePostReq = {
      title,
      description: draft.description.trim(),
      cover: draft.cover.trim(),
      type: draft.type,
      draftContent: draft.content,
      categoryId: draft.categoryId ? Number(draft.categoryId) : 0,
      labelIds: draft.labelIds,
      visibility: draft.visibility,
      visibleRoleIds: draft.visibility === 'roles' ? draft.visibleRoleIds : [],
      top: draft.top,
    }

    if (isCreating || publish)
      request.publish = publish

    setSaving(true)

    try {
      const savedPost = isCreating
        ? await createPost(request)
        : await updatePost(post.id, request)

      const savedDraft = createInitialDraft(savedPost)
      savedDraftRef.current = JSON.stringify(savedDraft)
      draftRef.current = savedDraft
      localStorage.removeItem(storageKey)
      setDraft(savedDraft)
      setAutoSaveMessage(publish ? t('publishedMessage') : t('savedToDrafts'))
      toast.success(publish ? t('publishedMessage') : t('draftSaved'))

      if (isCreating) {
        router.replace(`/console/posts/edit/${savedPost.id}`)
      }
      else {
        router.refresh()
      }
    }
    catch {
      toast.error(isCreating ? t('createFailed') : t('saveFailed'))
    }
    finally {
      setSaving(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="mx-auto max-w-5xl space-y-6">
      <header className="flex flex-wrap items-center justify-between gap-4 border-b border-black/10 pb-5 dark:border-white/10">
        <div>
          <h1 className="text-2xl font-semibold">
            {isCreating ? t('createTitle') : t('editTitle')}
          </h1>
          {!isCreating && (
            <p className="mt-1 text-sm text-neutral-500">
              {post.status === 'published' ? t('published') : t('draft')}
            </p>
          )}
          {autoSaveMessage && (
            <p className="mt-1 text-xs text-neutral-500">{autoSaveMessage}</p>
          )}
        </div>
      </header>

      <label className="block space-y-2">
        <span className="text-sm">{t('titleField')}</span>
        <input
          name="title"
          value={draft.title}
          onChange={event => setDraft(current => ({
            ...current,
            title: event.target.value,
          }))}
          required
          maxLength={255}
          className="min-h-11 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
        />
      </label>

      <div className="grid gap-4 md:grid-cols-2">
        <label className="block space-y-2">
          <span className="text-sm">{t('type')}</span>
          <select
            name="type"
            value={draft.type}
            onChange={event => setDraft(current => ({
              ...current,
              type: event.target.value,
            }))}
            className="min-h-10 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
          >
            <option value="article">{t('article')}</option>
            <option value="page">{t('page')}</option>
          </select>
        </label>

        <label className="block space-y-2">
          <span className="text-sm">{t('category')}</span>
          <select
            name="categoryId"
            value={draft.categoryId}
            onChange={event => setDraft(current => ({
              ...current,
              categoryId: event.target.value,
            }))}
            className="min-h-10 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
          >
            <option value="">{t('uncategorized')}</option>
            {categories.map(category => (
              <option key={category.id} value={category.id}>{category.name}</option>
            ))}
          </select>
        </label>

        <label className="block space-y-2">
          <span className="text-sm">{t('cover')}</span>
          <input
            name="cover"
            type="url"
            value={draft.cover}
            onChange={event => setDraft(current => ({
              ...current,
              cover: event.target.value,
            }))}
            maxLength={512}
            className="min-h-10 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
          />
        </label>
      </div>

      <label className="block space-y-2">
        <span className="text-sm">{t('description')}</span>
        <textarea
          name="description"
          value={draft.description}
          onChange={event => setDraft(current => ({
            ...current,
            description: event.target.value,
          }))}
          rows={3}
          maxLength={1000}
          className="w-full rounded-md border border-black/15 p-3 dark:border-white/15"
        />
      </label>

      {labels.length > 0 && (
        <fieldset className="space-y-3">
          <legend className="text-sm">{t('labels')}</legend>
          <div className="flex flex-wrap gap-3">
            {labels.map(label => (
              <label key={label.id} className="inline-flex items-center gap-2 rounded-sm border border-black/10 px-3 py-2 text-sm dark:border-white/10">
                <input
                  name="labelIds"
                  type="checkbox"
                  value={label.id}
                  checked={draft.labelIds.includes(label.id)}
                  onChange={() => toggleLabel(label.id)}
                />
                {label.name}
              </label>
            ))}
          </div>
        </fieldset>
      )}

      <section aria-labelledby="content-title">
        <h2 id="content-title" className="mb-2 text-sm">{t('content')}</h2>
        <SimpleEditor
          value={draft.content}
          onChange={content => setDraft(current => ({
            ...current,
            content,
          }))}
          disabled={saving}
        />
      </section>

      <div className="space-y-4">
        <label className="block max-w-sm space-y-2">
          <span className="text-sm">{t('visibility')}</span>
          <select
            name="visibility"
            value={draft.visibility}
            onChange={event => setDraft(current => ({
              ...current,
              visibility: event.target.value as PostVisibility,
            }))}
            className="min-h-10 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
          >
            <option value="public">{t('public')}</option>
            <option value="roles">{t('roles')}</option>
            <option value="private">{t('private')}</option>
          </select>
        </label>

        {draft.visibility === 'roles' && (
          <fieldset className="space-y-3">
            <legend className="text-sm">{t('visibleRoles')}</legend>
            <div className="flex flex-wrap gap-3">
              {roleOptions.map(role => (
                <label key={role.id} className="inline-flex items-center gap-2 rounded-sm border border-black/10 px-3 py-2 text-sm dark:border-white/10">
                  <input
                    name="visibleRoleIds"
                    type="checkbox"
                    value={role.id}
                    checked={draft.visibleRoleIds.includes(role.id)}
                    onChange={() => toggleVisibleRole(role.id)}
                  />
                  {role.name}
                </label>
              ))}
            </div>
          </fieldset>
        )}

        <div className="flex flex-wrap gap-6">
          <label className="flex items-center gap-2">
            <input
              name="top"
              type="checkbox"
              checked={draft.top}
              onChange={event => setDraft(current => ({
                ...current,
                top: event.target.checked,
              }))}
            />
            {t('top')}
          </label>
        </div>
      </div>

      <footer className="flex flex-wrap justify-end gap-3 border-t border-black/10 pt-5 dark:border-white/10">
        <button
          type="button"
          onClick={handlePreview}
          disabled={saving}
          className="inline-flex min-h-10 items-center gap-2 rounded-md border border-black/15 px-4 text-sm disabled:opacity-50 dark:border-white/15"
        >
          <Eye className="size-4" aria-hidden="true" />
          {t('preview')}
        </button>

        <button
          type="submit"
          value="draft"
          disabled={saving}
          className="inline-flex min-h-10 items-center gap-2 rounded-md border border-black/15 px-4 text-sm disabled:opacity-50 dark:border-white/15"
        >
          <Save className="size-4" aria-hidden="true" />
          {t('saveDraft')}
        </button>

        <button
          type="submit"
          value="publish"
          disabled={saving}
          className="inline-flex min-h-10 items-center gap-2 rounded-md bg-black px-4 text-sm text-white disabled:opacity-50 dark:bg-white dark:text-black"
        >
          <Send className="size-4" aria-hidden="true" />
          {t('publish')}
        </button>
      </footer>

      <PostPreviewDialog
        preview={preview}
        onClose={() => setPreview(null)}
      />
    </form>
  )
}
