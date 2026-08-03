import type { PostSummary } from '@/models/post'
import Link from 'next/link'

interface PostListProps {
  posts: PostSummary[]
}

// 创建一个国际化日期格式化工具实例，专门输出中文风格日期
const dateFormatter = new Intl.DateTimeFormat('zh-CN', {
  dateStyle: 'medium',
})

export default function PostList({ posts }: PostListProps) {
  if (posts.length === 0) {
    return <p>...</p>
  }
  return (
    <>
      {posts.map(post => (
        <article key={post.id}>
          <h2>
            <Link href={`/p/${post.id}`}>
              {post.title}
            </Link>
          </h2>

          <p className="mt-2 leading-7 text-neutral-600 dark:text-neutral-400">
            {post.excerpt}
          </p>

          <div className="mt-3 flex flex-wrap gap-x-4 gap-y-1 text-sm text-neutral-500">
            <span>{post.authorName}</span>
            <time dateTime={post.publishedAt}>
              {dateFormatter.format(new Date(post.publishedAt))}
            </time>
          </div>
        </article>
      ))}
    </>
  )
}
