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
  'mark',
  'ol',
  'p',
  'pre',
  's',
  'span',
  'strong',
  'ul',
])

const allowedAttributes: Record<string, Set<string>> = {
  a: new Set(['href', 'rel', 'title']),
  code: new Set(['class']),
  img: new Set(['alt', 'data-align', 'data-width', 'src', 'title', 'width']),
  mark: new Set(['data-color', 'style']),
  span: new Set(['class', 'style']),
}

const colorPattern = /^(?:#[0-9a-f]{3,8}|(?:rgb|hsl)a?\([\d.%\s,]+\)|[a-z]+)$/i

function safeStyle(tagName: string, value: string): string | null {
  const declarations = value
    .split(';')
    .map(item => item.trim())
    .filter(Boolean)

  if (tagName === 'span' && declarations.length === 1) {
    const separator = declarations[0].indexOf(':')
    const property = declarations[0].slice(0, separator).trim().toLowerCase()
    const color = declarations[0].slice(separator + 1).trim()
    if (separator > 0 && property === 'color' && colorPattern.test(color))
      return `color: ${color}`
  }

  if (tagName === 'mark') {
    const background = declarations.find((item) => {
      const separator = item.indexOf(':')
      return separator > 0
        && item.slice(0, separator).trim().toLowerCase() === 'background-color'
    })
    const separator = background?.indexOf(':') ?? -1
    const color = separator > 0
      ? background?.slice(separator + 1).trim()
      : undefined
    if (color && colorPattern.test(color))
      return `background-color: ${color}; color: inherit`
  }

  if (['h1', 'h2', 'h3', 'h4', 'p'].includes(tagName) && declarations.length === 1) {
    const separator = declarations[0].indexOf(':')
    const property = declarations[0].slice(0, separator).trim().toLowerCase()
    const align = declarations[0].slice(separator + 1).trim().toLowerCase()
    if (
      separator > 0
      && property === 'text-align'
      && ['left', 'center', 'right', 'justify'].includes(align)
    ) {
      return `text-align: ${align}`
    }
  }

  return null
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

    const tagAttributes = allowedAttributes[tagName] ?? (
      ['h1', 'h2', 'h3', 'h4', 'p'].includes(tagName)
        ? new Set(['style'])
        : new Set<string>()
    )
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

      const align = $(element).attr('data-align')
      if (align && !['left', 'center', 'right'].includes(align))
        $(element).removeAttr('data-align')

      const width = $(element).attr('data-width') ?? $(element).attr('width')
      if (width && !/^\d{2,4}$/.test(width)) {
        $(element).removeAttr('data-width')
        $(element).removeAttr('width')
      }
    }

    if (tagName === 'code' || tagName === 'span') {
      const className = $(element).attr('class')
      if (
        className
        && !className.split(/\s+/).every((value) => {
          if (tagName === 'code')
            return /^language-[a-z0-9-]+$/.test(value)

          return /^hljs-[a-z0-9-]+$/.test(value)
        })
      ) {
        $(element).removeAttr('class')
      }
    }

    const style = $(element).attr('style')
    if (style) {
      const sanitizedStyle = safeStyle(tagName, style)
      if (sanitizedStyle)
        $(element).attr('style', sanitizedStyle)
      else
        $(element).removeAttr('style')
    }

    if (tagName === 'mark') {
      const color = $(element).attr('data-color')
      if (color && !colorPattern.test(color))
        $(element).removeAttr('data-color')
    }
  })

  return $.html()
}
