import { MetadataRoute } from 'next'
import { apiClient } from '@/infrastructure/api/api-client'
import type { PaginationMeta } from '@/domain/value-objects/api-response.vo'
import type { Startup } from '@/domain/entities/startup.entity'
import { getSiteUrl } from '@/lib/seo'

const baseUrl = getSiteUrl()

/** Avoid runaway loops if the API omits pagination meta. */
const MAX_SITEMAP_PAGES = 500
const SITEMAP_PAGE_SIZE = 100
/** Listing page size used by /jobs and /startups UI (for pagination URLs). */
const LISTING_PAGE_SIZE = 12

async function fetchAllPages<T>(
  load: (page: number) => Promise<{ data?: T[]; meta?: PaginationMeta }>,
): Promise<{ items: T[]; totalCount: number }> {
  const out: T[] = []
  let totalCount = 0

  for (let page = 1; page <= MAX_SITEMAP_PAGES; page++) {
    const { data, meta } = await load(page)
    const batch = data ?? []
    out.push(...batch)
    if (meta?.totalCount) totalCount = meta.totalCount
    if (batch.length === 0) break
    if (meta && page >= meta.totalPages) break
    if (batch.length < SITEMAP_PAGE_SIZE) break
  }

  if (!totalCount) totalCount = out.length
  return { items: out, totalCount }
}

async function fetchJobsForSitemap() {
  try {
    return await fetchAllPages((page) =>
      apiClient.listJobs({
        page,
        page_size: SITEMAP_PAGE_SIZE,
        status: 'active',
      }),
    )
  } catch (error) {
    console.error('Error fetching jobs for sitemap:', error)
    return { items: [], totalCount: 0 }
  }
}

async function fetchStartupsForSitemap() {
  try {
    return await fetchAllPages((page) =>
      apiClient.listStartups({
        page,
        page_size: SITEMAP_PAGE_SIZE,
        status: 'active',
      }),
    )
  } catch (error) {
    console.error('Error fetching startups for sitemap:', error)
    return { items: [], totalCount: 0 }
  }
}

/** API list payloads use snake_case dates; entity types use camelCase. */
function startupLastModified(startup: Startup): Date {
  const s = startup as Startup & { updated_at?: string; created_at?: string }
  const iso = s.updatedAt ?? s.updated_at ?? s.createdAt ?? s.created_at
  return new Date(iso || Date.now())
}

function listingPageEntries(
  path: '/jobs' | '/startups',
  totalCount: number,
): MetadataRoute.Sitemap {
  const totalPages = Math.max(1, Math.ceil(totalCount / LISTING_PAGE_SIZE))
  const entries: MetadataRoute.Sitemap = []

  // Page 1 is already covered by the static /jobs or /startups entry.
  for (let page = 2; page <= totalPages && page <= 50; page++) {
    entries.push({
      url: `${baseUrl}${path}?page=${page}`,
      lastModified: new Date(),
      changeFrequency: 'daily',
      priority: 0.6,
    })
  }

  return entries
}

/** Build sitemap at request time so new startups/jobs appear without a rebuild. */
export const dynamic = 'force-dynamic'
export const revalidate = 300

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
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

  const [jobsResult, startupsResult] = await Promise.all([
    fetchJobsForSitemap(),
    fetchStartupsForSitemap(),
  ])

  const jobPages: MetadataRoute.Sitemap = jobsResult.items.map((job) => ({
    url: `${baseUrl}/jobs/${job.id}`,
    lastModified: job.updatedAt
      ? new Date(job.updatedAt)
      : new Date(job.createdAt || Date.now()),
    changeFrequency: 'weekly' as const,
    priority: 0.8,
  }))

  const startupPages: MetadataRoute.Sitemap = startupsResult.items
    .filter((startup) => startup.slug)
    .map((startup) => ({
      url: `${baseUrl}/startups/${startup.slug}`,
      lastModified: startupLastModified(startup),
      changeFrequency: 'weekly' as const,
      priority: 0.8,
    }))

  const jobListingPages = listingPageEntries('/jobs', jobsResult.totalCount)
  const startupListingPages = listingPageEntries('/startups', startupsResult.totalCount)

  return [
    ...staticPages,
    ...jobListingPages,
    ...startupListingPages,
    ...jobPages,
    ...startupPages,
  ]
}
