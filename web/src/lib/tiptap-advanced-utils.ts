const entities: Record<string, string> = {
  '&amp;': '&',
  '&apos;': '\'',
  '&gt;': '>',
  '&lt;': '<',
  '&nbsp;': ' ',
  '&quot;': '"',
}

export function htmlToPlainText(html: string): string {
  return html
    .replace(/<br\s*\/?>/gi, ' ')
    .replace(/<\/(?:div|h[1-6]|li|p|pre|blockquote)>/gi, ' ')
    .replace(/<[^>]*>/g, '')
    .replace(/&(amp|apos|gt|lt|nbsp|quot);/g, value => entities[value] ?? value)
    .replace(/\s+/g, ' ')
    .trim()
}

export function isHTMLContentEmpty(html: string): boolean {
  return !/<img\b/i.test(html) && htmlToPlainText(html).length === 0
}

export function truncateText(value: string, maxLength: number): string {
  return value.length > maxLength
    ? `${value.slice(0, maxLength).trimEnd()}...`
    : value
}

export function isSafeImageURL(value: string): boolean {
  const url = value.trim()

  if (url.startsWith('/') && !url.startsWith('//'))
    return true

  try {
    const parsed = new URL(url)
    return parsed.protocol === 'http:' || parsed.protocol === 'https:'
  }
  catch {
    return false
  }
}
