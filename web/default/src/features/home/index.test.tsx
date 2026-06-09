import { describe, expect, mock, test } from 'bun:test'
import { renderToStaticMarkup } from 'react-dom/server'

mock.module('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}))

mock.module('@tanstack/react-router', () => ({
  useNavigate: () => () => {},
}))

mock.module('@/stores/auth-store', () => ({
  useAuthStore: () => ({
    auth: {
      user: null,
    },
  }),
}))

mock.module('./hooks', () => ({
  useHomePageContent: () => ({
    content: '',
    isLoaded: true,
    isUrl: false,
  }),
}))

mock.module('@/components/layout', () => ({
  PublicLayout: ({ children }: { children?: React.ReactNode }) => (
    <div>{children}</div>
  ),
}))

mock.module('@/components/layout/components/footer', () => ({
  Footer: () => <footer>Footer</footer>,
}))

mock.module('./components', () => ({
  Hero: () => <section>Get Started View Pricing Docs</section>,
  Stats: () => <section>Stats</section>,
  Features: () => <section>Features</section>,
  HowItWorks: () => <section>How It Works</section>,
  CTA: () => <section>CTA</section>,
}))

mock.module('@/features/auth/hooks/use-auth-redirect', () => ({
  useAuthRedirect: () => ({
    handleLoginSuccess: async () => {},
    redirectTo2FA: () => {},
  }),
}))

mock.module('@/features/auth/api', () => ({
  login: async () => ({ success: true, data: { id: 1 } }),
}))

mock.module('sonner', () => ({
  toast: {
    success: () => {},
    error: () => {},
  },
}))

import { Home } from './index'

describe('Home', () => {
  test('renders only the cat visual and a minimal sign-in form', () => {
    const html = renderToStaticMarkup(<Home />)

    expect(html).toContain('Username or Email')
    expect(html).toContain('Password')
    expect(html).toContain('Sign in')

    expect(html).not.toContain('Get Started')
    expect(html).not.toContain('View Pricing')
    expect(html).not.toContain('Docs')
    expect(html).not.toContain('Forgot password?')
    expect(html).not.toContain('Sign in with Passkey')
  })
})
