import { MetadataRoute } from 'next'
import { apiClient } from '@/infrastructure/api/api-client'
import type { PaginationMeta } from '@/domain/value-objects/api-response.vo'
import type { Startup } from '@/domain/entities/startup.entity'

const baseUrl = process.env.NEXT_PUBLIC_SITE_URL || process.env.NEXT_PUBLIC_APP_URL || 'https://joinus.ie'

/** Avoid runaway loops if the API omits pagination meta. */
const MAX_SITEMAP_PAGES = 500
const SITEMAP_PAGE_SIZE = 100

async function fetchAllPages<T>(
  load: (page: number) => Promise<{ data?: T[]; meta?: PaginationMeta }>,
): Promise<T[]> {
  const out: T[] = []
  for (let page = 1; page <= MAX_SITEMAP_PAGES; page++) {
    const { data, meta } = await load(page)
    const batch = data ?? []
    out.push(...batch)
    if (batch.length === 0) break
    if (meta && page >= meta.totalPages) break
    if (batch.length < SITEMAP_PAGE_SIZE) break
  }
  return out
}

async function fetchJobsForSitemap() {
  try {
    return await fetchAllPages((page) =>
      apiClient.listJobs({
        page,
        page_size: SITEMAP_PAGE_SIZE,
      }),
    )
  } catch (error) {
    console.error('Error fetching jobs for sitemap:', error)
    return []
  }
}

async function fetchStartupsForSitemap() {
  try {
    return await fetchAllPages((page) =>
      apiClient.listStartups({
        page,
        page_size: SITEMAP_PAGE_SIZE,
      }),
    )
  } catch (error) {
    console.error('Error fetching startups for sitemap:', error)
    return []
  }
}

/** API list payloads use snake_case dates; entity types use camelCase. */
function startupLastModified(startup: Startup): Date {
  const s = startup as Startup & { updated_at?: string; created_at?: string }
  const iso = s.updatedAt ?? s.updated_at ?? s.createdAt ?? s.created_at
  return new Date(iso || Date.now())
}

/** Build sitemap at request time so new startups/jobs appear without a rebuild. */
export const dynamic = 'force-dynamic'

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  // Static pages
  const staticPages: MetadataRoute.Sitemap = [
    {
      url: baseUrl,
      lastModified: new Date(),
      changeFrequency: 'daily',
      priority: 1,
    },
    {
      url: `${baseUrl}/jobs`,
      lastModified: new Date(),
      changeFrequency: 'hourly',
      priority: 0.9,
    },
    {
      url: `${baseUrl}/startups`,
      lastModified: new Date(),
      changeFrequency: 'hourly',
      priority: 0.9,
    },
    {
      url: `${baseUrl}/about`,
      lastModified: new Date(),
      changeFrequency: 'monthly',
      priority: 0.7,
    },
    {
      url: `${baseUrl}/contact`,
      lastModified: new Date(),
      changeFrequency: 'monthly',
      priority: 0.7,
    },
  ]

  const [jobs, startups] = await Promise.all([
    fetchJobsForSitemap(),
    fetchStartupsForSitemap(),
  ])

  const jobPages: MetadataRoute.Sitemap = jobs.map((job) => ({
    url: `${baseUrl}/jobs/${job.id}`,
    lastModified: job.updatedAt ? new Date(job.updatedAt) : new Date(job.createdAt || Date.now()),
    changeFrequency: 'weekly' as const,
    priority: 0.8,
  }))

  const startupPages: MetadataRoute.Sitemap = startups
    .filter((startup) => startup.slug)
    .map((startup) => ({
      url: `${baseUrl}/startups/${startup.slug}`,
      lastModified: startupLastModified(startup),
      changeFrequency: 'weekly' as const,
      priority: 0.8,
    }))

  return [...staticPages, ...jobPages, ...startupPages]
}
