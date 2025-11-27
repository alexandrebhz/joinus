export interface Job {
  id: string
  startupId: string
  startupName?: string
  title: string
  description: string
  requirements: string
  jobType: JobType
  locationType: LocationType
  city: string
  country: string
  salaryMin?: number
  salaryMax?: number
  currency: string
  applicationUrl?: string
  applicationEmail?: string
  status: JobStatus
  expiresAt?: string
  createdAt: string
  updatedAt: string
}

export type JobType = 'full_time' | 'part_time' | 'contract' | 'internship'
export type LocationType = 'remote' | 'hybrid' | 'onsite'
export type JobStatus = 'active' | 'closed'

