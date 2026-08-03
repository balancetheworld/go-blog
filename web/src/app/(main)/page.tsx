import type { PostSummary } from '@/models/post'
import PostList from '@/components/home/post-list'

const placeholderPosts: PostSummary[] = [
  {
    id: 1,
    title: '博客项目的第一篇文章',
    excerpt: '记录这个博客从后端接口到前端页面的搭建过程。',
    authorName: '站长',
    publishedAt: '2026-08-03T08:00:00+09:00',
  },
  {
    id: 2,
    title: '为什么选择 Next.js 和 Hertz',
    excerpt: '介绍前后端技术选型，以及两个应用之间如何通信。',
    authorName: '站长',
    publishedAt: '2026-08-02T08:00:00+09:00',
  },
]

export default function HomePage() {
  return (
    <section>
      <header className="mb-8">
        <h1 className="text-2xl font-semibold">最新文章</h1>
        <p className="mt-2 text-neutral-600 dark:text-neutral-400">
          浏览最近发布的公开内容
        </p>
      </header>

      <PostList posts={placeholderPosts} />
    </section>
  )
}
