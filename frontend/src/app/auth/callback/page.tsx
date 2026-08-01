'use client'

import { useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { Suspense } from 'react'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { apiClient } from '@/infrastructure/api/api-client'
import { useAuthStore } from '@/infrastructure/store/auth.store'
import { User, UserRole, UserStatus } from '@/domain/entities/user.entity'

function AuthCallbackContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { setUser } = useAuthStore()
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const run = async () => {
      const accessToken = searchParams.get('access_token')
      const refreshToken = searchParams.get('refresh_token')

      if (!accessToken || !refreshToken) {
        setError('Missing authentication tokens. Please try signing in again.')
        return
      }

      localStorage.setItem('access_token', accessToken)
      localStorage.setItem('refresh_token', refreshToken)

      try {
        const response = await apiClient.getCurrentUser()
        if (response.success && response.data) {
          const raw = response.data as unknown as Record<string, unknown>
          const user: User = {
            id: String(raw.id),
            email: String(raw.email),
            name: String(raw.name),
            role: String(raw.role) as UserRole,
            status: (String(raw.status || 'active') as UserStatus),
            createdAt: String(raw.created_at || new Date().toISOString()),
            updatedAt: String(raw.updated_at || new Date().toISOString()),
          }
          setUser(user)
          router.replace('/dashboard')
          return
        }
        setError('Could not load your profile after Google sign-in.')
      } catch {
        setError('Google sign-in succeeded but profile load failed. Try signing in again.')
      }
    }

    void run()
  }, [searchParams, setUser, router])

  return (
    <div className="w-full max-w-md text-center space-y-3">
      {error ? (
        <>
          <p className="text-error-600 text-sm">{error}</p>
          <a href="/login" className="text-primary-600 hover:text-primary-700 font-medium text-sm">
            Back to login
          </a>
        </>
      ) : (
        <p className="text-secondary-600 text-sm">Completing Google sign-in…</p>
      )}
    </div>
  )
}

export default function AuthCallbackPage() {
  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 flex items-center justify-center py-12 px-4">
        <Suspense fallback={<p className="text-secondary-600 text-sm">Loading…</p>}>
          <AuthCallbackContent />
        </Suspense>
      </main>
      <Footer />
    </div>
  )
}
