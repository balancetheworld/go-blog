import type { Post } from '@/models/post'
import {
  Eye,
  Heart,
  LockKeyhole,
  MessageCircle,
  Pin,
} from 'lucide-react'
import Image from 'next/image'
import { sanitizePostHTML } from '@/lib/post-html'

interface BlogPostProps {
  post: Post
}

const dateFormatter = new Intl.DateTimeFormat('zh-CN', {
  dateStyle: 'long',
})

export function BlogPost({ post }: BlogPostProps) {
  const authorName
    = post.author.nickname || post.author.username
  const content = sanitizePostHTML(post.content)

  return (
    <article className="mx-auto w-full max-w-3xl">
      <header className="border-b border-black/10 pb-8 dark:border-white/10">
        <div className="flex flex-wrap items-center gap-3 text-sm text-neutral-500">
          {post.category && (
            <span>{post.category.name}</span>
          )}

          {post.top && (
            <span className="inline-flex items-center gap-1">
              <Pin className="size-4" aria-hidden="true" />
              置顶
            </span>
          )}

          {post.isPrivate && (
            <span className="inline-flex items-center gap-1">
              <LockKeyhole className="size-4" aria-hidden="true" />
              私密
            </span>
          )}
        </div>

        <h1 className="mt-4 text-3xl font-semibold leading-tight">
          {post.title}
        </h1>

        {post.description && (
          <p className="mt-4 leading-7 text-neutral-600 dark:text-neutral-400">
            {post.description}
          </p>
        )}

        <div className="mt-5 flex flex-wrap gap-x-4 gap-y-2 text-sm text-neutral-500">
          <span>{authorName}</span>

          <time dateTime={post.publishedAt ?? post.createdAt}>
            {dateFormatter.format(new Date(post.publishedAt ?? post.createdAt))}
          </time>

          <span className="inline-flex items-center gap-1">
            <Eye className="size-4" aria-hidden="true" />
            {post.viewCount}
          </span>

          <span className="inline-flex items-center gap-1">
            <Heart className="size-4" aria-hidden="true" />
            {post.likeCount}
          </span>

          <span className="inline-flex items-center gap-1">
            <MessageCircle className="size-4" aria-hidden="true" />
            {post.commentCount}
          </span>
        </div>

        {post.cover && (
          <div className="relative mt-8 aspect-video overflow-hidden rounded-md">
            <Image
              src={post.cover}
              alt={post.title}
              fill
              priority
              sizes="(max-width: 768px) 100vw, 768px"
              className="object-cover"
            />
          </div>
        )}
      </header>

      <div
        className="break-words py-8 leading-8 [&_a]:underline [&_blockquote]:my-6 [&_blockquote]:border-l-2 [&_blockquote]:pl-4 [&_code]:rounded-sm [&_code]:bg-black/5 [&_code]:px-1 [&_h1]:mt-10 [&_h1]:text-3xl [&_h2]:mt-9 [&_h2]:text-2xl [&_h3]:mt-8 [&_h3]:text-xl [&_hr]:my-8 [&_img]:my-8 [&_img]:max-w-full [&_img]:rounded-md [&_li]:my-1 [&_ol]:my-5 [&_ol]:list-decimal [&_ol]:pl-6 [&_p]:my-5 [&_pre]:my-6 [&_pre]:overflow-x-auto [&_pre]:rounded-md [&_pre]:bg-black/5 [&_pre]:p-4 [&_ul]:my-5 [&_ul]:list-disc [&_ul]:pl-6 dark:[&_code]:bg-white/10 dark:[&_pre]:bg-white/10"
        dangerouslySetInnerHTML={{ __html: content }}
      />

      {post.labels.length > 0 && (
        <footer className="flex flex-wrap gap-2 border-t border-black/10 pt-6 dark:border-white/10">
          {post.labels.map(label => (
            <span
              key={label.id}
              className="rounded-sm border border-black/10 px-2 py-1 text-sm dark:border-white/10"
            >
              {label.name}
            </span>
          ))}
        </footer>
      )}
    </article>
  )
}
