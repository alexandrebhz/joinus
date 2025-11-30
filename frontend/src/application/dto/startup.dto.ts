import { Startup, StartupStatus } from '@/domain/entities/startup.entity'

export interface CreateStartupRequest {
  name: string
  description: string
  website: string
  founded_year: number
  industry: string
  company_size: string
  location: string
  allow_public_join?: boolean
}

export interface UpdateStartupRequest {
  id: string
  name?: string
  description?: string
  logo_url?: string
  website?: string
  allow_public_join?: boolean
  join_code?: string
  industry?: string
  company_size?: string
  location?: string
}

export interface StartupListFilters {
  industry?: string
  status?: StartupStatus
  search?: string
  location?: string
  company_size?: string
  page?: number
  page_size?: number
}

export type StartupResponse = Startup

