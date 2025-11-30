import { SiteRepository } from '../../domain/repositories/SiteRepository';
import { CrawlRepository } from '../../domain/repositories/CrawlRepository';
import { CrawlSite } from '../../domain/entities/CrawlSite';
import { CreateSiteInput } from '../../domain/value-objects/CreateSiteInput';
import { UpdateSiteInput } from '../../domain/value-objects/UpdateSiteInput';
import { CrawlResult } from '../../domain/entities/CrawlResult';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081';

export class ApiClient implements SiteRepository, CrawlRepository {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ 
        error: response.statusText 
      }));
      throw new Error(error.error || `HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data.data || data;
  }

  // SiteRepository implementation
  async findAll(): Promise<CrawlSite[]> {
    return this.request<CrawlSite[]>('/api/v1/sites');
  }

  async findById(id: string): Promise<CrawlSite> {
    return this.request<CrawlSite>(`/api/v1/sites/${id}`);
  }

  async create(input: CreateSiteInput): Promise<CrawlSite> {
    return this.request<CrawlSite>('/api/v1/sites', {
      method: 'POST',
      body: JSON.stringify(input),
    });
  }

  async update(id: string, input: UpdateSiteInput): Promise<CrawlSite> {
    return this.request<CrawlSite>(`/api/v1/sites/${id}`, {
      method: 'PUT',
      body: JSON.stringify(input),
    });
  }

  async delete(id: string): Promise<void> {
    return this.request<void>(`/api/v1/sites/${id}`, {
      method: 'DELETE',
    });
  }

  // CrawlRepository implementation
  async execute(siteId: string): Promise<CrawlResult> {
    return this.request<CrawlResult>(`/api/v1/sites/${siteId}/crawl`, {
      method: 'POST',
    });
  }
}

