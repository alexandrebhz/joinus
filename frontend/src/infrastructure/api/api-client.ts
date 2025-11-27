import axios, { AxiosInstance, AxiosError } from 'axios'
import { IApiClient } from '@/application/ports/api-client.port'
import { AuthResponse, LoginRequest, RegisterRequest } from '@/application/dto/auth.dto'
import { JobResponse, CreateJobRequest, UpdateJobRequest, JobListFilters } from '@/application/dto/job.dto'
import { StartupResponse, CreateStartupRequest, UpdateStartupRequest, StartupListFilters } from '@/application/dto/startup.dto'
import { User } from '@/domain/entities/user.entity'
import { ApiResponse, ApiError } from '@/domain/value-objects/api-response.vo'

export class ApiClient implements IApiClient {
  private client: AxiosInstance
  private baseURL: string

  constructor(baseURL?: string) {
    this.baseURL = baseURL || process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
    this.client = axios.create({
      baseURL: `${this.baseURL}/api/v1`,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    // Request interceptor to add auth token
    this.client.interceptors.request.use(
      (config) => {
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
    return this.request<JobResponse[]>({
      method: 'GET',
      url: '/jobs',
      params: filters,
    })
  }

  async getJob(id: string): Promise<ApiResponse<JobResponse>> {
    return this.request<JobResponse>({
      method: 'GET',
      url: `/jobs/${id}`,
    })
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
    return this.request<StartupResponse[]>({
      method: 'GET',
      url: '/startups',
      params: filters,
    })
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
}

// Singleton instance
export const apiClient = new ApiClient()

