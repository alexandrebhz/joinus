import axios, { AxiosInstance, AxiosError } from 'axios'
import { IApiClient } from '@/application/ports/api-client.port'
import { AuthResponse, LoginRequest, RegisterRequest } from '@/application/dto/auth.dto'
import { JobResponse, CreateJobRequest, UpdateJobRequest, JobListFilters } from '@/application/dto/job.dto'
import { StartupResponse, CreateStartupRequest, UpdateStartupRequest, StartupListFilters } from '@/application/dto/startup.dto'
import { CreateContactRequest, ContactResponse } from '@/application/dto/contact.dto'
import { CreateCheckoutRequest, CheckoutResponse, BillingStatusResponse } from '@/application/dto/billing.dto'
import { User } from '@/domain/entities/user.entity'
import { ApiResponse, ApiError } from '@/domain/value-objects/api-response.vo'

export class ApiClient implements IApiClient {
  // Transform snake_case API response to camelCase frontend format
  private transformJobData(job: any): JobResponse {
    return {
      id: job.id,
      startupId: job.startup_id,
      startupName: job.startup_name,
      startupSlug: job.startup_slug,
      title: job.title,
      description: job.description,
      requirements: job.requirements,
      jobType: job.job_type,
      locationType: job.location_type,
      city: job.city,
      country: job.country,
      salaryMin: job.salary_min,
      salaryMax: job.salary_max,
      currency: job.currency,
      applicationUrl: job.application_url,
      applicationEmail: job.application_email,
      status: job.status,
      expiresAt: job.expires_at,
      boostedUntil: job.boosted_until,
      createdAt: job.created_at,
      updatedAt: job.updated_at,
    }
  }

  // Attach plan/billing fields (snake_case from the API) onto the startup payload
  private transformStartupData(startup: any): StartupResponse {
    return {
      ...startup,
      plan: startup.plan || 'free',
      planExpiresAt: startup.plan_expires_at,
    }
  }

  private client: AxiosInstance
  private baseURL: string

