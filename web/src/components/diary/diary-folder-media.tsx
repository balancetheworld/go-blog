import type { DiaryFolderMedia as DiaryFolderMediaValue } from '@/lib/diary-folders'
import Image from 'next/image'

interface DiaryFolderMediaProps {
  media: DiaryFolderMediaValue | null
}

export function DiaryFolderMedia({ media }: DiaryFolderMediaProps) {
  if (!media)
    return null

  return (
    <span className="folder-media" aria-hidden="true">
      {media.darkType === 'video'
        ? (
            <video
              className="folder-media-source folder-media-dark"
              src={media.dark}
              autoPlay
              muted
              loop
              playsInline
              preload="metadata"
            />
          )
        : (
            <Image
              className="folder-media-source folder-media-dark"
              src={media.dark}
              alt=""
              fill
              sizes="520px"
              unoptimized
            />
          )}
      {media.lightType === 'video'
        ? (
            <video
              className="folder-media-source folder-media-light"
              src={media.light}
              autoPlay
              muted
              loop
              playsInline
              preload="metadata"
            />
          )
        : (
            <Image
              className="folder-media-source folder-media-light"
              src={media.light}
              alt=""
              fill
              sizes="520px"
              unoptimized
            />
          )}
      <span className="folder-media-shade" />
    </span>
  )
}
