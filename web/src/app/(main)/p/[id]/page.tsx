import type { Metadata } from 'next'
import { notFound } from 'next/navigation'
import { CommentInput } from '@/components/comment/comment-input'
import {
  getPublicPost,
  PostUnavailableError,
} from '@/lib/api/post'
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

  // 正则校验id是否为合法正整数（不能以0开头、只能数字）
  // ^\[1-9\] 首位1-9，\\d\* 后面任意数字，$结尾
  if (!/^[1-9]\d*$/.test(id)) {
    // id格式非法，返回空元数据，使用网站默认SEO
    return {}
  }
  try {
    const post = await getPublicPost(id)

    return post
      ? {
          title: post.title,
          description: post.content.slice(0, 120),
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
  if (!/^[1-9]\d*$/.test(id)) {
    notFound()
  }

  let post
  try {
    post = await getPublicPost(id)
  }
  catch (error) {
    if (error instanceof PostUnavailableError) {
      notFound()
    }

    throw error
  }

  if (!post) {
    notFound()
  }

  return (
    <>
      <BlogPost post={post} />
      <CommentInput />
    </>
  )
}
