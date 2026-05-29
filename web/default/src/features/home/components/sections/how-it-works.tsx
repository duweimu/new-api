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
import { Settings, Zap, BarChart3 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { AnimateInView } from '@/components/animate-in-view'

export function HowItWorks() {
  const { t } = useTranslation()

  const steps = [
    {
      num: '1',
      title: t('Configure'),
      icon: <Settings className='size-7' strokeWidth={2} />,
    },
    {
      num: '2',
      title: t('Connect'),
      icon: <Zap className='size-7' strokeWidth={2} />,
    },
    {
      num: '3',
      title: t('Monitor'),
      icon: <BarChart3 className='size-7' strokeWidth={2} />,
    },
  ]

  return (
    <section className='relative z-10 overflow-hidden border-t-2 border-foreground bg-background px-6 py-24 text-foreground md:py-32'>
      <div
        aria-hidden
        className='absolute top-28 left-1/2 h-1 w-[80vw] max-w-5xl -translate-x-1/2 border-t-2 border-dashed border-foreground/25 md:top-[17.5rem]'
      />
      <div className='mx-auto max-w-6xl'>
        <AnimateInView className='mb-16 text-center md:mb-20'>
          <p className='mb-3 inline-flex rotate-[1deg] rounded-full border-2 border-foreground bg-background px-4 py-1.5 text-xs font-black tracking-[0.18em] uppercase shadow-[4px_4px_0_var(--foreground)]'>
            {t('How It Works')}
          </p>
          <h2 className='text-3xl font-black tracking-tight md:text-5xl'>
            {t('Three steps to get started')}
          </h2>
        </AnimateInView>

        <div className='grid gap-8 md:grid-cols-3 md:gap-12'>
          {steps.map((step, index) => (
            <AnimateInView
              key={step.num}
              delay={index * 150}
              animation='fade-up'
              className='relative flex flex-col items-center text-center'
            >
              <div className='relative mb-7'>
                <div className='flex size-24 items-center justify-center rounded-[1.75rem] border-2 border-foreground bg-background shadow-[7px_7px_0_var(--foreground)]'>
                  {step.icon}
                </div>
                <div className='absolute -top-3 -right-3 flex size-9 items-center justify-center rounded-full border-2 border-foreground bg-foreground text-sm font-black text-background'>
                  {step.num}
                </div>
                <PawTrail className='absolute -bottom-5 left-1/2 -translate-x-1/2 text-foreground/28' />
              </div>
              <h3 className='text-lg font-black'>{step.title}</h3>
            </AnimateInView>
          ))}
        </div>
      </div>
    </section>
  )
}

function PawTrail(props: { className?: string }) {
  return (
    <svg
      aria-hidden='true'
      viewBox='0 0 116 24'
      fill='none'
      className={`h-6 w-28 ${props.className ?? ''}`}
    >
      {[8, 34, 60, 86, 108].map((x, index) => (
        <g key={x} transform={`translate(${x} ${index % 2 === 0 ? 2 : 6})`}>
          <ellipse cx='0' cy='14' rx='5' ry='4' fill='currentColor' />
          <circle cx='-5' cy='8' r='2.3' fill='currentColor' />
          <circle cx='0' cy='6' r='2.3' fill='currentColor' />
          <circle cx='5' cy='8' r='2.3' fill='currentColor' />
        </g>
      ))}
    </svg>
  )
}
