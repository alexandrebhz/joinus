export interface Startup {
  id: string
  name: string
  slug: string
  description: string
  logoUrl?: string
  website: string
  foundedYear: number
  industry: string
  companySize: string
  location: string
  allowPublicJoin: boolean
  memberCount?: number
  jobCount?: number
  status: StartupStatus
  createdAt: string
  updatedAt: string
}

export type StartupStatus = 'active' | 'inactive'

