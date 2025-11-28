import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { JobCard } from '@/presentation/components/job/job-card'
import { JobFilters } from '@/presentation/components/job/job-filters'
import { JobFiltersMobile } from '@/presentation/components/job/job-filters-mobile'
import { Pagination } from '@/presentation/components/ui/pagination'
import { apiClient } from '@/infrastructure/api/api-client'
import { JobSortControls } from '@/presentation/components/job/job-sort-controls'

interface JobsPageProps {
  searchParams: {
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
  }
}

async function getJobs(filters: any) {
  try {
    const params: any = {
      page: filters.page ? parseInt(filters.page) : 1,
      page_size: 12,
    }
    
    if (filters.search) params.search = filters.search
    if (filters.job_type) params.job_type = filters.job_type
    if (filters.location_type) params.location_type = filters.location_type
    if (filters.country) params.country = filters.country
    if (filters.city) params.city = filters.city
    if (filters.salary_min) params.salary_min = parseInt(filters.salary_min)
    if (filters.salary_max) params.salary_max = parseInt(filters.salary_max)
    if (filters.currency) params.currency = filters.currency
    if (filters.order_by) params.order_by = filters.order_by
    if (filters.order_dir) params.order_dir = filters.order_dir
    
    const response = await apiClient.listJobs(params)
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

export default async function JobsPage({ searchParams }: JobsPageProps) {
  const { jobs, meta } = await getJobs(searchParams)
  
  const activeFiltersCount = [
    searchParams.search,
    searchParams.job_type,
    searchParams.location_type,
    searchParams.country,
    searchParams.city,
    searchParams.salary_min,
    searchParams.salary_max,
    searchParams.currency,
    searchParams.order_by && searchParams.order_by !== 'created_at' ? searchParams.order_by : null,
    searchParams.order_dir && searchParams.order_dir !== 'DESC' ? searchParams.order_dir : null,
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
          </div>

          <div className="flex flex-col lg:flex-row gap-8">
            {/* Filters Sidebar */}
            <aside className="lg:w-80 flex-shrink-0">
              {/* Mobile Filters Toggle */}
              <JobFiltersMobile />
              
              {/* Desktop Filters */}
              <div className="hidden md:block">
                <JobFilters />
              </div>
            </aside>

            {/* Jobs List */}
            <div className="flex-1">
              {/* Sort Controls */}
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

                  {/* Pagination */}
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

