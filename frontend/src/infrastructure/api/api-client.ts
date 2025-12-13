import axios, { AxiosInstance, AxiosError } from 'axios'
import { IApiClient } from '@/application/ports/api-client.port'
import { AuthResponse, LoginRequest, RegisterRequest } from '@/application/dto/auth.dto'
import { JobResponse, CreateJobRequest, UpdateJobRequest, JobListFilters } from '@/application/dto/job.dto'
import { StartupResponse, CreateStartupRequest, UpdateStartupRequest, StartupListFilters } from '@/application/dto/startup.dto'
import { CreateContactRequest, ContactResponse } from '@/application/dto/contact.dto'
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
      createdAt: job.created_at,
      updatedAt: job.updated_at,
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
    // Safely get the URL, handling undefined values
    const rawUrl = baseURL || process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
    
    // Log the raw environment variable value for debugging (only in browser)
    if (typeof window !== 'undefined') {
      console.log('[ApiClient] Environment check:', {
        'process.env.NEXT_PUBLIC_API_URL': process.env.NEXT_PUBLIC_API_URL,
        'rawUrl parameter': rawUrl,
        'window.location.origin': window.location.origin
      })
    }
    
    const normalizedUrl = this.normalizeApiUrl(rawUrl)
    
    // Warn if URL was normalized (indicates misconfiguration)
    if (rawUrl !== normalizedUrl && typeof window !== 'undefined') {
      console.warn('[ApiClient] API URL was normalized:', {
        original: rawUrl,
        normalized: normalizedUrl,
        message: 'NEXT_PUBLIC_API_URL appears to be incorrectly formatted. Please set it to the backend URL only (e.g., https://joinus-production.up.railway.app)'
      })
    }
    
    this.baseURL = normalizedUrl
    
    // Always log in production to help debug URL issues
    if (typeof window !== 'undefined') {
      console.log('[ApiClient] Final baseURL:', this.baseURL)
    }
    
    // Ensure baseURL is absolute (starts with http:// or https://)
    // This is critical - axios will treat URLs without protocol as relative to current origin
    let absoluteBaseURL: string
    if (this.baseURL.startsWith('http://') || this.baseURL.startsWith('https://')) {
      absoluteBaseURL = `${this.baseURL}/api/v1`
    } else {
      // If no protocol, assume HTTPS for production
      absoluteBaseURL = `https://${this.baseURL}/api/v1`
      console.warn('[ApiClient] baseURL missing protocol, assuming HTTPS:', absoluteBaseURL)
    }
    
    // Final validation - ensure it's truly absolute
    if (!absoluteBaseURL.match(/^https?:\/\//)) {
      throw new Error(`Invalid API baseURL: ${absoluteBaseURL}. Must start with http:// or https://`)
    }
    
    this.client = axios.create({
      baseURL: absoluteBaseURL,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    // Request interceptor to add auth token and log URLs for debugging
    this.client.interceptors.request.use(
      (config) => {
        // Log the actual URL being used (helpful for debugging)
        if (typeof window !== 'undefined') {
          const fullUrl = config.baseURL && config.url
            ? `${config.baseURL}${config.url.startsWith('/') ? '' : '/'}${config.url}`
            : config.url
          
          // Log in production if URL looks wrong
          if (fullUrl && fullUrl.includes('joinus.ie') && fullUrl.includes('railway.app')) {
            console.error('[ApiClient] Detected malformed URL:', {
              fullUrl,
              baseURL: config.baseURL,
              url: config.url,
              envVar: process.env.NEXT_PUBLIC_API_URL,
              normalizedBaseURL: this.baseURL
            })
          }
        }
        
        if (typeof window !== 'undefined') {
          const token = localStorage.getItem('access_token')
          if (token) {
            config.headers.Authorization = `Bearer ${token}`
          }
        }
        return config
      },
      (error) => Promise.reject(error)
    )

    // Response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      async (error: AxiosError<ApiError>) => {
        if (error.response?.status === 401 && typeof window !== 'undefined') {
          // Try to refresh token
          const refreshToken = localStorage.getItem('refresh_token')
          if (refreshToken) {
            try {
              const response = await this.refreshToken(refreshToken)
              if (response.success && response.data.access_token) {
                localStorage.setItem('access_token', response.data.access_token)
                // Retry original request
                if (error.config) {
                  error.config.headers.Authorization = `Bearer ${response.data.access_token}`
                  return this.client.request(error.config)
                }
              }
            } catch {
              // Refresh failed, clear tokens and redirect to login
              localStorage.removeItem('access_token')
              localStorage.removeItem('refresh_token')
              if (window.location.pathname !== '/login') {
                window.location.href = '/login'
              }
            }
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
  async register(data: RegisterRequest): Promise<ApiResponse<AuthResponse>> {
    const response = await this.request<{
      access_token: string
      refresh_token: string
      user: any
    }>({
      method: 'POST',
      url: '/auth/register',
      data,
    })
    
    // Store tokens
    if (typeof window !== 'undefined' && response.data) {
      localStorage.setItem('access_token', response.data.access_token)
      localStorage.setItem('refresh_token', response.data.refresh_token)
    }
    
    // Transform to match AuthResponse interface
    return {
      ...response,
      data: {
        accessToken: response.data.access_token,
        refreshToken: response.data.refresh_token,
        user: response.data.user,
      },
    }
  }

  async login(data: LoginRequest): Promise<ApiResponse<AuthResponse>> {
    const response = await this.request<{
      access_token: string
      refresh_token: string
      user: any
    }>({
      method: 'POST',
      url: '/auth/login',
      data,
    })
    
    // Store tokens and transform response
    if (typeof window !== 'undefined' && response.data) {
      localStorage.setItem('access_token', response.data.access_token)
      localStorage.setItem('refresh_token', response.data.refresh_token)
    }
    
    // Transform to match AuthResponse interface
    return {
      ...response,
      data: {
        accessToken: response.data.access_token,
        refreshToken: response.data.refresh_token,
        user: response.data.user,
      },
    }
  }

  async refreshToken(refreshToken: string): Promise<ApiResponse<{ access_token: string }>> {
    const response = await this.request<{ access_token: string }>({
      method: 'POST',
      url: '/auth/refresh',
      data: { refresh_token: refreshToken },
    })
    return response
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
    return this.request<StartupResponse>({
      method: 'GET',
      url: `/startups/${id}`,
    })
  }

  async getStartupBySlug(slug: string): Promise<ApiResponse<StartupResponse>> {
    return this.request<StartupResponse>({
      method: 'GET',
      url: `/startups/slug/${slug}`,
    })
  }

  async createStartup(data: CreateStartupRequest): Promise<ApiResponse<StartupResponse>> {
    return this.request<StartupResponse>({
      method: 'POST',
      url: '/startups',
      data,
    })
  }

  async updateStartup(data: UpdateStartupRequest): Promise<ApiResponse<StartupResponse>> {
    const { id, ...updateData } = data
    return this.request<StartupResponse>({
      method: 'PUT',
      url: `/startups/${id}`,
      data: updateData,
    })
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
}

// Singleton instance
export const apiClient = new ApiClient()

