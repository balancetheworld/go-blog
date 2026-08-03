export interface PostSummary {
  id: number
  title: string
  excerpt: string
  authorName: string
  publishedAt: string
}

export interface PostDetail {
  id: number
  title: string
  content: string
  authorName: string
  publishedAt: string
  updatedAt: string
}
