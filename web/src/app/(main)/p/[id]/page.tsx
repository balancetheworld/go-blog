import type { Metadata } from 'next'
import { notFound } from 'next/navigation'
import {
  getPost,
  PostServerError,
} from '@/api/post.server'
import { CommentSection } from '@/components/comment'
import { BlogPost } from './blog-post'

// 定义页面组件入参的TS接口 PostPageProps
interface PostPageProps {
  // params：Next.js 动态路由参数，App Router 下 params 是 Promise 类型
  params: Promise<{
    id: string // 动态路由路径里的id参数，字符串格式
  }>
}

export async function generateMetadata({
  params,
}: PostPageProps): Promise<Metadata> {
  // 新版Next.js params是Promise，必须await解析，拿到路径里的文章id
  const { id } = await params

  try {
    const post = await getPost(id)

    return post
      ? {
          title: post.title,
          description: post.description,
        }
      : {}
  }
  catch {
    return {}
  }
}

export default async function PostPage({
  params,
}: PostPageProps) {
  const { id } = await params

  let post
  try {
    post = await getPost(id)
  }
  catch (error) {
    if (error instanceof PostServerError && error.status === 404) {
      notFound()
    }

    if (
      error instanceof PostServerError
      && (error.status === 401 || error.status === 403)
    ) {
      return (
        <section className="border-y border-black/10 py-12 dark:border-white/10">
          <h1 className="text-2xl font-semibold">无权访问这篇文章</h1>
          <p className="mt-3 text-neutral-600 dark:text-neutral-400">
            请使用有权限的账号登录后重试。
          </p>
        </section>
      )
    }

    throw error
  }

  if (!post) {
    notFound()
  }

  return (
    <>
      <BlogPost post={post} />
      <CommentSection
        targetType={post.type === 'page' ? 'page' : 'post'}
        targetID={post.id}
        targetAuthorID={post.author.id}
      />
    </>
  )
}
