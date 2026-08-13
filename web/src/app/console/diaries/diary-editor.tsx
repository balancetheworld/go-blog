'use client'

import type { FormEvent } from 'react'
import type {
  CreateDiaryRequest,
  Diary,
  DiaryFolder,
  DiaryVisibility,
} from '@/models/diary'
import type { RoleOption } from '@/models/role'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { createDiary, updateDiary } from '@/api/diary'
import { SimpleEditor } from '@/components/tiptap-editor/simple-editor'
import { isHTMLContentEmpty } from '@/lib/tiptap-advanced-utils'

interface DiaryEditorProps {
  diary?: Diary
  folders: DiaryFolder[]
  roleOptions: RoleOption[]
}

interface DiaryDraft {
  title: string
  description: string
  cover: string
  folderId: string
  visibility: DiaryVisibility
  visibleRoleIds: number[]
  content: string
}

function createInitialDraft(diary?: Diary): DiaryDraft {
  return {
    title: diary?.title ?? '',
    description: diary?.description ?? '',
    cover: diary?.cover ?? '',
    folderId: diary?.folder ? String(diary.folder.id) : '',
    visibility: diary?.visibility ?? 'public',
    visibleRoleIds: diary?.visibleRoles.map(role => role.id) ?? [],
    content: diary?.draftContent ?? diary?.content ?? '',
  }
}

export function DiaryEditor({
  diary,
  folders,
  roleOptions,
}: DiaryEditorProps) {
  const router = useRouter()
  const [draft, setDraft] = useState(() => createInitialDraft(diary))
  const [saving, setSaving] = useState(false)
  const isCreating = diary === undefined

  function toggleVisibleRole(roleID: number) {
    setDraft(current => ({
      ...current,
      visibleRoleIds: current.visibleRoleIds.includes(roleID)
        ? current.visibleRoleIds.filter(id => id !== roleID)
        : [...current.visibleRoleIds, roleID],
    }))
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (saving)
      return

    const submitter = (event.nativeEvent as SubmitEvent)
      .submitter as HTMLButtonElement | null
    const publish = submitter?.value === 'publish'

    if (publish && isHTMLContentEmpty(draft.content)) {
      toast.error('日记正文不能为空')
      return
    }
    if (draft.visibility === 'roles' && draft.visibleRoleIds.length === 0) {
      toast.error('请选择至少一个可见身份')
      return
    }

    const request: CreateDiaryRequest = {
      title: draft.title.trim(),
      description: draft.description.trim(),
      cover: draft.cover.trim(),
      folderId: draft.folderId ? Number(draft.folderId) : undefined,
      draftContent: draft.content,
      publish,
      visibility: draft.visibility,
      visibleRoleIds: draft.visibility === 'roles' ? draft.visibleRoleIds : [],
    }

    setSaving(true)
    try {
      const savedDiary = isCreating
        ? await createDiary(request)
        : await updateDiary(diary.id, {
            ...request,
            publish: publish ? true : undefined,
            clearFolder: !draft.folderId,
          })

      toast.success(publish ? '日记已发布' : '草稿已保存')
      if (isCreating)
        router.replace(`/console/diaries/edit/${savedDiary.id}`)
      else
        router.refresh()
    }
    catch {
      toast.error(isCreating ? '创建日记失败' : '保存日记失败')
    }
    finally {
      setSaving(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="mx-auto max-w-5xl space-y-6">
      <header className="border-b border-black/10 pb-5 dark:border-white/10">
        <h1 className="text-2xl font-semibold">
          {isCreating ? '新建日记' : '编辑日记'}
        </h1>
        {!isCreating && (
          <p className="mt-1 text-sm text-neutral-500">
            {diary.status === 'published' ? '已发布' : '草稿'}
          </p>
        )}
      </header>

      <label className="block space-y-2">
        <span className="text-sm">标题</span>
        <input
          value={draft.title}
          maxLength={200}
          onChange={event => setDraft(current => ({ ...current, title: event.target.value }))}
          className="min-h-11 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
        />
      </label>

      <label className="block space-y-2">
        <span className="text-sm">摘要</span>
        <textarea
          value={draft.description}
          maxLength={500}
          rows={3}
          onChange={event => setDraft(current => ({ ...current, description: event.target.value }))}
          className="w-full resize-y rounded-md border border-black/15 p-3 dark:border-white/15"
        />
      </label>

      <div className="grid gap-4 md:grid-cols-2">
        <label className="block space-y-2">
          <span className="text-sm">封面地址</span>
          <input
            value={draft.cover}
            maxLength={500}
            onChange={event => setDraft(current => ({ ...current, cover: event.target.value }))}
            className="min-h-11 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
          />
        </label>
        <label className="block space-y-2">
          <span className="text-sm">文件夹</span>
          <select
            value={draft.folderId}
            onChange={event => setDraft(current => ({ ...current, folderId: event.target.value }))}
            className="min-h-11 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
          >
            <option value="">未分类</option>
            {folders.map(folder => (
              <option key={folder.id} value={folder.id}>{folder.name}</option>
            ))}
          </select>
        </label>
      </div>

      <fieldset className="space-y-3 border-y border-black/10 py-5 dark:border-white/10">
        <legend className="text-sm">可见范围</legend>
        <div className="flex flex-wrap gap-4 text-sm">
          {(['public', 'roles', 'private'] as const).map(value => (
            <label key={value} className="inline-flex items-center gap-2">
              <input
                type="radio"
                name="visibility"
                value={value}
                checked={draft.visibility === value}
                onChange={() => setDraft(current => ({ ...current, visibility: value }))}
              />
              {value === 'public' ? '公开' : value === 'roles' ? '指定身份' : '私密'}
            </label>
          ))}
        </div>
        {draft.visibility === 'roles' && (
          <div className="flex flex-wrap gap-3">
            {roleOptions.map(role => (
              <label key={role.id} className="inline-flex items-center gap-2 text-sm">
                <input
                  type="checkbox"
                  checked={draft.visibleRoleIds.includes(role.id)}
                  onChange={() => toggleVisibleRole(role.id)}
                />
                {role.name}
              </label>
            ))}
          </div>
        )}
      </fieldset>

      <div className="space-y-2">
        <span className="text-sm">正文</span>
        <SimpleEditor
          value={draft.content}
          onChange={content => setDraft(current => ({ ...current, content }))}
          minHeight={420}
          ariaLabel="日记正文编辑器"
        />
      </div>

      <div className="flex justify-end gap-3 border-t border-black/10 pt-5 dark:border-white/10">
        <button
          type="submit"
          value="draft"
          disabled={saving}
          className="min-h-10 rounded-md border border-black/15 px-4 text-sm disabled:opacity-50 dark:border-white/15"
        >
          保存草稿
        </button>
        <button
          type="submit"
          value="publish"
          disabled={saving}
          className="min-h-10 rounded-md bg-black px-4 text-sm text-white disabled:opacity-50 dark:bg-white dark:text-black"
        >
          {saving ? '保存中' : '发布'}
        </button>
      </div>
    </form>
  )
}
