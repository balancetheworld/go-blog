'use client'

import { gsap } from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import { usePathname } from 'next/navigation'
import { useLayoutEffect } from 'react'

export function DesignAnimations() {
  const pathname = usePathname()

  useLayoutEffect(() => {
    const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    if (reducedMotion)
      return

    gsap.registerPlugin(ScrollTrigger)
    const cleanups: Array<() => void> = []
    const context = gsap.context(() => {
      gsap.utils.toArray<HTMLElement>('.article-card').forEach((card) => {
        gsap.from(card, {
          opacity: 0,
          x: -50,
          duration: 0.8,
          ease: 'expo.out',
          scrollTrigger: {
            trigger: card,
            start: 'top 85%',
            toggleActions: 'play none none reverse',
          },
        })
      })

      gsap.utils.toArray<HTMLElement>('.timeline-item').forEach((item, index) => {
        gsap.from(item, {
          opacity: 0,
          y: 28,
          duration: 0.7,
          delay: index * 0.04,
          ease: 'expo.out',
          scrollTrigger: {
            trigger: item,
            start: 'top 88%',
            toggleActions: 'play none none reverse',
          },
        })
      })

      gsap.utils.toArray<HTMLElement>('.moment-card').forEach((card) => {
        gsap.from(card, {
          opacity: 0,
          y: 30,
          duration: 0.7,
          ease: 'expo.out',
          scrollTrigger: {
            trigger: card,
            start: 'top 88%',
            toggleActions: 'play none none reverse',
          },
        })
      })

      gsap.utils.toArray<HTMLElement>('.section-header').forEach((header) => {
        gsap.from(header, {
          opacity: 0,
          y: 40,
          duration: 0.8,
          ease: 'expo.out',
          scrollTrigger: {
            trigger: header,
            start: 'top 85%',
            toggleActions: 'play none none reverse',
          },
        })
      })

      gsap.utils.toArray<HTMLElement>('.magnetic').forEach((button) => {
        function handleMove(event: MouseEvent) {
          const rect = button.getBoundingClientRect()
          const x = event.clientX - rect.left - rect.width / 2
          const y = event.clientY - rect.top - rect.height / 2
          gsap.to(button, {
            x: x * 0.3,
            y: y * 0.3,
            duration: 0.4,
            ease: 'power2.out',
          })
        }

        function handleLeave() {
          gsap.to(button, {
            x: 0,
            y: 0,
            duration: 0.6,
            ease: 'elastic.out(1, 0.4)',
          })
        }

        button.addEventListener('mousemove', handleMove)
        button.addEventListener('mouseleave', handleLeave)
        cleanups.push(() => {
          button.removeEventListener('mousemove', handleMove)
          button.removeEventListener('mouseleave', handleLeave)
        })
      })

      gsap.utils.toArray<HTMLElement>('[data-tilt]').forEach((card) => {
        function handleMove(event: MouseEvent) {
          const rect = card.getBoundingClientRect()
          const x = (event.clientX - rect.left) / rect.width
          const y = (event.clientY - rect.top) / rect.height
          gsap.to(card, {
            rotateX: (y - 0.5) * -10,
            rotateY: (x - 0.5) * 10,
            transformPerspective: 1000,
            duration: 0.5,
            ease: 'power2.out',
          })
        }

        function handleLeave() {
          gsap.to(card, {
            rotateX: 0,
            rotateY: 0,
            duration: 0.8,
            ease: 'elastic.out(1, 0.4)',
          })
        }

        card.addEventListener('mousemove', handleMove)
        card.addEventListener('mouseleave', handleLeave)
        cleanups.push(() => {
          card.removeEventListener('mousemove', handleMove)
          card.removeEventListener('mouseleave', handleLeave)
        })
      })
    })

    window.requestAnimationFrame(() => ScrollTrigger.refresh())

    return () => {
      cleanups.forEach(cleanup => cleanup())
      context.revert()
    }
  }, [pathname])

  return null
}
