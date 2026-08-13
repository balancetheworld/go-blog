'use client'

import type { CSSProperties } from 'react'
import { gsap } from 'gsap'
import { useEffect, useRef } from 'react'

export function DesignBackground() {
  const orbitRef = useRef<HTMLDivElement>(null)
  const trackRef = useRef<HTMLDivElement>(null)
  const fillRef = useRef<HTMLSpanElement>(null)
  const thumbRef = useRef<HTMLSpanElement>(null)

  useEffect(() => {
    const reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    let lastProgress = 0
    let rotation = 0

    function update() {
      const scroll = window.scrollY
      const maxScroll = Math.max(1, document.documentElement.scrollHeight - window.innerHeight)
      const progress = Math.min(1, scroll / maxScroll)
      const trackHeight = trackRef.current?.clientHeight ?? 0
      const thumbHeight = thumbRef.current?.offsetHeight ?? 40
      const thumbY = progress * Math.max(trackHeight - thumbHeight, 0)
      const scrollDelta = progress - lastProgress
      const rotationDelta = Math.max(-90, Math.min(90, scrollDelta * 1800))

      lastProgress = progress
      rotation += rotationDelta

      document.documentElement.style.setProperty('--cat-scroll', `${progress * 260}px`)
      document.documentElement.style.setProperty('--cat-progress', progress.toFixed(3))

      if (reducedMotion) {
        if (orbitRef.current)
          orbitRef.current.style.opacity = '1'
        if (fillRef.current)
          fillRef.current.style.height = `${progress * 100}%`
        if (thumbRef.current)
          thumbRef.current.style.transform = `translateY(${thumbY}px) rotate(${progress * 720}deg)`
        return
      }

      if (orbitRef.current) {
        gsap.to(orbitRef.current, {
          autoAlpha: scroll > 12 ? 1 : 0.72,
          duration: 0.42,
          ease: 'expo.out',
          overwrite: true,
        })
      }
      if (fillRef.current) {
        gsap.to(fillRef.current, {
          height: `${progress * 100}%`,
          duration: 0.52,
          ease: 'expo.out',
          overwrite: true,
        })
      }
      if (thumbRef.current) {
        gsap.to(thumbRef.current, {
          y: thumbY,
          rotation,
          duration: 0.58,
          ease: 'expo.out',
          overwrite: 'auto',
        })
      }
    }

    update()
    window.addEventListener('scroll', update, { passive: true })
    window.addEventListener('resize', update, { passive: true })
    window.addEventListener('load', update, { passive: true })

    return () => {
      window.removeEventListener('scroll', update)
      window.removeEventListener('resize', update)
      window.removeEventListener('load', update)
      gsap.killTweensOf([orbitRef.current, fillRef.current, thumbRef.current])
    }
  }, [])

  return (
    <>
      <div className="cat-parallax-scene" aria-hidden="true">
        <svg className="cat-parallax-pattern" viewBox="0 0 1440 1000" fill="none">
          <g className="cat-parallax-layer doodle doodle-cat-blue" style={{ '--depth': 0.28 } as CSSProperties}>
            <g transform="translate(14 124)">
              <path d="M164 618 L220 588 L254 616 L296 594 L332 638 L330 748 L294 784 L204 778 L158 724 Z" fill="#00D4E8" />
              <path d="M164 618 L220 588 L254 616 L296 594 L332 638 L330 748 L294 784 L204 778 L158 724 Z" stroke="#0A0A0F" strokeWidth="14" strokeLinejoin="round" />
              <path d="M229 675 L229 724 M283 675 L283 724" stroke="#0A0A0F" strokeWidth="12" strokeLinecap="round" />
              <path d="M232 744 C250 760 266 758 280 742" stroke="#0A0A0F" strokeWidth="11" strokeLinecap="round" />
              <path d="M185 678 L128 656 M185 708 L120 708 M186 738 L132 764 M327 678 L384 656 M327 708 L392 708 M326 738 L380 764" stroke="#0A0A0F" strokeWidth="10" strokeLinecap="round" />
            </g>
          </g>
          <g className="cat-parallax-layer doodle doodle-rabbit" style={{ '--depth': -0.42 } as CSSProperties}>
            <g transform="translate(70 -30)">
              <path d="M770 344 C742 268 750 174 792 166 C834 158 842 256 818 340 Z" fill="#FFFFFF" stroke="#0A0A0F" strokeWidth="14" strokeLinejoin="round" />
              <path d="M832 340 C854 254 902 172 936 194 C974 220 922 306 864 356 Z" fill="#FFFFFF" stroke="#0A0A0F" strokeWidth="14" strokeLinejoin="round" />
              <path d="M724 356 L816 318 L914 362 L926 454 L866 512 L752 498 L696 438 Z" fill="#FFFFFF" stroke="#0A0A0F" strokeWidth="14" strokeLinejoin="round" />
              <path d="M774 408 L774 408 M846 406 L846 406" stroke="#1B2560" strokeWidth="18" strokeLinecap="round" />
              <path d="M791 440 C801 456 819 456 829 440" stroke="#E8368F" strokeWidth="12" strokeLinecap="round" />
              <path d="M752 440 L676 418 M750 458 L664 458 M752 476 L668 494 M868 440 L944 418 M870 458 L956 458 M868 476 L952 494" stroke="#0A0A0F" strokeWidth="8" strokeLinecap="round" />
            </g>
          </g>
          <g className="cat-parallax-layer doodle doodle-flower" style={{ '--depth': 0.62 } as CSSProperties}>
            <g transform="translate(28 24)">
              <path d="M1118 502 C1078 468 1122 420 1164 454 C1170 398 1238 408 1230 462 C1280 428 1324 490 1270 520 C1328 554 1282 614 1232 574 C1238 632 1166 640 1162 582 C1118 622 1068 562 1118 502 Z" fill="#F7472E" stroke="#0A0A0F" strokeWidth="14" strokeLinejoin="round" />
              <circle cx="1194" cy="516" r="44" fill="#F7D64A" stroke="#0A0A0F" strokeWidth="12" />
              <path d="M1194 564 L1194 696 M1194 640 C1148 620 1124 650 1110 692 M1194 648 C1242 626 1276 646 1290 690" stroke="#0A0A0F" strokeWidth="11" strokeLinecap="round" />
            </g>
          </g>
          <g className="cat-parallax-layer doodle doodle-dog" style={{ '--depth': -0.28 } as CSSProperties}>
            <g transform="translate(40 10)">
              <path d="M1036 152 L1060 94 L1086 94 L1116 140 L1186 140 L1216 94 L1242 94 L1266 152 L1268 236 L1196 288 L1106 288 L1034 236 Z" fill="#F7D64A" stroke="#0A0A0F" strokeWidth="14" strokeLinejoin="round" />
              <path d="M1108 202 L1108 202 M1194 202 L1194 202" stroke="#0A0A0F" strokeWidth="16" strokeLinecap="round" />
              <path d="M1119 238 C1139 262 1165 260 1183 236" stroke="#0A0A0F" strokeWidth="12" strokeLinecap="round" />
              <path d="M1064 208 L1084 208 M1218 208 L1238 208" stroke="#0A0A0F" strokeWidth="10" strokeLinecap="round" />
            </g>
          </g>
          <g className="cat-parallax-layer doodle doodle-apple" style={{ '--depth': 0.5 } as CSSProperties}>
            <g transform="translate(-16 -10)">
              <path d="M210 188 C142 180 116 248 144 304 C168 356 254 366 302 322 C348 280 338 196 272 188 C252 174 228 174 210 188 Z" fill="#FF6A1A" stroke="#0A0A0F" strokeWidth="14" strokeLinejoin="round" />
              <path d="M244 178 L252 116 M262 158 L330 128 L344 178 L290 206 Z" fill="#00D4E8" stroke="#0A0A0F" strokeWidth="13" strokeLinejoin="round" />
            </g>
          </g>
          <g className="cat-parallax-layer doodle doodle-house" style={{ '--depth': -0.58 } as CSSProperties}>
            <g transform="translate(0 36)">
              <path d="M580 756 L714 664 L840 758 Z" fill="#1B2560" stroke="#0A0A0F" strokeWidth="14" strokeLinejoin="round" />
              <path d="M614 746 L810 746 L810 866 L614 866 Z" fill="#FFFFFF" stroke="#0A0A0F" strokeWidth="14" strokeLinejoin="round" />
              <path d="M662 790 L704 790 L704 866 L662 866 Z" fill="#00D4E8" stroke="#0A0A0F" strokeWidth="12" strokeLinejoin="round" />
            </g>
          </g>
          <g className="cat-parallax-layer doodle doodle-star" style={{ '--depth': 0.76 } as CSSProperties}>
            <g transform="translate(70 -30)">
              <path d="M444 134 L472 194 L538 202 L488 244 L504 310 L444 276 L384 310 L400 244 L350 202 L416 194 Z" fill="#FFFFFF" stroke="#0A0A0F" strokeWidth="13" strokeLinejoin="round" />
              <path d="M424 226 L424 226 M466 226 L466 226" stroke="#0A0A0F" strokeWidth="13" strokeLinecap="round" />
              <path d="M426 250 C436 260 454 260 464 250" stroke="#0A0A0F" strokeWidth="8" strokeLinecap="round" />
            </g>
          </g>
          <g className="cat-parallax-layer doodle doodle-swirl" style={{ '--depth': -0.72 } as CSSProperties}>
            <g transform="translate(24 90)">
              <path d="M76 404 C168 318 312 386 244 484 C202 544 94 526 112 442 C126 380 218 390 206 452 C198 492 144 488 152 452" stroke="#0A0A0F" strokeWidth="16" strokeLinecap="round" strokeLinejoin="round" />
              <path d="M64 536 L130 536 M92 504 L92 568 M196 544 L260 574 M218 514 L240 604" stroke="#00D4E8" strokeWidth="14" strokeLinecap="round" />
            </g>
          </g>
          <g className="cat-parallax-layer doodle doodle-sparks" style={{ '--depth': 0.9 } as CSSProperties}>
            <g transform="translate(-270 120)">
              <path d="M936 116 L954 156 L998 162 L964 190 L974 232 L936 210 L898 232 L908 190 L874 162 L918 156 Z" fill="#FFFFFF" stroke="#0A0A0F" strokeWidth="11" strokeLinejoin="round" />
            </g>
            <path d="M944 300 L974 300 M959 284 L959 316 M1064 386 L1100 422 M1098 386 L1064 422 M370 778 L414 778 M392 756 L392 800" stroke="#1B2560" strokeWidth="11" strokeLinecap="round" />
            <g transform="translate(-96 54)">
              <path d="M1014 668 C1018 632 1056 628 1068 652 C1086 614 1138 634 1128 674 C1168 676 1172 726 1128 734 L1026 734 C986 724 988 678 1014 668 Z" fill="#FFFFFF" stroke="#00D4E8" strokeWidth="12" strokeLinejoin="round" />
              <path d="M1049 702 C1067 718 1091 718 1109 700" stroke="#1B2560" strokeWidth="9" strokeLinecap="round" />
            </g>
          </g>
          <g className="cat-parallax-layer doodle doodle-dots" style={{ '--depth': 0.38 } as CSSProperties}>
            <circle cx="535" cy="436" r="9" fill="#E8368F" />
            <circle cx="562" cy="436" r="7" fill="#E8368F" />
            <circle cx="589" cy="436" r="9" fill="#E8368F" />
            <circle cx="548" cy="464" r="7" fill="#E8368F" />
            <circle cx="576" cy="464" r="8" fill="#E8368F" />
          </g>
        </svg>
      </div>

      <div ref={orbitRef} className="scroll-orbit" aria-hidden="true">
        <div ref={trackRef} className="scroll-orbit-track">
          <span ref={fillRef} className="scroll-orbit-fill" />
          <span ref={thumbRef} className="scroll-orbit-thumb">
            <svg className="scroll-orbit-icon scroll-orbit-icon-star" viewBox="0 0 64 64">
              <path d="M32 7C36 7 38 20 41 22C44 24 57 20 59 24C61 28 50 36 49 40C48 44 56 55 52 58C48 61 38 52 34 52C30 52 20 61 16 58C12 55 20 44 19 40C18 36 7 28 9 24C11 20 24 24 27 22C30 20 28 7 32 7Z" fill="#00D4FF" stroke="#F0F0F5" strokeWidth="3.4" strokeLinejoin="round" />
            </svg>
            <svg className="scroll-orbit-icon scroll-orbit-icon-flower" viewBox="0 0 72 72">
              <circle cx="36" cy="17" r="12" fill="#FF8F86" stroke="#1B2560" strokeWidth="5" />
              <circle cx="54" cy="31" r="12" fill="#FF8F86" stroke="#1B2560" strokeWidth="5" />
              <circle cx="47" cy="52" r="12" fill="#FF8F86" stroke="#1B2560" strokeWidth="5" />
              <circle cx="25" cy="52" r="12" fill="#FF8F86" stroke="#1B2560" strokeWidth="5" />
              <circle cx="18" cy="31" r="12" fill="#FF8F86" stroke="#1B2560" strokeWidth="5" />
              <circle cx="36" cy="36" r="14.5" fill="#F8E99A" stroke="#1B2560" strokeWidth="5" />
            </svg>
          </span>
        </div>
      </div>
    </>
  )
}
