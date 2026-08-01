import type { Metadata } from 'next'

const DEFAULT_SITE_URL = 'https://joinus.ie'
/** Google typically displays ~50–60 characters; keep a hard cap for Ahrefs. */
export const MAX_TITLE_LENGTH = 60

export function getSiteUrl(): string {
  const raw =
    process.env.NEXT_PUBLIC_SITE_URL ||
    process.env.NEXT_PUBLIC_APP_URL ||
    DEFAULT_SITE_URL
  return raw.replace(/\/+$/, '')
}

export function absoluteUrl(path = '/'): string {
  const base = getSiteUrl()
  if (!path || path === '/') return base
  return `${base}${path.startsWith('/') ? path : `/${path}`}`
}

/** Truncate a title near a word boundary without exceeding maxLen. */
export function truncateTitle(title: string, maxLen = MAX_TITLE_LENGTH): string {
  const normalized = title.replace(/\s+/g, ' ').trim()
  if (normalized.length <= maxLen) return normalized

  const sliced = normalized.slice(0, maxLen - 1)
  const lastSpace = sliced.lastIndexOf(' ')
  const base = lastSpace > maxLen * 0.5 ? sliced.slice(0, lastSpace) : sliced
  return `${base.trimEnd()}…`
}

/**
 * Query params that create alternate views of listing pages.
 * When any of these are present (except a clean page>1), we noindex and
 * canonical to the clean listing URL to avoid duplicate-content issues.
 */
const LISTING_FILTER_PARAMS = new Set([
  'search',
  'job_type',
  'location_type',
  'country',
  'city',
  'salary_min',
  'salary_max',
  'currency',
  'order_by',
  'order_dir',
  'industry',
  'status',
  'location',
  'company_size',
])

export function hasListingFilters(
  searchParams: Record<string, string | string[] | undefined>,
): boolean {
  return Object.keys(searchParams).some((key) => LISTING_FILTER_PARAMS.has(key))
}

export function listingPageNumber(
  searchParams: Record<string, string | string[] | undefined>,
): number {
  const raw = searchParams.page
  const value = Array.isArray(raw) ? raw[0] : raw
  const page = value ? parseInt(value, 10) : 1
  return Number.isFinite(page) && page > 1 ? page : 1
}

/** Canonical path for /jobs or /startups: clean path, or ?page=N when paginated. */
export function listingCanonicalPath(
  basePath: '/jobs' | '/startups',
  searchParams: Record<string, string | string[] | undefined>,
): string {
  if (hasListingFilters(searchParams)) return basePath
  const page = listingPageNumber(searchParams)
  return page > 1 ? `${basePath}?page=${page}` : basePath
}

export function shouldIndexListing(
  searchParams: Record<string, string | string[] | undefined>,
): boolean {
  return !hasListingFilters(searchParams)
}

type BuildMetadataOptions = {
  title: string
  description?: string
  /** Path or path+query for the preferred URL (e.g. `/jobs` or `/jobs?page=2`). */
  canonicalPath: string
  /** When false, emits noindex,follow. Default true. */
  index?: boolean
  /** Use absolute title (skip layout `%s | JoinUs` template). Default true. */
  absoluteTitle?: boolean
  openGraph?: Metadata['openGraph']
}

export function buildPageMetadata({
  title,
  description,
  canonicalPath,
  index = true,
  absoluteTitle = true,
  openGraph,
}: BuildMetadataOptions): Metadata {
  const shortTitle = truncateTitle(title)
  const canonical = absoluteUrl(canonicalPath)

  return {
    title: absoluteTitle ? { absolute: shortTitle } : shortTitle,
    ...(description ? { description } : {}),
    alternates: {
      canonical,
    },
    robots: {
      index,
      follow: true,
      googleBot: {
        index,
        follow: true,
      },
    },
    openGraph: {
      title: shortTitle,
      ...(description ? { description } : {}),
      url: canonical,
      ...openGraph,
    },
  }
}
