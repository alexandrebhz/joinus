import { notFound } from 'next/navigation'
import type { Metadata } from 'next'
import Image from 'next/image'
import Link from 'next/link'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { apiClient } from '@/infrastructure/api/api-client'
import { JobCard } from '@/presentation/components/job/job-card'
import { Button } from '@/presentation/components/ui/button'
import { Card, CardContent } from '@/presentation/components/ui/card'
import { Building2, MapPin, Calendar, Users, Briefcase, ExternalLink } from 'lucide-react'
import { buildPageMetadata, truncateTitle } from '@/lib/seo'

export const revalidate = 60

interface StartupDetailPageProps {
  params: Promise<{
    slug: string
  }>
}

async function getStartup(slug: string) {
  try {
    const response = await apiClient.getStartupBySlug(slug)
    return response.data
  } catch {
    return null
  }
}

async function getStartupJobs(startupId: string) {
  try {
    const response = await apiClient.listJobs({
      startup_id: startupId,
      status: 'active',
      page_size: 50,
    })
    return response.data || []
  } catch {
    return []
  }
}

export async function generateMetadata({ params }: StartupDetailPageProps): Promise<Metadata> {
  const { slug } = await params
  const startup = await getStartup(slug)

  if (!startup) {
    return buildPageMetadata({
      title: 'Startup Not Found',
      canonicalPath: `/startups/${slug}`,
      index: false,
    })
  }

  return buildPageMetadata({
    title: truncateTitle(`${startup.name} – Jobs & Company`),
    description:
      startup.description?.substring(0, 160) ||
      `Explore ${startup.name} and open roles on JoinUs.`,
    canonicalPath: `/startups/${startup.slug}`,
  })
}

export default async function StartupDetailPage({ params }: StartupDetailPageProps) {
  const { slug } = await params
  const startup = await getStartup(slug)

  if (!startup) {
    notFound()
  }

  const jobs = await getStartupJobs(startup.id)

  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1">
        <section className="bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-16">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <nav aria-label="Breadcrumb" className="mb-6 text-sm text-secondary-600">
              <ol className="flex flex-wrap items-center gap-2">
                <li>
                  <Link href="/startups" className="text-primary-600 hover:text-primary-700">
                    Startups
                  </Link>
                </li>
                <li aria-hidden="true">/</li>
                <li className="text-secondary-900 font-medium">{startup.name}</li>
              </ol>
            </nav>

            <div className="flex flex-col md:flex-row items-start md:items-center gap-8">
              {startup.logoUrl ? (
                <div className="relative h-24 w-24 rounded-xl overflow-hidden border-2 border-white shadow-lg flex-shrink-0">
                  <Image
                    src={startup.logoUrl}
                    alt={startup.name}
                    fill
                    className="object-cover"
                  />
                </div>
              ) : (
                <div className="h-24 w-24 rounded-xl bg-primary-100 flex items-center justify-center flex-shrink-0">
                  <Building2 className="h-12 w-12 text-primary-600" />
                </div>
              )}
              <div className="flex-1">
                <h1 className="text-4xl font-bold text-secondary-900 mb-2">{startup.name}</h1>
                <p className="text-lg text-secondary-600 mb-4">{startup.description}</p>
                <div className="flex flex-wrap gap-4 text-sm text-secondary-600">
                  <div className="flex items-center">
                    <MapPin className="h-4 w-4 mr-1" />
                    {startup.location}
                  </div>
                  <div className="flex items-center">
                    <Building2 className="h-4 w-4 mr-1" />
                    {startup.industry}
                  </div>
                  <div className="flex items-center">
                    <Calendar className="h-4 w-4 mr-1" />
                    Founded {startup.foundedYear}
                  </div>
                  {startup.memberCount !== undefined && (
                    <div className="flex items-center">
                      <Users className="h-4 w-4 mr-1" />
                      {startup.memberCount} members
                    </div>
                  )}
                </div>
              </div>
              {startup.website && (
                <a href={startup.website} target="_blank" rel="noopener noreferrer">
                  <Button variant="outline">
                    <ExternalLink className="h-4 w-4 mr-2" />
                    Visit Website
                  </Button>
                </a>
              )}
            </div>
          </div>
        </section>

        {jobs.length > 0 && (
          <section className="py-12">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8">
              <div className="mb-8 flex items-end justify-between gap-4">
                <div>
                  <h2 className="text-3xl font-bold text-secondary-900 mb-2">Open Positions</h2>
                  <p className="text-secondary-600">
                    {jobs.length} job{jobs.length !== 1 ? 's' : ''} available
                  </p>
                </div>
                <Link href="/jobs" className="text-sm text-primary-600 hover:text-primary-700">
                  Browse all jobs
                </Link>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {jobs.map((job) => (
                  <JobCard key={job.id} job={job} />
                ))}
              </div>
            </div>
          </section>
        )}

        {jobs.length === 0 && (
          <section className="py-12">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8">
              <Card>
                <CardContent className="pt-6 text-center py-12">
                  <Briefcase className="h-12 w-12 text-secondary-400 mx-auto mb-4" />
                  <h3 className="text-xl font-semibold text-secondary-900 mb-2">No open positions</h3>
                  <p className="text-secondary-600 mb-4">
                    This startup doesn&apos;t have any open positions at the moment.
                  </p>
                  <Link href="/jobs" className="text-primary-600 hover:text-primary-700">
                    Browse other jobs
                  </Link>
                </CardContent>
              </Card>
            </div>
          </section>
        )}
      </main>
      <Footer />
    </div>
  )
}
