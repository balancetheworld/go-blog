import type { PostDetail } from '@/models/post'

interface BlogPostProps {
  post: PostDetail
}

const dateFormatter = new Intl.DateTimeFormat('zh-CN', {
  dateStyle: 'long',
})

export function BlogPost({ post }: BlogPostProps) {
  return (
    <article className="mx-auto w-full max-w-3xl">
      <header className="border-b border-black/10 pb-6 dark:border-white/10">
        <h1 className="text-3xl font-semibold leading-tight">
          {post.title}
        </h1>

        <div className="mt-4 flex flex-wrap gap-x-4 gap-y-2 text-sm text-neutral-500">
          <span>{post.authorName}</span>
          <time dateTime={post.publishedAt}>
            {dateFormatter.format(new Date(post.publishedAt))}
          </time>
        </div>
      </header>

      <div className="whitespace-pre-wrap break-words py-8 leading-8">
        {post.content}
      </div>
    </article>
  )
}
