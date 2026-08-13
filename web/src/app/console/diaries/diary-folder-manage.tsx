'use client'

import type { FormEvent } from 'react'
import type { DiaryFolder } from '@/models/diary'
import { FolderPlus, Trash2 } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { createDiaryFolder, deleteDiaryFolder } from '@/api/diary'

interface DiaryFolderManageProps {
  folders: DiaryFolder[]
}

export function DiaryFolderManage({ folders }: DiaryFolderManageProps) {
  const router = useRouter()
  const [name, setName] = useState('')
  const [saving, setSaving] = useState(false)
  const [deletingID, setDeletingID] = useState<number | null>(null)

  async function handleCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const value = name.trim()
    if (!value || saving)
      return

    setSaving(true)
    try {
      await createDiaryFolder({
        name: value,
        description: '',
        cover: '',
        sort: folders.length,
        visibility: 'public',
        visibleRoleIds: [],
      })
      setName('')
      toast.success('文件夹已创建')
      router.refresh()
    }
    catch {
      toast.error('创建文件夹失败')
    }
    finally {
      setSaving(false)
    }
  }

  async function handleDelete(id: number) {
    if (deletingID !== null)
      return

    setDeletingID(id)
    try {
      await deleteDiaryFolder(id)
      toast.success('文件夹已删除')
      router.refresh()
    }
    catch {
      toast.error('删除文件夹失败')
    }
    finally {
      setDeletingID(null)
    }
  }

  return (
    <section aria-labelledby="diary-folders-title" className="space-y-4 border-t border-black/10 pt-6 dark:border-white/10">
      <div>
        <h2 id="diary-folders-title" className="text-lg font-semibold">日记文件夹</h2>
        <p className="mt-1 text-sm text-neutral-500">删除文件夹不会删除其中的日记</p>
      </div>
      <form onSubmit={handleCreate} className="flex max-w-xl gap-2">
        <input
          value={name}
          maxLength={100}
          required
          placeholder="文件夹名称"
          onChange={event => setName(event.target.value)}
          className="min-h-10 min-w-0 flex-1 rounded-md border border-black/15 px-3 dark:border-white/15"
        />
        <button
          type="submit"
          disabled={saving || !name.trim()}
          className="inline-flex min-h-10 items-center gap-2 rounded-md border border-black/15 px-4 text-sm disabled:opacity-50 dark:border-white/15"
        >
          <FolderPlus className="size-4" aria-hidden="true" />
          {saving ? '创建中' : '创建'}
        </button>
      </form>
      {folders.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {folders.map(folder => (
            <div
              key={folder.id}
              className="inline-flex min-h-10 items-center gap-2 rounded-md border border-black/10 px-3 text-sm dark:border-white/10"
            >
              <span>{folder.name}</span>
              <button
                type="button"
                disabled={deletingID !== null}
                onClick={() => void handleDelete(folder.id)}
                aria-label={`删除文件夹 ${folder.name}`}
                title="删除文件夹"
                className="inline-flex size-7 items-center justify-center text-red-600 disabled:opacity-50"
              >
                <Trash2 className="size-4" aria-hidden="true" />
              </button>
            </div>
          ))}
        </div>
      )}
    </section>
  )
}
