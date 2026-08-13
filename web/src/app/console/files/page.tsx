export default function FilesPage() {
  return (
    <section aria-labelledby="files-title" className="space-y-6">
      <header>
        <h1 id="files-title" className="text-2xl font-semibold">
          文件管理
        </h1>
        <p className="mt-1 text-sm text-neutral-500">
          文件接口完成后，这里将展示文件列表和上传操作。
        </p>
      </header>

      <div className="border-y border-black/10 py-12 text-center text-sm text-neutral-500 dark:border-white/10">
        暂无文件数据
      </div>
    </section>
  )
}
