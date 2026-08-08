import { load } from 'cheerio'

const allowedTags = new Set([
  'a',
  'blockquote',
  'br',
  'code',
  'em',
  'h1',
  'h2',
  'h3',
  'h4',
  'hr',
  'img',
  'li',
  'ol',
  'p',
  'pre',
  's',
  'strong',
  'ul',
])

const allowedAttributes: Record<string, Set<string>> = {
  a: new Set(['href', 'rel', 'title']),
  img: new Set(['alt', 'src', 'title']),
}

function isSafeURL(value: string, allowMailto: boolean): boolean {
  if (value.startsWith('/') || value.startsWith('#'))
    return true

  try {
    const url = new URL(value)
    return url.protocol === 'http:'
      || url.protocol === 'https:'
      || (allowMailto && url.protocol === 'mailto:')
  }
  catch {
    return false
  }
}

export function sanitizePostHTML(value: string): string {
  const $ = load(value, null, false)
  $('script, style, iframe, object, embed, form, input, button, textarea, select, link, meta').remove()

  $('*').each((_, element) => {
    if (!('tagName' in element) || !('attribs' in element))
      return

    const tagName = element.tagName.toLowerCase()
    if (!allowedTags.has(tagName)) {
      $(element).replaceWith($(element).contents())
      return
    }

    const tagAttributes = allowedAttributes[tagName] ?? new Set<string>()
    for (const attribute of Object.keys(element.attribs)) {
      if (!tagAttributes.has(attribute))
        $(element).removeAttr(attribute)
    }

    if (tagName === 'a') {
      const href = $(element).attr('href')
      if (href && !isSafeURL(href, true))
        $(element).removeAttr('href')
      if ($(element).attr('href'))
        $(element).attr('rel', 'nofollow noreferrer')
    }

    if (tagName === 'img') {
      const src = $(element).attr('src')
      if (src && !isSafeURL(src, false))
        $(element).removeAttr('src')
    }
  })

  return $.html()
}
