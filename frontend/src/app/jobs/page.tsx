import type { Metadata } from 'next'
import Link from 'next/link'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { JobCard } from '@/presentation/components/job/job-card'
import { JobFilters } from '@/presentation/components/job/job-filters'
import { JobFiltersMobile } from '@/presentation/components/job/job-filters-mobile'
import { Pagination } from '@/presentation/components/ui/pagination'
import { apiClient } from '@/infrastructure/api/api-client'
import { JobSortControls } from '@/presentation/components/job/job-sort-controls'
import {
  buildPageMetadata,
  listingCanonicalPath,
  listingPageNumber,
  shouldIndexListing,
} from '@/lib/seo'

export const revalidate = 60

interface JobsPageProps {
  searchParams: Promise<{
    search?: string
    job_type?: string
    location_type?: string
    country?: string
    city?: string
    salary_min?: string
    salary_max?: string
    currency?: string
    order_by?: string
    order_dir?: 'ASC' | 'DESC'
    page?: string
  }>
}

export async function generateMetadata({ searchParams }: JobsPageProps): Promise<Metadata> {
  const params = await searchParams
  const page = listingPageNumber(params)
  const index = shouldIndexListing(params)
  const title =
    page > 1 ? `Startup Jobs – Page ${page}` : 'Startup Jobs – Remote, Hybrid & Onsite'

  return buildPageMetadata({
    title,
    description:
      'Browse tech and startup jobs across remote, hybrid, and onsite roles. Filter by type, location, and salary on JoinUs.',
    canonicalPath: listingCanonicalPath('/jobs', params),
    index,
  })
}

async function getJobs(filters: Record<string, string | undefined>) {
  try {
    const params: Record<string, string | number> = {
      page: filters.page ? parseInt(filters.page, 10) : 1,
      page_size: 12,
      status: 'active',
    }

    if (filters.search) params.search = filters.search
    if (filters.job_type) params.job_type = filters.job_type
    if (filters.location_type) params.location_type = filters.location_type
    if (filters.country) params.country = filters.country
    if (filters.city) params.city = filters.city
    if (filters.salary_min) params.salary_min = parseInt(filters.salary_min, 10)
    if (filters.salary_max) params.salary_max = parseInt(filters.salary_max, 10)
    if (filters.currency) params.currency = filters.currency
    if (filters.order_by) params.order_by = filters.order_by
    if (filters.order_dir) params.order_dir = filters.order_dir

    const response = await apiClient.listJobs(params as any)
    return {
      jobs: response.data || [],
      meta: response.meta,
    }
  } catch {
    return {
      jobs: [],
      meta: undefined,
    }
  }
}

const HUB_LINKS = [
  { href: '/jobs?location_type=remote', label: 'Remote' },
  { href: '/jobs?location_type=hybrid', label: 'Hybrid' },
  { href: '/jobs?job_type=full_time', label: 'Full-time' },
  { href: '/jobs?job_type=part_time', label: 'Part-time' },
  { href: '/jobs?job_type=contract', label: 'Contract' },
  { href: '/startups', label: 'Browse startups' },
]

export default async function JobsPage({ searchParams }: JobsPageProps) {
  const params = await searchParams
  const { jobs, meta } = await getJobs(params)

  const activeFiltersCount = [
    params.search,
    params.job_type,
    params.location_type,
    params.country,
    params.city,
    params.salary_min,
    params.salary_max,
    params.currency,
    params.order_by && params.order_by !== 'created_at' ? params.order_by : null,
    params.order_dir && params.order_dir !== 'DESC' ? params.order_dir : null,
  ].filter(Boolean).length

  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 py-12">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-secondary-900 mb-2">Job Opportunities</h1>
            <p className="text-secondary-600">
              Discover your next career move
              {activeFiltersCount > 0 && (
                <span className="ml-2 text-primary-600">
                  ({activeFiltersCount} filter{activeFiltersCount !== 1 ? 's' : ''} active)
                </span>
              )}
            </p>
            <ul className="mt-4 flex flex-wrap gap-x-4 gap-y-2 text-sm">
              {HUB_LINKS.map((link) => (
                <li key={link.href}>
                  <Link href={link.href} className="text-primary-600 hover:text-primary-700">
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div className="flex flex-col lg:flex-row gap-8">
            <aside className="lg:w-80 flex-shrink-0">
              <JobFiltersMobile />
              <div className="hidden md:block">
                <JobFilters />
              </div>
            </aside>

            <div className="flex-1">
              <div className="mb-6">
                <JobSortControls />
              </div>

              {jobs.length > 0 ? (
                <>
                  <div className="mb-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 xl:grid-cols-3 gap-6">
                      {jobs.map((job) => (
                        <JobCard key={job.id} job={job} />
                      ))}
                    </div>
                  </div>

                  {meta && <Pagination meta={meta} />}
                </>
              ) : (
                <div className="text-center py-12">
                  <p className="text-secondary-600 text-lg mb-2">No jobs found</p>
                  <p className="text-secondary-500 text-sm">
                    Try adjusting your filters or search terms
                  </p>
                </div>
              )}
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  )
}
