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
import {
  Zap,
  Shield,
  Globe,
  Code,
  Gauge,
  DollarSign,
  Users,
  HeartHandshake,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { AnimateInView } from '@/components/animate-in-view'

interface FeaturesProps {
  className?: string
}

export function Features(_props: FeaturesProps) {
  const { t } = useTranslation()

  const features = [
    {
      id: 'fast',
      num: '01',
      title: t('Lightning Fast'),
      span: 'md:col-span-2',
      icon: <Zap className='size-5' strokeWidth={2} />,
      visual: <ModelStickerGrid />,
    },
    {
      id: 'secure',
      num: '02',
      title: t('Secure & Reliable'),
      span: 'md:col-span-1',
      icon: <Shield className='size-5' strokeWidth={2} />,
      visual: <ShieldStamp />,
    },
    {
      id: 'global',
      num: '03',
      title: t('Global Coverage'),
      span: 'md:col-span-1',
      icon: <Globe className='size-5' strokeWidth={2} />,
      visual: (
        <FlowSteps
          steps={[t('Load Balancing'), t('Rate Limiting'), t('Cost Tracking')]}
        />
      ),
    },
    {
      id: 'developer',
      num: '04',
      title: t('Developer Friendly'),
      span: 'md:col-span-2',
      icon: <Code className='size-5' strokeWidth={2} />,
      visual: <ToolStack label={t('Multi-protocol Compatible')} />,
    },
  ]

  const additionalFeatures = [
    {
      icon: <Gauge className='size-5' strokeWidth={2} />,
      title: t('High Performance'),
    },
    {
      icon: <DollarSign className='size-5' strokeWidth={2} />,
      title: t('Transparent Billing'),
    },
    {
      icon: <Users className='size-5' strokeWidth={2} />,
      title: t('Team Collaboration'),
    },
    {
      icon: <HeartHandshake className='size-5' strokeWidth={2} />,
      title: t('Open Source'),
    },
  ]

  return (
    <section className='relative z-10 overflow-hidden bg-[#fbfaf6] px-6 py-24 text-foreground md:py-32 dark:bg-[#101010]'>
      <div
        aria-hidden
        className='absolute inset-0 bg-[linear-gradient(90deg,var(--foreground)_1px,transparent_1px),linear-gradient(180deg,var(--foreground)_1px,transparent_1px)] bg-[size:64px_64px] opacity-[0.035]'
      />
      <div className='relative mx-auto max-w-6xl'>
        <AnimateInView className='mb-16 max-w-xl'>
          <p className='mb-3 inline-flex rotate-[-2deg] rounded-full border-2 border-foreground bg-background px-4 py-1.5 text-xs font-black tracking-[0.18em] uppercase shadow-[4px_4px_0_var(--foreground)]'>
            {t('Core Features')}
          </p>
          <h2 className='text-3xl leading-tight font-black tracking-tight md:text-5xl'>
            {t('Built for developers,')}
            <br />
            {t('designed for scale')}
          </h2>
        </AnimateInView>

        <div className='grid gap-5 md:grid-cols-3'>
          {features.map((feature, index) => (
            <AnimateInView
              key={feature.id}
              delay={index * 100}
              animation='scale-in'
              className={`group relative overflow-hidden rounded-[1.5rem] border-2 border-foreground bg-background p-7 shadow-[8px_8px_0_var(--foreground)] transition-transform duration-300 hover:-translate-y-1 md:p-8 ${feature.span}`}
            >
              <div
                aria-hidden
                className='absolute inset-0 bg-[radial-gradient(circle_at_1px_1px,var(--foreground)_1px,transparent_0)] bg-[size:18px_18px] opacity-[0.035]'
              />
              <div className='relative'>
                <div className='mb-4 flex items-center gap-3'>
                  <span className='flex size-9 items-center justify-center rounded-full border-2 border-foreground bg-foreground text-[11px] font-black text-background tabular-nums'>
                    {feature.num}
                  </span>
                  <span className='flex size-10 items-center justify-center rounded-xl border-2 border-foreground bg-background'>
                    {feature.icon}
                  </span>
                  <h3 className='text-base font-black'>{feature.title}</h3>
                </div>
                {feature.visual}
              </div>
            </AnimateInView>
          ))}
        </div>

        <div className='mt-12 grid grid-cols-2 gap-5 md:grid-cols-4'>
          {additionalFeatures.map((feature, index) => (
            <AnimateInView
              key={feature.title}
              delay={index * 100}
              animation='fade-up'
              className='relative flex min-h-48 rotate-[-1deg] flex-col rounded-[1.25rem] border-2 border-foreground bg-background p-5 shadow-[6px_6px_0_var(--foreground)] even:rotate-[1deg]'
            >
              <div className='mb-4 flex size-12 items-center justify-center rounded-xl border-2 border-foreground bg-foreground text-background'>
                {feature.icon}
              </div>
              <h3 className='text-sm font-black'>{feature.title}</h3>
            </AnimateInView>
          ))}
        </div>
      </div>
    </section>
  )
}

function ModelStickerGrid() {
  return (
    <div className='mt-5 grid grid-cols-3 gap-2'>
      {['OpenAI', 'Claude', 'Gemini', 'DeepSeek', 'Qwen', 'Llama'].map(
        (name, index) => (
          <div
            key={name}
            className='flex items-center justify-center rounded-full border-2 border-foreground bg-background px-3 py-2 text-[11px] font-black shadow-[3px_3px_0_var(--foreground)]'
            style={{ transform: `rotate(${index % 2 === 0 ? -2 : 2}deg)` }}
          >
            {name}
          </div>
        )
      )}
    </div>
  )
}

function ShieldStamp() {
  return (
    <div className='mt-5 flex items-center justify-center'>
      <div className='relative flex size-24 rotate-[-8deg] items-center justify-center rounded-full border-4 border-double border-foreground bg-background'>
        <Shield className='size-10' strokeWidth={2} />
        <span className='absolute -right-2 -bottom-2 flex size-8 items-center justify-center rounded-full border-2 border-foreground bg-foreground text-xs font-black text-background'>
          <span aria-hidden='true'>
            OK
          </span>
        </span>
      </div>
    </div>
  )
}

function FlowSteps(props: { steps: string[] }) {
  return (
    <div className='mt-5 space-y-3'>
      {props.steps.map((step, index) => (
        <div key={step} className='flex items-center gap-3'>
          <span className='flex size-8 items-center justify-center rounded-full border-2 border-foreground bg-background text-xs font-black'>
            {index + 1}
          </span>
          <span className='h-0.5 flex-1 border-t-2 border-dashed border-foreground/50' />
          <span className='text-xs font-black text-foreground/68'>{step}</span>
        </div>
      ))}
    </div>
  )
}

function ToolStack(props: { label: string }) {
  return (
    <div className='mt-5 flex flex-wrap items-center gap-3'>
      {['API', 'SDK', 'CLI', 'Docs'].map((name) => (
        <span
          key={name}
          className='flex size-11 items-center justify-center rounded-full border-2 border-foreground bg-background font-mono text-[10px] font-black shadow-[3px_3px_0_var(--foreground)]'
        >
          {name}
        </span>
      ))}
      <div className='flex items-center gap-2 rounded-full border-2 border-dashed border-foreground/50 px-4 py-2 text-xs font-black text-foreground/68'>
        <Code className='size-4' />
        {props.label}
      </div>
    </div>
  )
}