  private normalizeApiUrl(url: string | undefined): string {
    // Handle undefined/null/empty values
    if (!url || typeof url !== 'string') {
      return 'http://localhost:8080'
    }
    
    // Remove trailing slashes
    let normalized = url.trim().replace(/\/+$/, '')
    
    // If URL contains the frontend domain incorrectly prepended, try to extract the backend URL
    // Example: https://joinus.ie/joinus-production.up.railway.app -> https://joinus-production.up.railway.app
    // Pattern 1: https://frontend.com/https://backend.com
    const frontendDomainMatch1 = normalized.match(/https?:\/\/[^\/]+\/(https?:\/\/.+)/)
    if (frontendDomainMatch1) {
      normalized = frontendDomainMatch1[1]
    }
    
    // Pattern 2: https://frontend.com/backend-domain.com (without protocol in path)
    // Check if URL looks like frontend-domain/backend-domain pattern
    const urlParts = normalized.split('/')
    if (urlParts.length >= 4 && urlParts[2] && urlParts[3]) {
      const frontendDomain = urlParts[2] // e.g., joinus.ie
      const backendDomain = urlParts[3] // e.g., joinus-production.up.railway.app
      
      // If backend domain looks like a Railway/backend domain and frontend domain is different
      if (backendDomain.includes('.railway.app') || backendDomain.includes('.up.railway.app')) {
        // Extract just the backend domain and reconstruct URL
        normalized = `https://${backendDomain}`
      }
    }
    
    // Ensure URL has protocol
    if (!normalized.match(/^https?:\/\//)) {
      // For Railway/production domains, use HTTPS. For localhost, use HTTP
      if (normalized.includes('.railway.app') || normalized.includes('.up.railway.app') || normalized.includes('railway')) {
        normalized = `https://${normalized}`
      } else if (normalized.includes('localhost') || normalized.includes('127.0.0.1')) {
        normalized = `http://${normalized}`
      } else {
        // Default to HTTPS for production domains
        normalized = `https://${normalized}`
      }
    }
    
    return normalized
  }

  constructor(baseURL?: string) {
    // Browser: same-origin BFF proxy (HttpOnly cookies). SSR: direct backend URL.
    const isBrowser = typeof window !== 'undefined'
    let absoluteBaseURL: string

    if (isBrowser && !baseURL) {
      this.baseURL = '/api/backend'
      absoluteBaseURL = '/api/backend'
    } else {
      const rawUrl = baseURL || process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
      const normalizedUrl = this.normalizeApiUrl(rawUrl)
      this.baseURL = normalizedUrl
      if (normalizedUrl.startsWith('http://') || normalizedUrl.startsWith('https://')) {
        absoluteBaseURL = `${normalizedUrl}/api/v1`
      } else {
        absoluteBaseURL = `https://${normalizedUrl}/api/v1`
      }
      if (!absoluteBaseURL.match(/^https?:\/\//)) {
        throw new Error(`Invalid API baseURL: ${absoluteBaseURL}`)
      }
    }

    this.client = axios.create({
      baseURL: absoluteBaseURL,
      headers: {
        'Content-Type': 'application/json',
      },
      withCredentials: isBrowser,
    })

    this.client.interceptors.request.use(
      (config) => {
        if (!isBrowser) {
          const internalKey = process.env.API_INTERNAL_KEY
          if (internalKey) {
            config.headers['X-Internal-Key'] = internalKey
          }
        }
        return config
      },
      (error) => Promise.reject(error)
    )

    this.client.interceptors.response.use(
      (response) => response,
      async (error: AxiosError<ApiError>) => {
        if (error.response?.status === 401 && isBrowser) {
          const requestUrl = error.config?.url || ''
          // AuthSessionProvider probes /me on every page; a 401 means "logged out", not "force login".
          const isSessionProbe =
            requestUrl === '/me' ||
            requestUrl.endsWith('/me') ||
            requestUrl.includes('/me?')
          const path = window.location.pathname
          const isProtectedRoute = path.startsWith('/dashboard')
          // Only hard-redirect when an authenticated action fails on a protected page.
          // Public pages (home, jobs, etc.) stay put; dashboard route guards also cover this.
          if (!isSessionProbe && isProtectedRoute) {
            window.location.href = '/login'
          }
        }
        return Promise.reject(error)
      }
    )
  }

  private async request<T>(config: { method: string; url: string; data?: any; params?: any }): Promise<ApiResponse<T>> {
    try {
      const response = await this.client.request<ApiResponse<T>>(config)
      return response.data
    } catch (error) {
      const axiosError = error as AxiosError<ApiError>
      throw new Error(axiosError.response?.data?.error || axiosError.message || 'An error occurred')
    }
  }

  // Auth methods
  private toAuthResponse(response: ApiResponse<{
    access_token?: string
    refresh_token?: string
    user: any
  }>): ApiResponse<AuthResponse> {
    // Tokens are stored in HttpOnly cookies by the BFF proxy and redacted from JSON.
    return {
      ...response,
      data: {
        accessToken: response.data.access_token || '',
        refreshToken: response.data.refresh_token || '',
        user: response.data.user,
      },
    }
  }

  async register(data: RegisterRequest): Promise<ApiResponse<AuthResponse>> {
    const response = await this.request<{
      access_token?: string
      refresh_token?: string
      user: any
    }>({
      method: 'POST',
      url: '/auth/register',
      data,
    })
    return this.toAuthResponse(response)
  }

  async login(data: LoginRequest): Promise<ApiResponse<AuthResponse>> {
    const response = await this.request<{
      access_token?: string
      refresh_token?: string
      user: any
    }>({
      method: 'POST',
      url: '/auth/login',
      data,
    })
    return this.toAuthResponse(response)
  }

  async refreshToken(refreshToken: string): Promise<ApiResponse<{ access_token: string }>> {
    return this.request<{ access_token: string }>({
      method: 'POST',
      url: '/auth/refresh',
      data: { refresh_token: refreshToken },
    })
  }

  async logout(): Promise<void> {
    try {
      await this.request({ method: 'POST', url: '/auth/logout' })
    } catch {
      // clear cookies even if network fails
      if (typeof window !== 'undefined') {
        await fetch('/api/auth/session', { method: 'DELETE' })
      }
    }
  }

  /** Exchange a one-time OAuth login code for a session (tokens in HttpOnly cookies). */
  async exchangeOAuthCode(code: string): Promise<ApiResponse<AuthResponse>> {
    const response = await this.request<{
      access_token?: string
      refresh_token?: string
      user: any
    }>({
      method: 'POST',
      url: '/auth/oauth/exchange',
      data: { code },
    })
    return this.toAuthResponse(response)
  }

  async getCurrentUser(): Promise<ApiResponse<User>> {
    return this.request<User>({
      method: 'GET',
      url: '/me',
    })
  }

  // Job methods
  async listJobs(filters?: JobListFilters): Promise<ApiResponse<JobResponse[]>> {
    const response = await this.request<any>({
      method: 'GET',
      url: '/jobs',
      params: filters,
    })
    
    // Transform snake_case API response to camelCase frontend format
    if (response.data && Array.isArray(response.data)) {
      response.data = response.data.map(job => this.transformJobData(job))
    }
    
    // Transform pagination meta from snake_case to camelCase
    if (response.meta) {
      const meta = response.meta as any
      response.meta = {
        page: meta.page,
        pageSize: meta.page_size,
        totalPages: meta.total_pages,
        totalCount: meta.total_count,
      }
    }
    
    return response as ApiResponse<JobResponse[]>
  }

  async getJob(id: string): Promise<ApiResponse<JobResponse>> {
    const response = await this.request<any>({
      method: 'GET',
      url: `/jobs/${id}`,
    })
    
    // Transform snake_case API response to camelCase frontend format
    if (response.data) {
      response.data = this.transformJobData(response.data)
    }
    
    return response as ApiResponse<JobResponse>
  }

  async createJob(data: CreateJobRequest): Promise<ApiResponse<JobResponse>> {
    return this.request<JobResponse>({
      method: 'POST',
      url: '/jobs',
      data,
    })
  }

  async updateJob(data: UpdateJobRequest): Promise<ApiResponse<JobResponse>> {
    const { id, ...updateData } = data
    return this.request<JobResponse>({
      method: 'PUT',
      url: `/jobs/${id}`,
      data: updateData,
    })
  }

  async deleteJob(id: string): Promise<void> {
    await this.request({
      method: 'DELETE',
      url: `/jobs/${id}`,
    })
  }

  // Startup methods
  async listStartups(filters?: StartupListFilters): Promise<ApiResponse<StartupResponse[]>> {
    const response = await this.request<any>({
      method: 'GET',
      url: '/startups',
      params: filters,
    })
    
    if (response.data && Array.isArray(response.data)) {
      response.data = response.data.map((startup) => this.transformStartupData(startup))
    }

    // Transform pagination meta from snake_case to camelCase
    if (response.meta) {
      const meta = response.meta as any
      response.meta = {
        page: meta.page,
        pageSize: meta.page_size,
        totalPages: meta.total_pages,
        totalCount: meta.total_count,
      }
    }
    
    return response as ApiResponse<StartupResponse[]>
  }

  async getStartup(id: string): Promise<ApiResponse<StartupResponse>> {
    const response = await this.request<any>({
      method: 'GET',
      url: `/startups/${id}`,
    })
    if (response.data) {
      response.data = this.transformStartupData(response.data)
    }
    return response as ApiResponse<StartupResponse>
  }

  async getStartupBySlug(slug: string): Promise<ApiResponse<StartupResponse>> {
    const response = await this.request<any>({
      method: 'GET',
      url: `/startups/slug/${slug}`,
    })
    if (response.data) {
      response.data = this.transformStartupData(response.data)
    }
    return response as ApiResponse<StartupResponse>
  }

  async createStartup(data: CreateStartupRequest): Promise<ApiResponse<StartupResponse>> {
    const response = await this.request<any>({
      method: 'POST',
      url: '/startups',
      data,
    })
    if (response.data) {
      response.data = this.transformStartupData(response.data)
    }
    return response as ApiResponse<StartupResponse>
  }

  async updateStartup(data: UpdateStartupRequest): Promise<ApiResponse<StartupResponse>> {
    const { id, ...updateData } = data
    const response = await this.request<any>({
      method: 'PUT',
      url: `/startups/${id}`,
      data: updateData,
    })
    if (response.data) {
      response.data = this.transformStartupData(response.data)
    }
    return response as ApiResponse<StartupResponse>
  }

  // File methods
  async uploadFile(file: File): Promise<ApiResponse<{ url: string; id: string }>> {
    const formData = new FormData()
    formData.append('file', file)
    
    return this.request<{ url: string; id: string }>({
      method: 'POST',
      url: '/upload',
      data: formData,
    })
  }

  // Contact methods
  async createContact(data: CreateContactRequest): Promise<ApiResponse<ContactResponse>> {
    return this.request<ContactResponse>({
      method: 'POST',
      url: '/contact',
      data,
    })
  }

  // Billing methods
  async createCheckout(data: CreateCheckoutRequest): Promise<ApiResponse<CheckoutResponse>> {
    const response = await this.request<any>({
      method: 'POST',
      url: '/billing/checkout',
      data,
    })
    return {
      ...response,
      data: {
        url: response.data.url,
        sessionId: response.data.session_id,
      },
    }
  }

  async getBillingStatus(): Promise<ApiResponse<BillingStatusResponse>> {
    const response = await this.request<any>({
      method: 'GET',
      url: '/billing/status',
    })
    const startups = (response.data?.startups || []).map((s: any) => ({
      startupId: s.startup_id,
      startupName: s.startup_name,
      plan: s.plan,
      planExpiresAt: s.plan_expires_at,
    }))
    return {
      ...response,
      data: { startups },
    }
  }
}

// Singleton instance
export const apiClient = new ApiClient()

