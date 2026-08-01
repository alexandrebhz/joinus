import { notFound } from 'next/navigation'
import { Metadata } from 'next'
import Link from 'next/link'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { apiClient } from '@/infrastructure/api/api-client'
import { JobCard } from '@/presentation/components/job/job-card'
import { Button } from '@/presentation/components/ui/button'
import { Card, CardContent } from '@/presentation/components/ui/card'
import { MapPin, Briefcase, DollarSign, Calendar, ExternalLink, Mail } from 'lucide-react'
import { formatDistanceToNow } from 'date-fns'
import { buildPageMetadata, truncateTitle } from '@/lib/seo'

export const revalidate = 60

interface JobDetailPageProps {
  params: Promise<{
    id: string
  }>
}

async function getJob(id: string) {
  try {
    const response = await apiClient.getJob(id)
    return response.data
  } catch (error: unknown) {
    console.error('Error fetching job:', error)
    return null
  }
}

async function getRelatedJobs(job: {
  id: string
  startupId: string
  jobType?: string
  locationType?: string
}) {
  try {
    const [sameStartup, similar] = await Promise.all([
      apiClient.listJobs({
        startup_id: job.startupId,
        status: 'active',
        page: 1,
        page_size: 6,
      }),
      apiClient.listJobs({
        job_type: job.jobType as any,
        location_type: job.locationType as any,
        status: 'active',
        page: 1,
        page_size: 6,
      }),
    ])

    const seen = new Set<string>([job.id])
    const related = []

    for (const candidate of [...(sameStartup.data || []), ...(similar.data || [])]) {
      if (seen.has(candidate.id)) continue
      seen.add(candidate.id)
      related.push(candidate)
      if (related.length >= 3) break
    }

    return related
  } catch {
    return []
  }
}

export async function generateMetadata({ params }: JobDetailPageProps): Promise<Metadata> {
  const { id } = await params
  const job = await getJob(id)

  if (!job) {
    return buildPageMetadata({
      title: 'Job Not Found',
      canonicalPath: `/jobs/${id}`,
      index: false,
    })
  }

  const company = job.startupName || 'Startup'
  const rawTitle = `${job.title} at ${company}`
  const title = truncateTitle(rawTitle, 60)

  return buildPageMetadata({
    title,
    description:
      job.description?.substring(0, 160) || `Job opening at ${company} on JoinUs`,
    canonicalPath: `/jobs/${job.id}`,
  })
}

