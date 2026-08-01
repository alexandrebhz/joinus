import { create } from 'zustand'
import { User } from '@/domain/entities/user.entity'
import { apiClient } from '@/infrastructure/api/api-client'

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  setUser: (user: User | null) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  setUser: (user) => {
    set({ user, isAuthenticated: !!user })
  },
  logout: () => {
    void apiClient.logout()
    set({ user: null, isAuthenticated: false })
  },
}))
