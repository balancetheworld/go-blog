import { getRequestConfig } from 'next-intl/server'
import { cookies } from 'next/headers'
import enMessages from '../../messages/en.json'
import zhCNMessages from '../../messages/zh-CN.json'
import { defaultLocale, isLocale, localeCookieName } from './config'

export default getRequestConfig(async () => {
  const cookieStore = await cookies()
  const cookieLocale = cookieStore.get(localeCookieName)?.value
  const locale = isLocale(cookieLocale) ? cookieLocale : defaultLocale

  return {
    locale,
    messages: locale === 'en' ? enMessages : zhCNMessages,
    timeZone: 'Asia/Tokyo',
  }
})
