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

export type UserRole = 'admin' | 'platform_admin' | 'startup_owner' | 'candidate' | 'member' | 'user'
export type UserStatus = 'active' | 'pending' | 'inactive'

