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
import { CherryStudio } from '@lobehub/icons'
import { Link } from '@tanstack/react-router'
import { ArrowRight, BookOpen } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { CatPlaygroundVisual } from '../cat-playground-visual'
import { useStatus } from '@/hooks/use-status'
import { HeroTerminalDemo } from '../hero-terminal-demo'

interface HeroProps {
  className?: string
  isAuthenticated?: boolean
}

const MoreIcon = () => (
  <svg
    className='text-foreground/55 group-hover:text-foreground size-6 shrink-0 transition-colors'
    viewBox='0 0 24 24'
    fill='none'
    xmlns='http://www.w3.org/2000/svg'
  >
    <circle cx='6' cy='12' r='2' fill='currentColor' />
    <circle cx='12' cy='12' r='2' fill='currentColor' />
    <circle cx='18' cy='12' r='2' fill='currentColor' />
  </svg>
)

export function Hero(props: HeroProps) {
  const { t } = useTranslation()
  const { status } = useStatus()
  const docsUrl =
    (status?.docs_link as string | undefined) || 'https://docs.newapi.pro'

  const renderDocsButton = () => {
    const isExternal = docsUrl.startsWith('http')
    const className =
      'group h-11 rounded-full border-2 border-foreground bg-background px-5 text-sm font-black text-foreground shadow-[4px_4px_0_var(--foreground)] transition-transform hover:-translate-y-0.5 hover:bg-background hover:shadow-[6px_6px_0_var(--foreground)]'

    if (isExternal) {
      return (
        <Button
          variant='outline'
          className={className}
          render={
            <a href={docsUrl} target='_blank' rel='noopener noreferrer' />
          }
        >
          <BookOpen className='size-4' />
          <span>{t('Docs')}</span>
        </Button>
      )
    }

    return (
      <Button variant='outline' className={className} render={<Link to={docsUrl} />}>
        <BookOpen className='size-4' />
        <span>{t('Docs')}</span>
      </Button>
    )
  }

  return (
    <section className='relative isolate overflow-hidden bg-[#fbfaf6] px-4 pt-20 pb-20 text-foreground md:px-6 md:pt-28 md:pb-24 dark:bg-[#101010]'>
      <div
        aria-hidden
        className='pointer-events-none absolute inset-0 -z-20 bg-[radial-gradient(circle_at_1px_1px,var(--foreground)_1px,transparent_0)] bg-[size:22px_22px] opacity-[0.045]'
      />
      <div
        aria-hidden
        className='pointer-events-none absolute inset-x-0 top-0 -z-10 h-[70%] bg-[linear-gradient(180deg,var(--background)_0%,transparent_100%)] opacity-60'
      />
      <div
        aria-hidden
        className='pointer-events-none absolute right-[8%] bottom-5 -z-10 h-36 w-36 rounded-full border-2 border-dashed border-foreground/15'
      />

      <div className='mx-auto grid max-w-7xl grid-cols-1 items-center gap-12 lg:grid-cols-[minmax(0,0.82fr)_minmax(0,1.18fr)] lg:gap-10'>
        <div className='relative z-20 flex min-w-0 flex-col items-start text-left'>
          <div
            className='landing-animate-fade-up mb-6 inline-flex rotate-[-1deg] items-center gap-2 rounded-full border-2 border-foreground bg-background px-4 py-2 text-[11px] font-black tracking-[0.16em] text-foreground uppercase opacity-0 shadow-[5px_5px_0_var(--foreground)]'
            style={{ animationDelay: '0ms' }}
          >
            <span className='relative flex size-2'>
              <span className='absolute inline-flex h-full w-full animate-ping rounded-full bg-foreground opacity-40' />
              <span className='relative inline-flex size-2 rounded-full bg-foreground' />
            </span>
            <span>{t('AI Application Infrastructure Foundation')}</span>
          </div>

          <h1
            className='landing-animate-fade-up max-w-2xl text-[clamp(2.5rem,5.4vw,4.35rem)] leading-[0.98] font-black tracking-tight text-balance opacity-0'
            style={{ animationDelay: '70ms' }}
          >
            {t('Unified API Gateway for')}
            <br />
            <span className='relative inline-block'>
              {t('Vast Range of AI Models')}
              <span
                aria-hidden
                className='absolute -bottom-2 left-0 h-3 w-full rounded-full bg-foreground/12'
              />
            </span>
          </h1>

          <div
            className='landing-animate-fade-up mt-9 flex flex-wrap items-center gap-3 opacity-0'
            style={{ animationDelay: '150ms' }}
          >
            {props.isAuthenticated ? (
              <>
                <Button
                  className='group h-11 rounded-full border-2 border-foreground bg-foreground px-5 text-sm font-black text-background shadow-[4px_4px_0_var(--muted-foreground)] transition-transform hover:-translate-y-0.5 hover:bg-foreground hover:shadow-[6px_6px_0_var(--muted-foreground)]'
                  render={<Link to='/dashboard' />}
                >
                  {t('Go to Dashboard')}
                  <ArrowRight className='ml-1.5 size-4 transition-transform duration-200 group-hover:translate-x-0.5' />
                </Button>
                {renderDocsButton()}
              </>
            ) : (
              <>
                <Button
                  className='group h-11 rounded-full border-2 border-foreground bg-foreground px-5 text-sm font-black text-background shadow-[4px_4px_0_var(--muted-foreground)] transition-transform hover:-translate-y-0.5 hover:bg-foreground hover:shadow-[6px_6px_0_var(--muted-foreground)]'
                  render={<Link to='/sign-up' />}
                >
                  {t('Get Started')}
                  <ArrowRight className='ml-1.5 size-4 transition-transform duration-200 group-hover:translate-x-0.5' />
                </Button>
                <Button
                  variant='outline'
                  className='h-11 rounded-full border-2 border-foreground bg-background px-5 text-sm font-black text-foreground shadow-[4px_4px_0_var(--foreground)] transition-transform hover:-translate-y-0.5 hover:bg-background hover:shadow-[6px_6px_0_var(--foreground)]'
                  render={<Link to='/pricing' />}
                >
                  {t('View Pricing')}
                </Button>
                {renderDocsButton()}
              </>
            )}
          </div>

          <div
            className='landing-animate-fade-up mt-10 w-full max-w-xl opacity-0'
            style={{ animationDelay: '220ms' }}
          >
            <div className='mb-4 flex flex-col gap-1'>
              <span className='text-[10px] font-black tracking-[0.22em] text-foreground/50 uppercase'>
                {t('Supported Applications')}
              </span>
            </div>
            <div className='flex flex-wrap items-center gap-3'>
              <a
                href='https://cherry-ai.com'
                target='_blank'
                rel='noopener noreferrer'
                className='group flex rotate-[-1deg] items-center gap-3 rounded-full border-2 border-foreground bg-background px-5 py-2.5 text-sm font-black text-foreground shadow-[4px_4px_0_var(--foreground)] transition-transform duration-300 hover:-translate-y-0.5'
              >
                <CherryStudio.Color
                  size={24}
                  className='shrink-0 grayscale'
                />
                <span>Cherry Studio</span>
              </a>

              <a
                href='https://ccswitch.io'
                target='_blank'
                rel='noopener noreferrer'
                className='group flex rotate-[1deg] items-center gap-3 rounded-full border-2 border-foreground bg-background px-5 py-2.5 text-sm font-black text-foreground shadow-[4px_4px_0_var(--foreground)] transition-transform duration-300 hover:-translate-y-0.5'
              >
                <span className='flex size-6 shrink-0 items-center justify-center rounded-md border-2 border-foreground bg-foreground font-mono text-[9px] text-background'>
                  CC
                </span>
                <span>CC Switch</span>
              </a>

              <div className='group flex rotate-[-2deg] cursor-default items-center gap-2.5 rounded-full border-2 border-dashed border-foreground/55 bg-background px-5 py-2.5 text-sm font-black text-foreground/70 transition-colors hover:border-foreground hover:text-foreground'>
                <MoreIcon />
                <span>{t('More Apps')}</span>
              </div>
            </div>
          </div>
        </div>

        <div
          className='landing-animate-fade-up relative z-10 min-w-0 opacity-0'
          style={{ animationDelay: '330ms' }}
        >
          <CatPlaygroundVisual />
          <div className='relative z-20 mx-auto mt-8 w-full max-w-2xl lg:-mt-20 lg:ml-10'>
            <HeroTerminalDemo />
          </div>
        </div>
      </div>
    </section>
  )
}
