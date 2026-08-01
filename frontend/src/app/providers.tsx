'use client'

import { ReactNode } from 'react'
import { AuthSessionProvider } from '@/presentation/components/auth/auth-session-provider'

export function Providers({ children }: { children: ReactNode }) {
  return <AuthSessionProvider>{children}</AuthSessionProvider>
}
