export interface User {
  id: string
  email: string
  name: string
  role: UserRole
  startupId?: string
  status: UserStatus
  createdAt: string
  updatedAt: string
}

export type UserRole = 'admin' | 'startup_owner' | 'member'
export type UserStatus = 'active' | 'pending' | 'inactive'

