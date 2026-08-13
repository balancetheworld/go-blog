import type { Metadata } from 'next'
import { Eye, LockKeyhole, MessageCircle, Users } from 'lucide-react'
import Image from 'next/image'
import { notFound } from 'next/navigation'
import { DiaryServerError, getDiary } from '@/api/diary.server'
import { CommentSection } from '@/components/comment'
import { sanitizePostHTML } from '@/lib/post-html'

interface DiaryDetailPageProps {
  params: Promise<{
    id: string
  }>
}

const dateFormatter = new Intl.DateTimeFormat('zh-CN', {
  dateStyle: 'long',
})

export async function generateMetadata({
  params,
}: DiaryDetailPageProps): Promise<Metadata> {
  const { id } = await params

  try {
    const diary = await getDiary(id)
    return {
      title: diary.title || '日记',
      description: diary.description,
    }
  }
  catch {
    return {}
  }
}

export default async function DiaryDetailPage({
  params,
}: DiaryDetailPageProps) {
  const { id } = await params

  let diary
  try {
    diary = await getDiary(id)
  }
  catch (error) {
    if (
      error instanceof DiaryServerError
      && [401, 403, 404].includes(error.status)
    ) {
      notFound()
    }

    throw error
  }

  const authorName = diary.author.nickname || diary.author.username
  const content = sanitizePostHTML(diary.content)

  return (
    <>
      <article className="mx-auto w-full max-w-3xl">
        <header className="border-b border-black/10 pb-8 dark:border-white/10">
          <div className="flex flex-wrap items-center gap-3 text-sm text-neutral-500">
            {diary.folder && <span>{diary.folder.name}</span>}
            {diary.visibility === 'private' && (
              <span className="inline-flex items-center gap-1">
                <LockKeyhole className="size-4" aria-hidden="true" />
                私密
              </span>
            )}
            {diary.visibility === 'roles' && (
              <span className="inline-flex items-center gap-1">
                <Users className="size-4" aria-hidden="true" />
                {diary.visibleRoles.map(role => role.name).join('、')}
              </span>
            )}
          </div>

          <h1 className="mt-4 text-3xl font-semibold leading-tight">
            {diary.title || '无标题日记'}
          </h1>
          {diary.description && (
            <p className="mt-4 leading-7 text-neutral-600 dark:text-neutral-400">
              {diary.description}
            </p>
          )}
          <div className="mt-5 flex flex-wrap gap-x-4 gap-y-2 text-sm text-neutral-500">
            <span>{authorName}</span>
            <time dateTime={diary.publishedAt ?? diary.createdAt}>
              {dateFormatter.format(new Date(diary.publishedAt ?? diary.createdAt))}
            </time>
            <span className="inline-flex items-center gap-1">
              <Eye className="size-4" aria-hidden="true" />
              {diary.viewCount}
            </span>
            <span className="inline-flex items-center gap-1">
              <MessageCircle className="size-4" aria-hidden="true" />
              {diary.commentCount}
            </span>
          </div>
          {diary.cover && (
            <div className="relative mt-8 aspect-video overflow-hidden rounded-md">
              <Image
                src={diary.cover}
                alt={diary.title || '日记封面'}
                fill
                priority
                sizes="(max-width: 768px) 100vw, 768px"
                className="object-cover"
              />
            </div>
          )}
        </header>

        <div
          className="break-words py-8 leading-8 [&_a]:underline [&_blockquote]:my-6 [&_blockquote]:border-l-2 [&_blockquote]:pl-4 [&_code]:rounded-sm [&_code]:bg-black/5 [&_code]:px-1 [&_h1]:mt-10 [&_h1]:text-3xl [&_h2]:mt-9 [&_h2]:text-2xl [&_h3]:mt-8 [&_h3]:text-xl [&_hr]:my-8 [&_img]:my-8 [&_img]:max-w-full [&_img]:rounded-md [&_img[data-align=center]]:mx-auto [&_img[data-align=left]]:mr-auto [&_img[data-align=right]]:ml-auto [&_li]:my-1 [&_ol]:my-5 [&_ol]:list-decimal [&_ol]:pl-6 [&_p]:my-5 [&_pre]:my-6 [&_pre]:overflow-x-auto [&_pre]:rounded-md [&_pre]:bg-neutral-900 [&_pre]:p-4 [&_pre]:text-neutral-100 [&_ul]:my-5 [&_ul]:list-disc [&_ul]:pl-6 dark:[&_code]:bg-white/10"
          dangerouslySetInnerHTML={{ __html: content }}
        />
      </article>

      <CommentSection
        targetType="diary"
        targetID={diary.id}
        targetAuthorID={diary.author.id}
      />
    </>
  )
}
