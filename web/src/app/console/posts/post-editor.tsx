'use client'

import type { FormEvent } from 'react'
import type {
  Category,
  CreatePostReq,
  Label,
  Post,
  PostVisibility,
} from '@/models/post'
import type { RoleOption } from '@/models/role'
import { EditorContent, useEditor } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import { Bold, Heading2, Italic, Save, Send } from 'lucide-react'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { toast } from 'sonner'
import { createPost, updatePost } from '@/api/post'

interface PostEditorProps {
  post?: Post
  categories: Category[]
  labels: Label[]
  roleOptions: RoleOption[]
}

export function PostEditor({
  post,
  categories,
  labels,
  roleOptions,
}: PostEditorProps) {
  const router = useRouter()
  const [saving, setSaving] = useState(false)
  const [visibility, setVisibility] = useState<PostVisibility>(
    post?.visibility ?? (post?.isPrivate ? 'private' : 'public'),
  )
  const isCreating = post === undefined

  const editor = useEditor({
    extensions: [StarterKit],
    content: post?.draftContent ?? post?.content ?? '',
    immediatelyRender: false,
  })

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (!editor || saving)
      return

    const form = new FormData(event.currentTarget)
    const submitter = (event.nativeEvent as SubmitEvent)
      .submitter as HTMLButtonElement | null
    const publish = submitter?.value === 'publish'
    const title = String(form.get('title') ?? '').trim()

    if (!title) {
      toast.error('请输入文章标题')
      return
    }

    if (publish && editor.isEmpty) {
      toast.error('文章正文不能为空')
      return
    }

    const categoryValue = String(form.get('categoryId') ?? '')
    const visibleRoleIds = form
      .getAll('visibleRoleIds')
      .map(value => Number(value))

    if (visibility === 'roles' && visibleRoleIds.length === 0) {
      toast.error('请选择至少一个可见身份')
      return
    }

    const request: CreatePostReq = {
      title,
      description: String(form.get('description') ?? '').trim(),
      cover: String(form.get('cover') ?? '').trim(),
      type: String(form.get('type') ?? 'article'),
      draftContent: editor.getHTML(),
      categoryId: categoryValue ? Number(categoryValue) : 0,
      labelIds: form.getAll('labelIds').map(value => Number(value)),
      visibility,
      visibleRoleIds: visibility === 'roles' ? visibleRoleIds : [],
      top: form.has('top'),
    }

    if (isCreating || publish)
      request.publish = publish

    setSaving(true)

    try {
      const savedPost = isCreating
        ? await createPost(request)
        : await updatePost(post.id, request)

      toast.success(publish ? '文章已发布' : '草稿已保存')

      if (isCreating) {
        router.replace(`/console/posts/edit/${savedPost.id}`)
      }
      else {
        router.refresh()
      }
    }
    catch {
      toast.error(isCreating ? '创建文章失败' : '保存文章失败')
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
            {isCreating ? '新建文章' : '编辑文章'}
          </h1>
          {!isCreating && (
            <p className="mt-1 text-sm text-neutral-500">
              {post.status === 'published' ? '已发布' : '草稿'}
            </p>
          )}
        </div>
      </header>

      <label className="block space-y-2">
        <span className="text-sm">标题</span>
        <input
          name="title"
          defaultValue={post?.title ?? ''}
          required
          maxLength={255}
          className="min-h-11 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
        />
      </label>

      <div className="grid gap-4 md:grid-cols-2">
        <label className="block space-y-2">
          <span className="text-sm">类型</span>
          <select
            name="type"
            defaultValue={post?.type ?? 'article'}
            className="min-h-10 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
          >
            <option value="article">文章</option>
            <option value="page">页面</option>
          </select>
        </label>

        <label className="block space-y-2">
          <span className="text-sm">分类</span>
          <select
            name="categoryId"
            defaultValue={post?.categoryId ?? ''}
            className="min-h-10 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
          >
            <option value="">未分类</option>
            {categories.map(category => (
              <option key={category.id} value={category.id}>{category.name}</option>
            ))}
          </select>
        </label>

        <label className="block space-y-2">
          <span className="text-sm">封面地址</span>
          <input
            name="cover"
            type="url"
            defaultValue={post?.cover ?? ''}
            maxLength={512}
            className="min-h-10 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
          />
        </label>
      </div>

      <label className="block space-y-2">
        <span className="text-sm">摘要</span>
        <textarea
          name="description"
          defaultValue={post?.description ?? ''}
          rows={3}
          maxLength={1000}
          className="w-full rounded-md border border-black/15 p-3 dark:border-white/15"
        />
      </label>

      {labels.length > 0 && (
        <fieldset className="space-y-3">
          <legend className="text-sm">标签</legend>
          <div className="flex flex-wrap gap-3">
            {labels.map(label => (
              <label key={label.id} className="inline-flex items-center gap-2 rounded-sm border border-black/10 px-3 py-2 text-sm dark:border-white/10">
                <input
                  name="labelIds"
                  type="checkbox"
                  value={label.id}
                  defaultChecked={post?.labels.some(item => item.id === label.id) ?? false}
                />
                {label.name}
              </label>
            ))}
          </div>
        </fieldset>
      )}

      <section aria-labelledby="content-title">
        <h2 id="content-title" className="mb-2 text-sm">正文</h2>

        <div className="rounded-md border border-black/15 dark:border-white/15">
          <div
            role="toolbar"
            aria-label="正文格式"
            className="flex gap-1 border-b border-black/15 p-2 dark:border-white/15"
          >
            <button
              type="button"
              aria-label="加粗"
              title="加粗"
              onClick={() => editor?.chain().focus().toggleBold().run()}
              className="inline-flex size-9 items-center justify-center"
            >
              <Bold className="size-4" aria-hidden="true" />
            </button>
            <button
              type="button"
              aria-label="斜体"
              title="斜体"
              onClick={() => editor?.chain().focus().toggleItalic().run()}
              className="inline-flex size-9 items-center justify-center"
            >
              <Italic className="size-4" aria-hidden="true" />
            </button>
            <button
              type="button"
              aria-label="二级标题"
              title="二级标题"
              onClick={() => editor?.chain().focus().toggleHeading({ level: 2 }).run()}
              className="inline-flex size-9 items-center justify-center"
            >
              <Heading2 className="size-4" aria-hidden="true" />
            </button>
          </div>

          <EditorContent
            editor={editor}
            className="[&_.ProseMirror]:min-h-80 [&_.ProseMirror]:p-4 [&_.ProseMirror]:outline-none"
          />
        </div>
      </section>

      <div className="space-y-4">
        <label className="block max-w-sm space-y-2">
          <span className="text-sm">可见范围</span>
          <select
            name="visibility"
            value={visibility}
            onChange={event => setVisibility(event.target.value as PostVisibility)}
            className="min-h-10 w-full rounded-md border border-black/15 px-3 dark:border-white/15"
          >
            <option value="public">公开</option>
            <option value="roles">指定身份</option>
            <option value="private">仅作者和管理员</option>
          </select>
        </label>

        {visibility === 'roles' && (
          <fieldset className="space-y-3">
            <legend className="text-sm">可见身份</legend>
            <div className="flex flex-wrap gap-3">
              {roleOptions.map(role => (
                <label key={role.id} className="inline-flex items-center gap-2 rounded-sm border border-black/10 px-3 py-2 text-sm dark:border-white/10">
                  <input
                    name="visibleRoleIds"
                    type="checkbox"
                    value={role.id}
                    defaultChecked={post?.visibleRoles.some(item => item.id === role.id) ?? false}
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
              defaultChecked={post?.top ?? false}
            />
            置顶文章
          </label>
        </div>
      </div>

      <footer className="flex flex-wrap justify-end gap-3 border-t border-black/10 pt-5 dark:border-white/10">
        <button
          type="submit"
          value="draft"
          disabled={saving}
          className="inline-flex min-h-10 items-center gap-2 rounded-md border border-black/15 px-4 text-sm disabled:opacity-50 dark:border-white/15"
        >
          <Save className="size-4" aria-hidden="true" />
          保存草稿
        </button>

        <button
          type="submit"
          value="publish"
          disabled={saving}
          className="inline-flex min-h-10 items-center gap-2 rounded-md bg-black px-4 text-sm text-white disabled:opacity-50 dark:bg-white dark:text-black"
        >
          <Send className="size-4" aria-hidden="true" />
          发布
        </button>
      </footer>
    </form>
  )
}
