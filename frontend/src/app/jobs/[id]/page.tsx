import { notFound } from 'next/navigation'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { apiClient } from '@/infrastructure/api/api-client'
import { Button } from '@/presentation/components/ui/button'
import { Card, CardContent } from '@/presentation/components/ui/card'
import { MapPin, Briefcase, DollarSign, Calendar, ExternalLink, Mail } from 'lucide-react'
import { formatDistanceToNow } from 'date-fns'

interface JobDetailPageProps {
  params: {
    id: string
  }
}

async function getJob(id: string) {
  try {
    const response = await apiClient.getJob(id)
    return response.data
  } catch {
    return null
  }
}

export default async function JobDetailPage({ params }: JobDetailPageProps) {
  const job = await getJob(params.id)

  if (!job) {
    notFound()
  }

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
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-secondary-900 mb-2">{job.title}</h1>
            {job.startupName && (
              <p className="text-lg text-primary-600 font-medium">{job.startupName}</p>
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
                      <span>Posted {formatDistanceToNow(new Date(job.createdAt), { addSuffix: true })}</span>
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
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  )
}

