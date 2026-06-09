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
import { useState } from 'react'
import { Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { login } from '@/features/auth/api'
import { useAuthRedirect } from '@/features/auth/hooks/use-auth-redirect'

export function HomeLoginForm() {
  const { t } = useTranslation()
  const { handleLoginSuccess, redirectTo2FA } = useAuthRedirect()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const isSubmitDisabled =
    isLoading || username.trim().length === 0 || password.length === 0

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (isSubmitDisabled) {
      return
    }

    setIsLoading(true)

    try {
      const response = await login({
        username: username.trim(),
        password,
      })

      if (!response.success) {
        toast.error(response.message || t('Login failed'))
        return
      }

      if (response.data?.require_2fa) {
        redirectTo2FA()
        return
      }

      await handleLoginSuccess(response.data as { id?: number } | null, '/dashboard')
    } catch (_error) {
      // Errors are handled by the global interceptor.
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className='flex flex-col gap-4'>
      <label className='flex flex-col gap-2 text-sm font-medium text-foreground'>
        <span>{t('Username or Email')}</span>
        <Input
          name='username'
          autoComplete='username'
          value={username}
          onChange={(event) => setUsername(event.target.value)}
        />
      </label>

      <label className='flex flex-col gap-2 text-sm font-medium text-foreground'>
        <span>{t('Password')}</span>
        <Input
          type='password'
          name='password'
          autoComplete='current-password'
          value={password}
          onChange={(event) => setPassword(event.target.value)}
        />
      </label>

      <Button type='submit' size='lg' disabled={isSubmitDisabled} className='mt-1 w-full'>
        {isLoading ? <Loader2 className='animate-spin' /> : null}
        {t('Sign in')}
      </Button>
    </form>
  )
}
