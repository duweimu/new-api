/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { cn } from '@/lib/utils'

interface CatPlaygroundVisualProps {
  className?: string
}

export function CatPlaygroundVisual(props: CatPlaygroundVisualProps) {
  return (
    <div
      aria-hidden='true'
      className={cn(
        'pointer-events-none relative mx-auto w-full origin-center',
        props.className
      )}
    >
      <style>{`
        @media (prefers-reduced-motion: no-preference) {
          .cat-playground-frame {
            animation: cat-playground-float 5.8s ease-in-out infinite;
          }
          .cat-playground-image {
            animation: cat-playground-breathe 6.4s ease-in-out infinite;
          }
          .cat-playground-thread path:first-child {
            animation: cat-playground-thread 4.4s linear infinite;
          }
          .cat-playground-sticker {
            animation: cat-playground-sticker 3.8s ease-in-out infinite;
          }
          .cat-playground-paw {
            animation: cat-playground-paw 4.6s ease-in-out infinite;
          }
          .cat-playground-yarn {
            animation: cat-playground-yarn 4.2s ease-in-out infinite;
          }
        }

        @keyframes cat-playground-float {
          0%, 100% { transform: translateY(0) rotate(-0.3deg); }
          50% { transform: translateY(-8px) rotate(0.5deg); }
        }

        @keyframes cat-playground-breathe {
          0%, 100% { transform: scale(1.13); }
          50% { transform: scale(1.18) translateY(-3px); }
        }

        @keyframes cat-playground-thread {
          to { stroke-dashoffset: -64; }
        }

        @keyframes cat-playground-sticker {
          0%, 100% { transform: translateY(0) rotate(var(--sticker-rotate)); }
          50% { transform: translateY(-5px) rotate(0deg); }
        }

        @keyframes cat-playground-paw {
          0%, 100% { transform: translateY(0) rotate(var(--paw-rotate)) scale(var(--paw-scale)); }
          50% { transform: translateY(-6px) rotate(var(--paw-rotate)) scale(var(--paw-scale)); }
        }

        @keyframes cat-playground-yarn {
          0%, 100% { transform: translate(0, 0) rotate(0deg); }
          50% { transform: translate(10px, -7px) rotate(18deg); }
        }
      `}</style>
      <div className='absolute -inset-4 -z-10 rotate-[-1deg] rounded-[2rem] border-2 border-dashed border-foreground/15 bg-[radial-gradient(circle_at_20%_20%,var(--foreground)_1px,transparent_1.5px)] bg-[size:18px_18px] opacity-70' />
      <ThreadPath className='cat-playground-thread absolute top-[11%] left-[7%] h-24 w-56 rotate-[-6deg] text-foreground/55' />
      <ThreadPath className='cat-playground-thread absolute right-[3%] bottom-[16%] h-20 w-48 rotate-[9deg] scale-x-[-1] text-foreground/45' />

      <div className='cat-playground-frame relative overflow-hidden rounded-[2rem] border-2 border-foreground bg-background shadow-[12px_12px_0_var(--foreground)]'>
        <div className='absolute inset-0 bg-[radial-gradient(circle_at_center,var(--foreground)_0.8px,transparent_1px)] bg-[size:16px_16px] opacity-[0.06]' />
        <img
          src='/home/cat-api-playground.png'
          alt=''
          className='cat-playground-image relative z-10 aspect-[16/9] w-full origin-center object-cover object-center grayscale'
          draggable={false}
        />
        <span className='cat-playground-yarn absolute right-[42%] bottom-[15%] z-20 size-8 rounded-full border-2 border-foreground bg-background shadow-[3px_3px_0_var(--foreground)] md:size-10'>
          <span className='absolute top-1/2 left-1/2 h-[2px] w-8 -translate-x-1/2 -translate-y-1/2 rotate-[-24deg] bg-foreground md:w-10' />
          <span className='absolute top-1/2 left-1/2 h-[2px] w-7 -translate-x-1/2 -translate-y-1/2 rotate-[26deg] bg-foreground md:w-9' />
          <span className='absolute top-1/2 left-1/2 h-[2px] w-6 -translate-x-1/2 -translate-y-1/2 rotate-[78deg] bg-foreground md:w-8' />
        </span>
      </div>

      <Sticker
        className='top-4 left-4 [--sticker-rotate:-8deg]'
        label='API'
      />
      <Sticker
        className='right-5 bottom-7 [--sticker-rotate:7deg]'
        label='SSE'
      />
      <Sticker
        className='top-10 right-[19%] [--sticker-rotate:4deg]'
        label='JSON'
      />
      <PawPrint className='cat-playground-paw absolute -bottom-5 left-10 text-foreground/85 [--paw-rotate:-10deg] [--paw-scale:1]' />
      <PawPrint className='cat-playground-paw absolute -top-5 right-8 text-foreground/70 [--paw-rotate:14deg] [--paw-scale:0.75]' />
    </div>
  )
}

function Sticker(props: { className?: string; label: string }) {
  return (
    <div
      className={cn(
        'cat-playground-sticker absolute z-20 rounded-full border-2 border-foreground bg-background px-3 py-1 font-mono text-[11px] font-black tracking-[0.18em] text-foreground shadow-[4px_4px_0_var(--foreground)]',
        props.className
      )}
    >
      {props.label}
    </div>
  )
}

function PawPrint(props: { className?: string }) {
  return (
    <svg
      viewBox='0 0 80 72'
      fill='none'
      className={cn('h-14 w-16', props.className)}
    >
      <ellipse cx='39' cy='49' rx='18' ry='15' fill='currentColor' />
      <ellipse cx='20' cy='30' rx='9' ry='12' fill='currentColor' />
      <ellipse cx='34' cy='18' rx='9' ry='13' fill='currentColor' />
      <ellipse cx='50' cy='18' rx='9' ry='13' fill='currentColor' />
      <ellipse cx='64' cy='30' rx='9' ry='12' fill='currentColor' />
    </svg>
  )
}

function ThreadPath(props: { className?: string }) {
  return (
    <svg viewBox='0 0 240 100' fill='none' className={props.className}>
      <path
        d='M8 68C42 18 70 94 105 48C135 9 159 86 192 42C207 22 221 21 232 28'
        stroke='currentColor'
        strokeWidth='5'
        strokeLinecap='round'
        strokeDasharray='2 14'
      />
      <path
        d='M8 68C42 18 70 94 105 48C135 9 159 86 192 42C207 22 221 21 232 28'
        stroke='currentColor'
        strokeWidth='1.5'
        strokeLinecap='round'
      />
    </svg>
  )
}
