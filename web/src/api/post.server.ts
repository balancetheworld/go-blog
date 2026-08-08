// 导入TS类型定义，仅用于类型校验，不会打包进运行代码
import type {
  Category, // 分类数据TS类型
  Label,
  Post, // 单篇帖子详情TS类型
  PostListRequest, // 查询帖子列表的筛选参数类型
  PostListResponse, // 帖子分页列表返回结果类型
} from '@/models/post'
import type { Resp } from '@/models/resp' // 后端统一标准返回体泛型类型

// node内置环境变量对象，读取.env里的后端接口地址
import process from 'node:process'
// 工具函数：把后端返回的下划线蛇形命名字段（category_id）转前端驼峰（categoryId）
import { snakeToCamelObj } from 'field-conv'
// Next.js 13+ 服务端组件API，获取当前请求携带的Cookie（登录凭证）
import { cookies } from 'next/headers'
import { cache } from 'react'

// 读取环境变量中的后端服务地址，无配置则默认本地8888端口
const backendUrl = process.env.BACKEND_URL ?? 'http://localhost:8888'

/**
 * 自定义接口请求错误类
 * 封装后端接口异常信息：HTTP状态码 + 错误提示文本
 * 调用接口出错时抛出该错误，页面/组件统一捕获处理弹窗、跳转404等
 */
export class PostServerError extends Error {
  constructor(
    // HTTP响应状态码 404/403/400/500
    public readonly status: number,
    message: string,
  ) {
    super(message)
    // 自定义错误名称，区分原生Error，方便错误判断
    this.name = 'PostServerError'
  }
}

/**
 * 通用服务端GET请求封装函数
 * 作用：统一处理Next服务端向后端发GET请求的通用逻辑
 * @param path 接口路径，如 /post/list
 * @returns 后端返回的data泛型数据 T
 */
async function serverGet<T>(path: string): Promise<T> {
  // 获取当前页面请求携带的cookie（登录token存放在cookie里，同步传给后端鉴权）
  const cookieStore = await cookies()

  // 发起网络请求，拼接完整后端接口地址
  const response = await fetch(
    `${backendUrl}/api/v1${path}`,
    {
      // no-store：不启用Next数据缓存，每次请求都实时拉取最新数据（适合后台管理、实时内容）
      cache: 'no-store',
      headers: {
        // 把当前用户cookie完整带给后端，实现登录态鉴权
        cookie: cookieStore.toString(),
      },
    },
  )

  // 解析后端返回json，解析失败赋值null防止报错崩溃
  const json = await response.json().catch(() => null)
  // 如果解析成功，将后端蛇形字段统一转为前端驼峰格式，强转为标准返回结构Resp<T>
  const body = json
    ? snakeToCamelObj(json) as Resp<T>
    : null

  // HTTP状态码非2xx（404/403/500等），抛出自定义业务错误
  if (!response.ok) {
    throw new PostServerError(
      response.status,
      // 优先使用后端返回的message，无则使用默认提示
      body?.message ?? `Request failed: ${response.status}`,
    )
  }

  // 响应体为空 或者 业务数据data为null，抛出异常
  if (!body || body.data === null) {
    throw new PostServerError(
      response.status,
      body?.message ?? 'Response data is empty',
    )
  }

  // 校验全部通过，返回后端真实业务数据
  return body.data
}

/**
 * 根据前端传入的列表筛选参数，拼接URL查询字符串
 * @param req 帖子列表筛选条件结构体
 * @returns 拼接好query参数的接口路径 /post/list?page=1&keyword=xxx
 */
function createPostListPath(req: PostListRequest): string {
  // 内置URL参数工具，安全拼接query参数，自动处理编码
  const params = new URLSearchParams()

  // 分页页码，有值则加入参数
  if (req.page !== undefined)
    params.set('page', String(req.page))
  // 单页条数（后端字段蛇形page_size）
  if (req.pageSize !== undefined)
    params.set('page_size', String(req.pageSize))
  // 搜索关键词
  if (req.keyword)
    params.set('keyword', req.keyword)
  // 帖子类型
  if (req.type)
    params.set('type', req.type)
  // 分类ID
  if (req.categoryId !== undefined)
    params.set('category_id', String(req.categoryId))
  if (req.labelId !== undefined)
    params.set('label_id', String(req.labelId))
  // 作者ID
  if (req.authorId !== undefined)
    params.set('author_id', String(req.authorId))
  // 帖子状态
  if (req.status)
    params.set('status', req.status)
  // 排序规则
  if (req.sort)
    params.set('sort', req.sort)

  // 拼接查询参数字符串
  const query = params.toString()

  // 有筛选参数拼接完整url，无参数直接返回基础接口路径
  return query
    ? `/post/list?${query}`
    : '/post/list'
}

/**
 * 根据slug别名或ID获取单篇帖子详情
 * @param slugOrID 帖子数字ID / 友好链接slug
 * @returns 单篇帖子完整数据
 */
export const getPost = cache(async (
  slugOrID: string,
): Promise<Post> => {
  // 对传入的标识做URL编码，防止特殊符号破坏路由
  return serverGet<Post>(
    `/post/p/${encodeURIComponent(slugOrID)}`,
  )
})

/**
 * 获取帖子分页列表接口
 * @param req 筛选条件，默认空对象{}，不传参数时查询全部
 * @returns 分页帖子列表+总数等分页信息
 */
export async function listPosts(
  req: PostListRequest = {},
): Promise<PostListResponse> {
  // 调用参数拼接函数生成接口路径，再走通用GET请求
  return serverGet<PostListResponse>(
    createPostListPath(req),
  )
}

/**
 * 查询全部分类列表
 * @returns 分类数组
 */
export async function listCategories(): Promise<Category[]> {
  return serverGet<Category[]>('/post/categories')
}

export async function listLabels(): Promise<Label[]> {
  return serverGet<Label[]>('/post/labels')
}

/**
 * 获取随机一篇帖子
 * @returns 随机帖子详情
 */
export async function getRandomPost(): Promise<Post> {
  return serverGet<Post>('/post/random')
}
