import Link from 'next/link'

export default function PostNotFound() {
  return (
    <section>
      <h1 className="text-2xl font-semibold">
        文章不可访问
      </h1>
      <p className="mt-3 text-neutral-600 dark:text-neutral-400">
        文章不存在、尚未公开或你没有查看权限。
      </p>
      <Link href="/" className="mt-6 inline-block font-medium">
        返回首页
      </Link>
    </section>
  )
}
