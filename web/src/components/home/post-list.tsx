import type { PostListItem } from '@/models/post'
import { Eye, MessageCircle, Pin } from 'lucide-react'
import Image from 'next/image'
import Link from 'next/link'

interface PostListProps {
  posts: PostListItem[]
}

const dateFormatter = new Intl.DateTimeFormat('zh-CN', {
  dateStyle: 'medium',
})

export default function PostList({ posts }: PostListProps) {
  if (posts.length === 0) {
    return (
      <div className="border-y border-black/10 py-14 text-center text-neutral-500 dark:border-white/10">
        暂无符合条件的文章
      </div>
    )
  }

  return (
    <div className="divide-y divide-black/10 border-y border-black/10 dark:divide-white/10 dark:border-white/10">
      {posts.map(post => (
        <article key={post.id} className="grid gap-5 py-7 sm:grid-cols-[minmax(0,1fr)_180px]">
          <div className="min-w-0">
            <div className="flex flex-wrap items-center gap-2 text-sm text-neutral-500">
              {post.category && <span>{post.category.name}</span>}
              {post.top && (
                <span className="inline-flex items-center gap-1">
                  <Pin className="size-3.5" aria-hidden="true" />
                  置顶
                </span>
              )}
            </div>

            <h2 className="mt-2 text-xl font-semibold leading-8">
              <Link href={`/p/${post.slug || post.id}`}>{post.title}</Link>
            </h2>

            {post.description && (
              <p className="mt-2 line-clamp-2 leading-7 text-neutral-600 dark:text-neutral-400">
                {post.description}
              </p>
            )}

            {post.labels.length > 0 && (
              <div className="mt-3 flex flex-wrap gap-2">
                {post.labels.map(label => (
                  <Link
                    key={label.id}
                    href={`/?label=${label.id}`}
                    className="rounded-sm border border-black/10 px-2 py-1 text-xs text-neutral-600 dark:border-white/10 dark:text-neutral-400"
                  >
                    {label.name}
                  </Link>
                ))}
              </div>
            )}

            <div className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-neutral-500">
              <span>{post.author.nickname || post.author.username}</span>
              <time dateTime={post.publishedAt ?? post.createdAt}>
                {dateFormatter.format(new Date(post.publishedAt ?? post.createdAt))}
              </time>
              <span className="inline-flex items-center gap-1">
                <Eye className="size-4" aria-hidden="true" />
                {post.viewCount}
              </span>
              <span className="inline-flex items-center gap-1">
                <MessageCircle className="size-4" aria-hidden="true" />
                {post.commentCount}
              </span>
            </div>
          </div>

          {post.cover && (
            <Link href={`/p/${post.slug || post.id}`} className="relative order-first aspect-video overflow-hidden rounded-md sm:order-none sm:aspect-[3/2]">
              <Image
                src={post.cover}
                alt=""
                fill
                sizes="(max-width: 640px) 100vw, 180px"
                className="object-cover"
              />
            </Link>
          )}
        </article>
      ))}
    </div>
  )
}
