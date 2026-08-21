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

  // 用于资源加载的 useEffect
  useEffect(() => {
    let disposed = false
    let idleID: number | null = null
    let fallbackTimer: ReturnType<typeof setTimeout> | null = null

    // 初始化拿到 rive 的函数，异步导入，最后绑定给 riveRef
    function initialize() {
      void import('@rive-app/canvas').then(({ Rive: RiveConstructor }) => {
        const canvas = canvasRef.current
        if (disposed || !canvas)
          return

        const rive: Rive = new RiveConstructor({
          src: '/design-assets/25691-49048-interactive-icon-set.riv',
          canvas,
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
    }
    // 计划导入函数。 把 `initialize` 初始化函数，安排到浏览器主线程空闲的时候再执行。
    function scheduleInitialize() { // 浏览器原生 API：把低优先级任务，放到浏览器主线程空闲间隙去执行，不抢占动画、点击、渲染的时间，用来做性能优化。
      if ('requestIdleCallback' in window) {
        idleID = window.requestIdleCallback(initialize, { // `requestIdleCallback()` 返回的任务 ID。 作用：保存浏览器空闲回调句柄，后续可以用 `cancelIdleCallback(idleID)` 取消任务。
          timeout: 2000,
        })
        return
      }
      fallbackTimer = setTimeout(initialize, 200) // 降级方案
    }
    // 如果已经加载完了，就直接执行
    if (document.readyState === 'complete') {
      scheduleInitialize()
    }
    else {
      window.addEventListener('load', scheduleInitialize, { once: true })
    }
    const handleResize = () => riveRef.current?.resizeDrawingSurfaceToCanvas()
    window.addEventListener('resize', handleResize, { passive: true })

    return () => {
      disposed = true
      window.removeEventListener('load', scheduleInitialize)
      window.removeEventListener('resize', handleResize)

      if (idleID !== null && 'cancelIdleCallback' in window)
        window.cancelIdleCallback(idleID)

      if (fallbackTimer)
        clearTimeout(fallbackTimer)

      if (timerRef.current)
        clearTimeout(timerRef.current)

      riveRef.current?.cleanup()
      riveRef.current = null
      loadedRef.current = false
    }
  }, [])

  // 用于动画执行的 useEffect 外部playkey 变化驱动动画加载
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
