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
import { useEffect } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { useAuthStore } from '@/stores/auth-store'
import { CatPlaygroundVisual } from './components/cat-playground-visual'
import { HomeLoginForm } from './components/home-login-form'

export function Home() {
  const { auth } = useAuthStore()
  const navigate = useNavigate()

  useEffect(() => {
    if (auth.user) {
      navigate({ to: '/dashboard', replace: true })
    }
  }, [auth.user, navigate])

  if (auth.user) {
    return <main className='min-h-screen bg-[#fbfaf6] dark:bg-[#101010]' />
  }

  return (
    <main className='min-h-screen overflow-hidden bg-[#fbfaf6] px-4 py-8 text-foreground md:px-6 md:py-10 dark:bg-[#101010]'>
      <div
        aria-hidden
        className='pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_1px_1px,var(--foreground)_1px,transparent_0)] bg-[size:22px_22px] opacity-[0.045]'
      />

      <div className='relative mx-auto grid min-h-[calc(100vh-4rem)] w-full max-w-6xl items-center gap-10 lg:grid-cols-[minmax(0,1.2fr)_24rem] lg:gap-14'>
        <section className='min-w-0'>
          <CatPlaygroundVisual className='max-w-4xl' />
        </section>

        <section className='min-w-0'>
          <div className='border-border bg-background/92 rounded-2xl border p-6 shadow-[0_18px_60px_-28px_rgba(0,0,0,0.35)] backdrop-blur'>
            <HomeLoginForm />
          </div>
        </section>
      </div>
    </main>
  )
}
