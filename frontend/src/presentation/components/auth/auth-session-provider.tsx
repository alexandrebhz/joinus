'use client'

import { ReactNode, useEffect, useState } from 'react'
import { apiClient } from '@/infrastructure/api/api-client'
import { useAuthStore } from '@/infrastructure/store/auth.store'

export function AuthSessionProvider({ children }: { children: ReactNode }) {
  const setUser = useAuthStore((s) => s.setUser)
  const [ready, setReady] = useState(false)

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      try {
        const res = await apiClient.getCurrentUser()
        if (!cancelled && res.data) {
          setUser(res.data as any)
        }
      } catch {
        if (!cancelled) setUser(null)
      } finally {
        if (!cancelled) setReady(true)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [setUser])

  // Block protected UI until cookie session is checked (avoids login redirect flash).
  if (!ready) {
    return null
  }

  return <>{children}</>
}
