import { Job, JobType, LocationType, JobStatus } from '@/domain/entities/job.entity'

export interface CreateJobRequest {
  startup_id: string
  title: string
  description: string
  requirements: string
  job_type: JobType
  location_type: LocationType
  city?: string
  country: string
  salary_min?: number
  salary_max?: number
  currency: string
  application_url?: string
  application_email?: string
  expires_at?: string
}

export interface UpdateJobRequest {
  id: string
  title?: string
  description?: string
  requirements?: string
  job_type?: JobType
  location_type?: LocationType
  city?: string
  country?: string
  salary_min?: number
  salary_max?: number
  currency?: string
  application_url?: string
  application_email?: string
  status?: JobStatus
  expires_at?: string
}

export interface JobListFilters {
  startup_id?: string
  job_type?: JobType
  location_type?: LocationType
  status?: JobStatus
  search?: string
  country?: string
  city?: string
  salary_min?: number
  salary_max?: number
  currency?: string
  order_by?: string
  order_dir?: 'ASC' | 'DESC'
  page?: number
  page_size?: number
}

export type JobResponse = Job