export default async function JobDetailPage({ params }: JobDetailPageProps) {
  const { id } = await params
  const job = await getJob(id)

  if (!job) {
    notFound()
  }

  const relatedJobs = await getRelatedJobs(job)
  const startupHref = job.startupSlug ? `/startups/${job.startupSlug}` : null

  const formatSalary = () => {
    if (!job.salaryMin && !job.salaryMax) return null
    const min = job.salaryMin?.toLocaleString() || 'N/A'
    const max = job.salaryMax?.toLocaleString() || 'N/A'
    return `${job.currency} ${min} - ${max}`
  }

  const getLocationText = () => {
    if (job.locationType === 'remote') return 'Remote'
    if (job.locationType === 'hybrid') return `Hybrid • ${job.city || job.country}`
    return `${job.city ? `${job.city}, ` : ''}${job.country}`
  }

  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 py-12">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-4xl">
          <nav aria-label="Breadcrumb" className="mb-6 text-sm text-secondary-600">
            <ol className="flex flex-wrap items-center gap-2">
              <li>
                <Link href="/jobs" className="text-primary-600 hover:text-primary-700">
                  Jobs
                </Link>
              </li>
              <li aria-hidden="true">/</li>
              {startupHref && job.startupName ? (
                <>
                  <li>
                    <Link href={startupHref} className="text-primary-600 hover:text-primary-700">
                      {job.startupName}
                    </Link>
                  </li>
                  <li aria-hidden="true">/</li>
                </>
              ) : null}
              <li className="text-secondary-900 font-medium line-clamp-1">{job.title}</li>
            </ol>
          </nav>

          <div className="mb-8">
            <h1 className="text-4xl font-bold text-secondary-900 mb-2">{job.title}</h1>
            {job.startupName && startupHref && (
              <Link href={startupHref}>
                <p className="text-lg text-primary-600 hover:text-primary-700 font-medium transition-colors">
                  {job.startupName}
                </p>
              </Link>
            )}
            {job.startupName && !startupHref && (
              <p className="text-lg text-secondary-700 font-medium">{job.startupName}</p>
            )}
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2 space-y-6">
              <Card>
                <CardContent className="pt-6">
                  <h2 className="text-2xl font-semibold mb-4">Job Description</h2>
                  <div className="prose max-w-none">
                    <p className="whitespace-pre-wrap text-secondary-700">{job.description}</p>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardContent className="pt-6">
                  <h2 className="text-2xl font-semibold mb-4">Requirements</h2>
                  <div className="prose max-w-none">
                    <p className="whitespace-pre-wrap text-secondary-700">{job.requirements}</p>
                  </div>
                </CardContent>
              </Card>
            </div>

            <div className="space-y-6">
              <Card>
                <CardContent className="pt-6">
                  <div className="space-y-4">
                    <div className="flex items-center text-secondary-700">
                      <MapPin className="h-5 w-5 mr-3 text-secondary-400" />
                      <span>{getLocationText()}</span>
                    </div>
                    <div className="flex items-center text-secondary-700">
                      <Briefcase className="h-5 w-5 mr-3 text-secondary-400" />
                      <span className="capitalize">{job.jobType.replace('_', ' ')}</span>
                    </div>
                    {formatSalary() && (
                      <div className="flex items-center text-secondary-700">
                        <DollarSign className="h-5 w-5 mr-3 text-secondary-400" />
                        <span>{formatSalary()}</span>
                      </div>
                    )}
                    <div className="flex items-center text-secondary-700">
                      <Calendar className="h-5 w-5 mr-3 text-secondary-400" />
                      <span>
                        Posted{' '}
                        {job.createdAt && !isNaN(new Date(job.createdAt).getTime())
                          ? formatDistanceToNow(new Date(job.createdAt), { addSuffix: true })
                          : 'recently'}
                      </span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardContent className="pt-6">
                  <div className="space-y-4">
                    {job.applicationUrl ? (
                      <a href={job.applicationUrl} target="_blank" rel="noopener noreferrer">
                        <Button className="w-full" size="lg">
                          <ExternalLink className="h-4 w-4 mr-2" />
                          Apply Now
                        </Button>
                      </a>
                    ) : job.applicationEmail ? (
                      <a href={`mailto:${job.applicationEmail}?subject=Application for ${job.title}`}>
                        <Button className="w-full" size="lg">
                          <Mail className="h-4 w-4 mr-2" />
                          Apply via Email
                        </Button>
                      </a>
                    ) : (
                      <Button className="w-full" size="lg" disabled>
                        Application method not specified
                      </Button>
                    )}
                    <Link
                      href="/jobs"
                      className="block text-center text-sm text-primary-600 hover:text-primary-700"
                    >
                      Browse all jobs
                    </Link>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>

          {relatedJobs.length > 0 && (
            <section className="mt-12">
              <div className="mb-6 flex items-end justify-between gap-4">
                <div>
                  <h2 className="text-2xl font-bold text-secondary-900">Related jobs</h2>
                  <p className="text-secondary-600 text-sm">More roles you might like</p>
                </div>
                <Link href="/jobs" className="text-sm text-primary-600 hover:text-primary-700">
                  View all
                </Link>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {relatedJobs.map((related) => (
                  <JobCard key={related.id} job={related} />
                ))}
              </div>
            </section>
          )}
        </div>
      </main>
      <Footer />
    </div>
  )
}
