import type { LucideIcon } from 'lucide-react'
import {
  FilePenLine,
  FileText,
  MessageSquare,
  Users,
} from 'lucide-react'

interface Statistic {
  label: string
  value: string
  icon: LucideIcon
}

const statistics: Statistic[] = [
  {
    label: '文章总数',
    value: '--',
    icon: FileText,
  },
  {
    label: '评论总数',
    value: '--',
    icon: MessageSquare,
  },
  {
    label: '用户总数',
    value: '--',
    icon: Users,
  },
  {
    label: '草稿总数',
    value: '--',
    icon: FilePenLine,
  },
]

export function Dashboard() {
  return (
    <div className="space-y-8">
      <section aria-labelledby="dashboard-title">
        <h1 id="dashboard-title" className="text-2xl font-semibold">
          数据概览
        </h1>

        <div className="mt-6 grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {statistics.map(({
            icon: Icon,
            label,
            value,
          }) => (
            <div
              key={label}
              className="min-h-24 rounded-lg border border-black/10 p-4 dark:border-white/10"
            >
              <div className="flex items-center justify-between gap-4">
                <span className="text-sm text-neutral-500">
                  {label}
                </span>
                <Icon aria-hidden="true" size={18} />
              </div>
              <p className="mt-3 text-2xl font-semibold">
                {value}
              </p>
            </div>
          ))}
        </div>
      </section>

      <div className="grid gap-8 xl:grid-cols-2">
        <section aria-labelledby="recent-posts-title">
          <h2 id="recent-posts-title" className="text-lg font-semibold">
            最近文章
          </h2>
          <div className="mt-4 border-t border-black/10 py-6 text-sm text-neutral-500 dark:border-white/10">
            暂无文章数据
          </div>
        </section>

        <section aria-labelledby="recent-comments-title">
          <h2 id="recent-comments-title" className="text-lg font-semibold">
            最近评论
          </h2>
          <div className="mt-4 border-t border-black/10 py-6 text-sm text-neutral-500 dark:border-white/10">
            暂无评论数据
          </div>
        </section>
      </div>
    </div>
  )
}
