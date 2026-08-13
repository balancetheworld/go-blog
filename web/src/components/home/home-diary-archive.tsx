'use client'

import type { Diary, DiaryFolder } from '@/models/diary'
import { gsap } from 'gsap'
import { useTranslations } from 'next-intl'
import Link from 'next/link'
import { useLayoutEffect, useRef, useState } from 'react'
import { DiaryList } from '@/components/diary'
import { DiaryFolderMedia } from '@/components/diary/diary-folder-media'
import { allDiaryFolderMedia, getDiaryFolderMedia, getDiaryFolderStyle, isDiaryFolderSlug } from '@/lib/diary-folders'

interface HomeDiaryArchiveProps {
  diaries: Diary[]
  folders: DiaryFolder[]
}

export function HomeDiaryArchive({ diaries, folders }: HomeDiaryArchiveProps) {
  const t = useTranslations('Home.diaries')
  const diaryT = useTranslations('Diary')
  const archiveRef = useRef<HTMLDivElement>(null)
  const folderViewRef = useRef<HTMLDivElement>(null)
  const listViewRef = useRef<HTMLDivElement>(null)
  const [selectedFolderID, setSelectedFolderID] = useState<number | null>(null)
  const [listOpen, setListOpen] = useState(false)

  const visibleDiaries = selectedFolderID === null
    ? diaries
    : diaries.filter(diary => diary.folder?.id === selectedFolderID)
  const selectedFolder = folders.find(folder => folder.id === selectedFolderID)

  useLayoutEffect(() => {
    if (!listOpen || !listViewRef.current)
      return

    const cards = listViewRef.current.querySelectorAll('.diary-card')
    const floats = archiveRef.current?.querySelectorAll('.diary-float, .diary-orbit')
    const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    if (reducedMotion)
      return

    gsap.fromTo(cards, {
      opacity: 0,
      y: 34,
      rotateX: -12,
      rotateY: (index: number) => index % 2 ? 8 : -8,
      scale: 0.94,
    }, {
      opacity: 1,
      y: 0,
      rotateX: 0,
      rotateY: 0,
      scale: 1,
      duration: 0.58,
      stagger: 0.045,
      ease: 'expo.out',
    })

    if (floats?.length) {
      gsap.fromTo(floats, {
        opacity: 0.22,
        filter: 'blur(2px)',
      }, {
        opacity: 0.58,
        filter: 'blur(0px)',
        duration: 0.72,
        stagger: 0.05,
        ease: 'elastic.out(1, 0.55)',
      })
    }
  }, [listOpen, selectedFolderID])

  function setActiveFolder(folder: HTMLElement) {
    const folderView = folderViewRef.current
    if (!folderView)
      return
    folderView.classList.add('has-active')
    folderView.querySelectorAll('.diary-folder').forEach((item) => {
      item.classList.toggle('is-active', item === folder)
    })
  }

  function clearActiveFolders() {
    const folderView = folderViewRef.current
    if (!folderView || listOpen)
      return
    folderView.classList.remove('has-active')
    folderView.querySelectorAll('.diary-folder').forEach(item => item.classList.remove('is-active'))
  }

  function openFolder(folderID: number | null, folder: HTMLElement) {
    const archive = archiveRef.current
    const folderView = folderViewRef.current
    if (!archive || !folderView)
      return

    const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    const allFolders = Array.from(folderView.querySelectorAll<HTMLElement>('.diary-folder'))
    const otherFolders = allFolders.filter(item => item !== folder)
    const papers = folder.querySelectorAll('.folder-paper, .folder-front')

    setActiveFolder(folder)
    folder.classList.add('is-opening')

    if (reducedMotion) {
      setSelectedFolderID(folderID)
      setListOpen(true)
      folder.classList.remove('is-opening')
      return
    }

    gsap.timeline({
      defaults: { ease: 'expo.out' },
      onComplete: () => {
        setSelectedFolderID(folderID)
        setListOpen(true)
        folder.classList.remove('is-opening')
        gsap.set([folder, folderView, ...Array.from(papers), ...otherFolders], { clearProps: 'all' })
      },
    })
      .to(otherFolders, {
        opacity: 0.28,
        y: 18,
        scale: 0.94,
        duration: 0.28,
        stagger: 0.025,
      }, 0)
      .to(folder, {
        y: -12,
        scale: 1.04,
        duration: 0.32,
      }, 0)
      .to(papers, {
        x: 0,
        y: 0,
        rotate: 0,
        rotateY: 0,
        scale: 1,
        opacity: 1,
        duration: 0.46,
        stagger: 0.05,
      }, 0.12)
      .to(folderView, {
        opacity: 0,
        y: -24,
        duration: 0.28,
      }, 0.48)
  }

  function closeFolder() {
    const archive = archiveRef.current
    const folderView = folderViewRef.current
    const listView = listViewRef.current
    if (!archive || !folderView || !listView || !listOpen)
      return

    const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    const cards = Array.from(listView.querySelectorAll<HTMLElement>('.diary-card'))
    const foldersElements = Array.from(folderView.querySelectorAll<HTMLElement>('.diary-folder'))

    if (reducedMotion) {
      setListOpen(false)
      setSelectedFolderID(null)
      clearActiveFolders()
      return
    }

    gsap.timeline({
      defaults: { ease: 'power2.inOut' },
      onComplete: () => {
        setListOpen(false)
        setSelectedFolderID(null)
        folderView.classList.remove('has-active')
        foldersElements.forEach(folder => folder.classList.remove('is-active'))
        gsap.set(folderView, { opacity: 0, y: 18 })
        gsap.set([listView, ...cards, ...foldersElements], { clearProps: 'all' })
        gsap.to(folderView, {
          opacity: 1,
          y: 0,
          duration: 0.42,
          ease: 'expo.out',
          onComplete: () => {
            gsap.set(folderView, { clearProps: 'all' })
          },
        })
      },
    })
      .to(cards, {
        opacity: 0,
        y: 34,
        scale: 0.96,
        duration: 0.28,
        stagger: 0.025,
      }, 0)
      .to(listView, {
        opacity: 0,
        duration: 0.24,
      }, 0.08)
  }

  function handlePointerMove(event: React.PointerEvent<HTMLDivElement>) {
    const archive = archiveRef.current
    if (!archive || !window.matchMedia('(hover: hover) and (pointer: fine)').matches)
      return
    const rect = archive.getBoundingClientRect()
    const x = ((event.clientX - rect.left) / rect.width - 0.5) * 2
    const y = ((event.clientY - rect.top) / rect.height - 0.5) * 2
    archive.style.setProperty('--diary-tilt-x', `${x * 3}deg`)
    archive.style.setProperty('--diary-tilt-y', `${y * -2}deg`)
  }

  function handlePointerLeave() {
    const archive = archiveRef.current
    if (!archive)
      return
    archive.style.setProperty('--diary-tilt-x', '0deg')
    archive.style.setProperty('--diary-tilt-y', '0deg')
    clearActiveFolders()
  }

  const folderItems = [
    { id: null, name: t('all'), description: t('count', { count: diaries.length }), style: 'diary-folder-all', media: allDiaryFolderMedia },
    ...folders.map(folder => ({
      id: folder.id,
      name: isDiaryFolderSlug(folder.slug) ? diaryT(`folders.${folder.slug}`) : folder.name,
      description: folder.description || t('view'),
      style: getDiaryFolderStyle(folder.slug),
      media: getDiaryFolderMedia(folder.slug),
    })),
  ]

  return (
    <div
      ref={archiveRef}
      className={`diary-archive${listOpen ? ' is-list-view' : ''}`}
      onPointerMove={handlePointerMove}
      onPointerLeave={handlePointerLeave}
    >
      <div className="diary-topline">
        <div className="diary-months">
          <span className="diary-month active">{t('recent')}</span>
        </div>
        <Link href="/diary" className="diary-back">{t('viewAll')}</Link>
      </div>
      <div ref={folderViewRef} className="diary-folder-view">
        {folderItems.map(item => (
          <button
            key={item.id ?? 'all'}
            className={`diary-folder ${item.style}`}
            type="button"
            onPointerEnter={event => setActiveFolder(event.currentTarget)}
            onFocus={event => setActiveFolder(event.currentTarget)}
            onClick={event => openFolder(item.id, event.currentTarget)}
          >
            <span className="folder-paper paper-one" />
            <span className="folder-paper paper-two" />
            <span className="folder-back" />
            <span className="folder-front">
              <DiaryFolderMedia media={item.media} />
            </span>
            <span className="folder-tab">{item.name}</span>
            <span className="folder-meta">{item.description}</span>
          </button>
        ))}
      </div>
      <div className="diary-atmosphere" aria-hidden="true">
        <span className="diary-float diary-float-photo" />
        <span className="diary-float diary-float-ticket" />
        <span className="diary-float diary-float-note" />
        <span className="diary-float diary-float-dot" />
        <span className="diary-orbit">
          <i />
          <i />
          <i />
        </span>
      </div>
      <div ref={listViewRef} className="diary-list-view">
        <div className="diary-list-head">
          <button className="diary-back" type="button" onClick={closeFolder}>{t('back')}</button>
          <h3 className="diary-list-title">
            {selectedFolder && isDiaryFolderSlug(selectedFolder.slug)
              ? diaryT(`folders.${selectedFolder.slug}`)
              : selectedFolder?.name || t('allDiaries')}
          </h3>
        </div>
        <DiaryList diaries={visibleDiaries} />
      </div>
    </div>
  )
}
