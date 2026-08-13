import Link from 'next/link'

export function ConsoleAccessDenied() {
  return (
    <main className="flex min-h-screen items-center justify-center p-6">
      <div className="w-full max-w-md border-y border-black/10 py-10 text-center dark:border-white/10">
        <h1 className="text-2xl font-semibold">无权访问后台</h1>
        <p className="mt-3 text-sm text-neutral-500">
          当前账号没有管理员或编辑者权限。
        </p>
        <Link
          href="/"
          className="mt-6 inline-flex min-h-10 items-center rounded-md bg-black px-4 text-sm text-white dark:bg-white dark:text-black"
        >
          返回首页
        </Link>
      </div>
    </main>
  )
}
