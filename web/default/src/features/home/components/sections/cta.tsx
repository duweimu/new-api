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
import { Link } from '@tanstack/react-router'
import { ArrowRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { AnimateInView } from '@/components/animate-in-view'
import { Button } from '@/components/ui/button'

interface CTAProps {
  className?: string
  isAuthenticated?: boolean
}

export function CTA(props: CTAProps) {
  const { t } = useTranslation()

  if (props.isAuthenticated) {
    return null
  }

  return (
    <section className='relative z-10 overflow-hidden bg-[#fbfaf6] px-6 py-24 text-foreground md:py-32 dark:bg-[#101010]'>
      <div
        aria-hidden
        className='absolute inset-x-0 top-1/2 h-1 -translate-y-1/2 border-t-2 border-dashed border-foreground/20'
      />
      <div
        aria-hidden
        className='absolute top-10 left-[8%] size-24 rotate-12 rounded-full border-2 border-dashed border-foreground/20'
      />
      <div
        aria-hidden
        className='absolute right-[8%] bottom-12 size-32 rounded-full border-2 border-foreground/10 bg-[radial-gradient(circle,var(--foreground)_1px,transparent_1.5px)] bg-[size:14px_14px]'
      />

      <AnimateInView
        className='relative mx-auto max-w-3xl rounded-[2rem] border-2 border-foreground bg-background px-6 py-12 text-center shadow-[10px_10px_0_var(--foreground)] md:px-12 md:py-16'
        animation='scale-in'
      >
        <CatToken />
        <h2 className='text-3xl leading-tight font-black tracking-tight md:text-5xl'>
          {t('Ready to simplify')}
          <br />
          <span className='relative inline-block'>
            {t('your AI integration?')}
            <span
              aria-hidden
              className='absolute -bottom-2 left-0 h-3 w-full rounded-full bg-foreground/12'
            />
          </span>
        </h2>
        <div className='mt-8 flex flex-wrap items-center justify-center gap-3'>
          <Button
            className='group rounded-full border-2 border-foreground bg-foreground px-5 font-black text-background shadow-[4px_4px_0_var(--muted-foreground)] transition-transform hover:-translate-y-0.5 hover:bg-foreground hover:shadow-[6px_6px_0_var(--muted-foreground)]'
            render={<Link to='/sign-up' />}
          >
            {t('Get Started')}
            <ArrowRight className='ml-1 size-3.5 transition-transform duration-200 group-hover:translate-x-0.5' />
          </Button>
          <Button
            variant='outline'
            className='rounded-full border-2 border-foreground bg-background px-5 font-black text-foreground shadow-[4px_4px_0_var(--foreground)] transition-transform hover:-translate-y-0.5 hover:bg-background hover:shadow-[6px_6px_0_var(--foreground)]'
            render={<Link to='/pricing' />}
          >
            {t('View Pricing')}
          </Button>
        </div>
      </AnimateInView>
    </section>
  )
}

function CatToken() {
  return (
    <div
      aria-hidden='true'
      className='absolute -top-8 left-1/2 flex size-16 -translate-x-1/2 rotate-[-8deg] items-center justify-center rounded-full border-2 border-foreground bg-background shadow-[5px_5px_0_var(--foreground)]'
    >
      <span className='relative flex size-8 rounded-full bg-foreground'>
        <span className='absolute -top-2 left-1 size-4 rotate-[-20deg] rounded-tl-full bg-foreground' />
        <span className='absolute -top-2 right-1 size-4 rotate-[20deg] rounded-tr-full bg-foreground' />
        <span className='absolute top-3 left-2 size-1 rounded-full bg-background' />
        <span className='absolute top-3 right-2 size-1 rounded-full bg-background' />
      </span>
    </div>
  )
}
