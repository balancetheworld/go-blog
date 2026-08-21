import type { Metadata } from 'next'
import { getTranslations } from 'next-intl/server'
import { notFound } from 'next/navigation'
import {
  getPost,
  PostServerError,
} from '@/api/post.server'
import { CommentSection } from '@/components/comment'
import { sanitizePostHTML } from '@/lib/post-html'
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
  const t = await getTranslations('Post')
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
        <main className="article-detail-main">
          <section className="article-detail-shell article-empty visible">
            <h1>{t('forbiddenTitle')}</h1>
            <p>{t('forbiddenDescription')}</p>
          </section>
        </main>
      )
    }

    throw error
  }

  if (!post) {
    notFound()
  }

  const sanitizedPost = {
    ...post,
    content: sanitizePostHTML(post.content),
  }

  return (
    <main className="article-detail-main">
      <BlogPost post={sanitizedPost}>
        <CommentSection
          targetType={post.type === 'page' ? 'page' : 'post'}
          targetID={post.id}
          targetAuthorID={post.author.id}
        />
      </BlogPost>
    </main>
  )
}
