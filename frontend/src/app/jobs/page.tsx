import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { JobCard } from '@/presentation/components/job/job-card'
import { JobFilters } from '@/presentation/components/job/job-filters'
import { JobFiltersMobile } from '@/presentation/components/job/job-filters-mobile'
import { apiClient } from '@/infrastructure/api/api-client'

interface JobsPageProps {
  searchParams: {
    search?: string
    job_type?: string
    location_type?: string
    page?: string
  }
}

async function getJobs(filters: any) {
  try {
    const response = await apiClient.listJobs({
      ...filters,
      page: filters.page ? parseInt(filters.page) : 1,
      page_size: 12,
    })
    return response.data || []
  } catch {
    return []
  }
}

export default async function JobsPage({ searchParams }: JobsPageProps) {
  const jobs = await getJobs(searchParams)
  
  const activeFiltersCount = [
    searchParams.search,
    searchParams.job_type,
    searchParams.location_type,
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
              {jobs.length > 0 ? (
                <>
                  <div className="mb-4 text-sm text-secondary-600">
                    {jobs.length} job{jobs.length !== 1 ? 's' : ''} found
                  </div>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 xl:grid-cols-3 gap-6">
                    {jobs.map((job) => (
                      <JobCard key={job.id} job={job} />
                    ))}
                  </div>
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

