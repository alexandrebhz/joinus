import { Suspense } from 'react'
import type { Metadata } from 'next'
import Link from 'next/link'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { StartupCard } from '@/presentation/components/startup/startup-card'
import { StartupFilters } from '@/presentation/components/startup/startup-filters'
import { Pagination } from '@/presentation/components/ui/pagination'
import { apiClient } from '@/infrastructure/api/api-client'
import { StartupStatus } from '@/domain/entities/startup.entity'
import {
  buildPageMetadata,
  listingCanonicalPath,
  listingPageNumber,
  shouldIndexListing,
} from '@/lib/seo'

export const revalidate = 60

interface StartupsPageProps {
  searchParams: Promise<{
    search?: string
    industry?: string
    status?: StartupStatus
    location?: string
    company_size?: string
    page?: string
  }>
}

export async function generateMetadata({ searchParams }: StartupsPageProps): Promise<Metadata> {
  const params = await searchParams
  const page = listingPageNumber(params)
  const index = shouldIndexListing(params)
  const title = page > 1 ? `Startups Hiring – Page ${page}` : 'Startups Hiring on JoinUs'

  return buildPageMetadata({
    title,
    description:
      'Explore innovative startups hiring on JoinUs. Browse companies by industry, size, and location.',
    canonicalPath: listingCanonicalPath('/startups', params),
    index,
  })
}

async function getStartups(filters: Record<string, string | undefined>) {
  try {
    const response = await apiClient.listStartups({
      ...filters,
      status: (filters.status as StartupStatus) || 'active',
      page: filters.page ? parseInt(filters.page, 10) : 1,
      page_size: 12,
    })
    return {
      startups: response.data || [],
      meta: response.meta,
    }
  } catch {
    return {
      startups: [],
      meta: undefined,
    }
  }
}

const HUB_LINKS = [
  { href: '/jobs', label: 'All jobs' },
  { href: '/jobs?location_type=remote', label: 'Remote jobs' },
  { href: '/about', label: 'About JoinUs' },
]

export default async function StartupsPage({ searchParams }: StartupsPageProps) {
  const params = await searchParams
  const { startups, meta } = await getStartups(params)

  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 py-12">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-secondary-900 mb-2">Startups</h1>
            <p className="text-secondary-600">Explore innovative companies</p>
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

          <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
            <div className="lg:col-span-1">
              <Suspense
                fallback={
                  <div className="border border-secondary-200 rounded-lg p-6 bg-white">
                    <div className="animate-pulse">
                      <div className="h-6 bg-secondary-200 rounded w-1/3 mb-4"></div>
                      <div className="space-y-4">
                        <div className="h-10 bg-secondary-100 rounded"></div>
                        <div className="h-10 bg-secondary-100 rounded"></div>
                        <div className="h-10 bg-secondary-100 rounded"></div>
                      </div>
                    </div>
                  </div>
                }
              >
                <StartupFilters />
              </Suspense>
            </div>

            <div className="lg:col-span-3">
              {startups.length > 0 ? (
                <>
                  <div className="mb-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                      {startups.map((startup) => (
                        <StartupCard key={startup.id} startup={startup} />
                      ))}
                    </div>
                  </div>

                  {meta && <Pagination meta={meta} basePath="/startups" />}
                </>
              ) : (
                <div className="text-center py-12">
                  <p className="text-secondary-600 text-lg mb-2">No startups found</p>
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
