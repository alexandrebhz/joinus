import axios, { AxiosInstance, AxiosError } from 'axios'

export interface PaginationConfig {
  type: 'query_param' | 'url_pattern' | 'link_follow' | 'api_pagination'
  param_name?: string
  start_page?: number
  increment?: number
  max_pages?: number
  next_page_selector?: string
  url_pattern?: string
  api_config?: {
    endpoint: string
    page_param: string
    page_size: number
    max_pages: number
  }
}

export interface JobURLRule {
  type: 'relative' | 'absolute' | 'attribute'
  selector: string
  attribute?: string
  base_url?: string
}

export interface FieldRule {
  selector: string
  type: 'text' | 'html' | 'attribute' | 'regex'
  attribute?: string
  regex_pattern?: string
  required: boolean
  default_value?: string
  transformations?: string[]
}

export interface ExtractionRules {
  job_list_selector: string
  job_detail_url: JobURLRule
  fields: Record<string, FieldRule>
}

export interface CrawlSite {
  id: string
  name: string
  base_url: string
  backend_startup_id: string
  active: boolean
  schedule: string
  last_crawled_at?: string
  next_crawl_at?: string
  crawl_interval: 'daily' | 'weekly' | 'custom'
  pagination_config: PaginationConfig
  extraction_rules: ExtractionRules
  deduplication_key: 'url' | 'composite' | 'external_id'
  request_delay: number
  user_agent: string
  created_at: string
  updated_at: string
}

export interface CreateSiteInput {
  name: string
  base_url: string
  backend_startup_id: string
  schedule: string
  crawl_interval: 'daily' | 'weekly' | 'custom'
  pagination_config: PaginationConfig
  extraction_rules: ExtractionRules
  deduplication_key?: 'url' | 'composite' | 'external_id'
  request_delay?: number
  user_agent?: string
}

export interface UpdateSiteInput {
  name?: string
  base_url?: string
  backend_startup_id?: string
  active?: boolean
  schedule?: string
  crawl_interval?: 'daily' | 'weekly' | 'custom'
  pagination_config?: PaginationConfig
  extraction_rules?: ExtractionRules
  deduplication_key?: 'url' | 'composite' | 'external_id'
  request_delay?: number
  user_agent?: string
}

export interface ApiResponse<T> {
  data: T
  error?: string
}

export interface CrawlResult {
  jobs_found: number
  jobs_saved: number
  jobs_skipped: number
  pages_crawled?: number
  errors?: string[]
}

export interface LogEntry {
  timestamp: string
  level: 'info' | 'warning' | 'error'
  message: string
}

export interface CrawlLog {
  id: string
  site_id: string
  status: 'running' | 'completed' | 'failed'
  started_at: string
  completed_at?: string
  duration_ms: number
  jobs_found: number
  jobs_saved: number
  jobs_skipped: number
  pages_crawled: number
  errors: string[]
  logs: LogEntry[]
  created_at: string
}

class CrawlerApiClient {
  private client: AxiosInstance

  constructor() {
    const baseURL = process.env.NEXT_PUBLIC_CRAWLER_API_URL || 'http://localhost:8081'
    this.client = axios.create({
      baseURL: `${baseURL}/api/v1`,
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  async listSites(): Promise<ApiResponse<CrawlSite[]>> {
    const response = await this.client.get<ApiResponse<CrawlSite[]>>('/sites')
    return response.data
  }

  async getSite(id: string): Promise<ApiResponse<CrawlSite>> {
    const response = await this.client.get<ApiResponse<CrawlSite>>(`/sites/${id}`)
    return response.data
  }

  async createSite(data: CreateSiteInput): Promise<ApiResponse<CrawlSite>> {
    const response = await this.client.post<ApiResponse<CrawlSite>>('/sites', data)
    return response.data
  }

  async updateSite(id: string, data: UpdateSiteInput): Promise<ApiResponse<CrawlSite>> {
    const response = await this.client.put<ApiResponse<CrawlSite>>(`/sites/${id}`, data)
    return response.data
  }

  async deleteSite(id: string): Promise<void> {
    await this.client.delete(`/sites/${id}`)
  }

  async executeCrawl(id: string): Promise<ApiResponse<CrawlResult>> {
    const response = await this.client.post<ApiResponse<CrawlResult>>(`/sites/${id}/crawl`)
    return response.data
  }

  async getCrawlLogs(siteId: string, limit?: number): Promise<ApiResponse<CrawlLog[]>> {
    const params = limit ? `?limit=${limit}` : ''
    const response = await this.client.get<ApiResponse<CrawlLog[]>>(`/sites/${siteId}/logs${params}`)
    return response.data
  }

  async getLatestCrawlLog(siteId: string): Promise<ApiResponse<CrawlLog>> {
    const response = await this.client.get<ApiResponse<CrawlLog>>(`/sites/${siteId}/logs/latest`)
    return response.data
  }
}

export const crawlerApiClient = new CrawlerApiClient()

