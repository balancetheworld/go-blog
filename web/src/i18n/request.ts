import { getRequestConfig } from 'next-intl/server'
import messages from '../../messages/zh-CN.json'

export default getRequestConfig(() => ({
  locale: 'zh-CN',
  messages,
}))
