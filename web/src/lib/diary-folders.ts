export const diaryFolderSlugs = [
  'travel',
  'daily',
  'inspiration',
  'notes',
  'favorites',
] as const

export type DiaryFolderSlug = typeof diaryFolderSlugs[number]

export interface DiaryFolderMedia {
  dark: string
  light: string
  darkType: 'image' | 'video'
  lightType: 'image' | 'video'
}

export const allDiaryFolderMedia: DiaryFolderMedia = {
  dark: '/design-assets/d-all.jpg',
  light: '/design-assets/l-all.jpg',
  darkType: 'image',
  lightType: 'image',
}

const diaryFolderStyles: Record<DiaryFolderSlug, string> = {
  travel: 'diary-folder-travel',
  daily: 'diary-folder-daily',
  inspiration: 'diary-folder-ideas',
  notes: 'diary-folder-memory',
  favorites: 'diary-folder-reading',
}

const diaryFolderMedia: Record<DiaryFolderSlug, DiaryFolderMedia> = {
  travel: {
    dark: '/design-assets/d-1.jpg',
    light: '/design-assets/l-1.jpg',
    darkType: 'image',
    lightType: 'image',
  },
  daily: {
    dark: '/design-assets/d-2.mp4',
    light: '/design-assets/l-2.jpg',
    darkType: 'video',
    lightType: 'image',
  },
  inspiration: {
    dark: '/design-assets/d-3.jpg',
    light: '/design-assets/l-3.jpg',
    darkType: 'image',
    lightType: 'image',
  },
  notes: {
    dark: '/design-assets/d-4.jpg',
    light: '/design-assets/l-4.jpg',
    darkType: 'image',
    lightType: 'image',
  },
  favorites: {
    dark: '/design-assets/d-5.jpg',
    light: '/design-assets/l-5.jpg',
    darkType: 'image',
    lightType: 'image',
  },
}

export function isDiaryFolderSlug(slug: string): slug is DiaryFolderSlug {
  return diaryFolderSlugs.includes(slug as DiaryFolderSlug)
}

export function getDiaryFolderStyle(slug: string): string {
  return isDiaryFolderSlug(slug)
    ? diaryFolderStyles[slug]
    : 'diary-folder-daily'
}

export function getDiaryFolderMedia(slug: string): DiaryFolderMedia | null {
  return isDiaryFolderSlug(slug) ? diaryFolderMedia[slug] : null
}
