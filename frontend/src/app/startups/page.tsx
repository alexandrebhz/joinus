import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { StartupCard } from '@/presentation/components/startup/startup-card'
import { StartupFilters } from '@/presentation/components/startup/startup-filters'
import { apiClient } from '@/infrastructure/api/api-client'
import { StartupStatus } from '@/domain/entities/startup.entity'

interface StartupsPageProps {
  searchParams: {
    search?: string
    industry?: string
    status?: StartupStatus
    location?: string
    company_size?: string
    page?: string
  }
}

async function getStartups(filters: any) {
  try {
    const response = await apiClient.listStartups({
      ...filters,
      page: filters.page ? parseInt(filters.page) : 1,
      page_size: 12,
    })
    return response.data || []
  } catch {
    return []
  }
}

export default async function StartupsPage({ searchParams }: StartupsPageProps) {
  const startups = await getStartups(searchParams)

  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 py-12">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-secondary-900 mb-2">Startups</h1>
            <p className="text-secondary-600">Explore innovative companies</p>
          </div>

          {/* Filters and Results Layout */}
          <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
            {/* Filters Sidebar */}
            <div className="lg:col-span-1">
              <StartupFilters />
            </div>

            {/* Startups Grid */}
            <div className="lg:col-span-3">
              {startups.length > 0 ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {startups.map((startup) => (
                    <StartupCard key={startup.id} startup={startup} />
                  ))}
                </div>
              ) : (
                <div className="text-center py-12">
                  <p className="text-secondary-600 text-lg">No startups found. Try adjusting your filters.</p>
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

