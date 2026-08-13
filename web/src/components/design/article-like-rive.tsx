'use client'

import type { Rive } from '@rive-app/canvas'
import { useEffect, useRef } from 'react'

interface ArticleLikeRiveProps {
  playKey: number
}

const animationName = '16_Bear_Anim'

export function ArticleLikeRive({ playKey }: ArticleLikeRiveProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const riveRef = useRef<import('@rive-app/canvas').Rive | null>(null)
  const loadedRef = useRef(false)
  const pendingPlayRef = useRef(false)
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => {
    let disposed = false

    void import('@rive-app/canvas').then(({ Rive: RiveConstructor }) => {
      if (disposed || !canvasRef.current)
        return

      const rive: Rive = new RiveConstructor({
        src: '/design-assets/25691-49048-interactive-icon-set.riv',
        canvas: canvasRef.current,
        artboard: '16_Bear',
        animations: animationName,
        autoplay: false,
        onLoad: () => {
          loadedRef.current = true
          rive.resizeDrawingSurfaceToCanvas()
          if (pendingPlayRef.current) {
            pendingPlayRef.current = false
            rive.play(animationName)
          }
        },
        onLoop: () => rive.stop(animationName),
      })

      riveRef.current = rive
    })

    const handleResize = () => riveRef.current?.resizeDrawingSurfaceToCanvas()
    window.addEventListener('resize', handleResize, { passive: true })

    return () => {
      disposed = true
      window.removeEventListener('resize', handleResize)
      if (timerRef.current)
        clearTimeout(timerRef.current)
      riveRef.current?.cleanup()
      riveRef.current = null
      loadedRef.current = false
    }
  }, [])

  useEffect(() => {
    if (playKey === 0)
      return

    const rive = riveRef.current
    if (!rive || !loadedRef.current) {
      pendingPlayRef.current = true
      return
    }

    rive.stop(animationName)
    rive.play(animationName)
    if (timerRef.current)
      clearTimeout(timerRef.current)
    timerRef.current = setTimeout(() => rive.stop(animationName), 1400)
  }, [playKey])

  return (
    <canvas
      ref={canvasRef}
      className="rive-like-canvas"
      width="160"
      height="160"
      aria-hidden="true"
    />
  )
}
