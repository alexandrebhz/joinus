import { AuthResponse, LoginRequest, RegisterRequest } from '../dto/auth.dto'
import { JobResponse, CreateJobRequest, UpdateJobRequest, JobListFilters } from '../dto/job.dto'
import { StartupResponse, CreateStartupRequest, UpdateStartupRequest, StartupListFilters } from '../dto/startup.dto'
import { CreateContactRequest, ContactResponse } from '../dto/contact.dto'
import { CreateCheckoutRequest, CheckoutResponse, BillingStatusResponse } from '../dto/billing.dto'
import { User } from '@/domain/entities/user.entity'
import { ApiResponse } from '@/domain/value-objects/api-response.vo'

export interface IApiClient {
  // Auth
  register(data: RegisterRequest): Promise<ApiResponse<AuthResponse>>
  login(data: LoginRequest): Promise<ApiResponse<AuthResponse>>
  refreshToken(refreshToken: string): Promise<ApiResponse<{ access_token: string }>>
  logout(): Promise<void>
  exchangeOAuthCode(code: string): Promise<ApiResponse<AuthResponse>>
  getCurrentUser(): Promise<ApiResponse<User>>

  // Jobs
  listJobs(filters?: JobListFilters): Promise<ApiResponse<JobResponse[]>>
  getJob(id: string): Promise<ApiResponse<JobResponse>>
  createJob(data: CreateJobRequest): Promise<ApiResponse<JobResponse>>
  updateJob(data: UpdateJobRequest): Promise<ApiResponse<JobResponse>>
  deleteJob(id: string): Promise<void>

  // Startups
  listStartups(filters?: StartupListFilters): Promise<ApiResponse<StartupResponse[]>>
  getStartup(id: string): Promise<ApiResponse<StartupResponse>>
  getStartupBySlug(slug: string): Promise<ApiResponse<StartupResponse>>
  createStartup(data: CreateStartupRequest): Promise<ApiResponse<StartupResponse>>
  updateStartup(data: UpdateStartupRequest): Promise<ApiResponse<StartupResponse>>

  // Files
  uploadFile(file: File): Promise<ApiResponse<{ url: string; id: string }>>

  // Contact
  createContact(data: CreateContactRequest): Promise<ApiResponse<ContactResponse>>

  // Billing
  createCheckout(data: CreateCheckoutRequest): Promise<ApiResponse<CheckoutResponse>>
  getBillingStatus(): Promise<ApiResponse<BillingStatusResponse>>
}

