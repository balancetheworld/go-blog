'use client'

import { useEffect, useRef } from 'react'

interface Particle {
  x: number
  y: number
  size: number
  baseSpeedX: number
  baseSpeedY: number
  speedX: number
  speedY: number
  baseOpacity: number
  opacity: number
  excitement: number
}

export function ParticleCanvas() {
  const canvasRef = useRef<HTMLCanvasElement>(null)

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas)
      return

    const context = canvas.getContext('2d')
    if (!context)
      return

    const canvasElement = canvas
    const drawingContext = context
    const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    if (reducedMotion)
      return

    let frame = 0
    let width = 0
    let height = 0
    let particles: Particle[] = []
    const mouse = { x: -1000, y: -1000 }
    let mouseThrottled = false

    function resize() {
      const ratio = Math.min(window.devicePixelRatio || 1, 2)
      width = window.innerWidth
      height = window.innerHeight
      canvasElement.width = width * ratio
      canvasElement.height = height * ratio
      drawingContext.setTransform(ratio, 0, 0, ratio, 0, 0)
      particles = Array.from(
        { length: width < 768 ? 50 : 120 },
        () => {
          const speedX = (Math.random() - 0.5) * 0.5
          const speedY = (Math.random() - 0.5) * 0.5
          const opacity = Math.random() * 0.5 + 0.1
          return {
            x: Math.random() * width,
            y: Math.random() * height,
            size: Math.random() * 2 + 0.5,
            baseSpeedX: speedX,
            baseSpeedY: speedY,
            speedX,
            speedY,
            baseOpacity: opacity,
            opacity,
            excitement: 0,
          }
        },
      )
    }

    function draw() {
      drawingContext.clearRect(0, 0, width, height)

      if (document.documentElement.dataset.theme === 'cat') {
        frame = window.requestAnimationFrame(draw)
        return
      }

      for (const particle of particles) {
        if (particle.excitement > 0) {
          particle.excitement *= 0.92
          if (particle.excitement < 0.01)
            particle.excitement = 0
        }

        const deltaX = particle.x - mouse.x
        const deltaY = particle.y - mouse.y
        const distance = Math.hypot(deltaX, deltaY)

        if (distance < 150 && distance > 0) {
          const force = (150 - distance) / 150
          const push = Math.min(force * 3, 2)
          particle.x += deltaX / distance * push
          particle.y += deltaY / distance * push
          particle.excitement = Math.min(1, particle.excitement + force * 0.3)
        }

        const speedBoost = 1 + particle.excitement * 2
        particle.x += particle.speedX * speedBoost
        particle.y += particle.speedY * speedBoost
        particle.speedX += (particle.baseSpeedX - particle.speedX) * 0.05
        particle.speedY += (particle.baseSpeedY - particle.speedY) * 0.05

        if (particle.x < -50)
          particle.x = width + 50
        if (particle.x > width + 50)
          particle.x = -50
        if (particle.y < -50)
          particle.y = height + 50
        if (particle.y > height + 50)
          particle.y = -50

        const targetOpacity = particle.baseOpacity + particle.excitement * 0.6
        particle.opacity += (targetOpacity - particle.opacity) * 0.15

        drawingContext.beginPath()
        drawingContext.arc(particle.x, particle.y, particle.size, 0, Math.PI * 2)
        drawingContext.fillStyle = `rgba(0, 212, 255, ${particle.opacity})`
        drawingContext.fill()
      }

      const averageExcitement = particles.reduce((sum, particle) => sum + particle.excitement, 0) / particles.length
      const maxDistance = 100 * (1 - averageExcitement * 0.3)

      for (let leftIndex = 0; leftIndex < particles.length; leftIndex++) {
        for (let rightIndex = leftIndex + 1; rightIndex < particles.length; rightIndex++) {
          const left = particles[leftIndex]
          const right = particles[rightIndex]
          const distance = Math.hypot(left.x - right.x, left.y - right.y)
          if (distance >= maxDistance)
            continue

          const opacity = (maxDistance - distance) / maxDistance * 0.15
          drawingContext.beginPath()
          drawingContext.strokeStyle = `rgba(0, 212, 255, ${opacity})`
          drawingContext.lineWidth = 0.5
          drawingContext.moveTo(left.x, left.y)
          drawingContext.lineTo(right.x, right.y)
          drawingContext.stroke()
        }
      }

      frame = window.requestAnimationFrame(draw)
    }

    function handleMouseMove(event: MouseEvent) {
      if (mouseThrottled)
        return
      mouse.x = event.clientX
      mouse.y = event.clientY
      mouseThrottled = true
      window.requestAnimationFrame(() => {
        mouseThrottled = false
      })
    }

    function handleMouseLeave() {
      mouse.x = -1000
      mouse.y = -1000
    }

    function handleTouchMove(event: TouchEvent) {
      const touch = event.touches[0]
      if (!touch || mouseThrottled)
        return
      mouse.x = touch.clientX
      mouse.y = touch.clientY
      mouseThrottled = true
      window.requestAnimationFrame(() => {
        mouseThrottled = false
      })
    }

    function handleBurst() {
      for (const particle of particles) {
        const angle = Math.random() * Math.PI * 2
        const jump = Math.random() * 30 + 10
        particle.x += Math.cos(angle) * jump
        particle.y += Math.sin(angle) * jump
        particle.excitement = 1
      }
    }

    resize()
    draw()
    window.addEventListener('resize', resize)
    window.addEventListener('mousemove', handleMouseMove, { passive: true })
    document.addEventListener('mouseleave', handleMouseLeave)
    window.addEventListener('touchmove', handleTouchMove, { passive: true })
    window.addEventListener('touchend', handleMouseLeave, { passive: true })
    window.addEventListener('design-theme-burst', handleBurst)

    return () => {
      window.cancelAnimationFrame(frame)
      window.removeEventListener('resize', resize)
      window.removeEventListener('mousemove', handleMouseMove)
      document.removeEventListener('mouseleave', handleMouseLeave)
      window.removeEventListener('touchmove', handleTouchMove)
      window.removeEventListener('touchend', handleMouseLeave)
      window.removeEventListener('design-theme-burst', handleBurst)
    }
  }, [])

  return <canvas ref={canvasRef} id="particle-canvas" aria-hidden="true" />
}
