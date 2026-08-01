'use client'

import { Suspense, useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { useAuthStore } from '@/infrastructure/store/auth.store'
import { apiClient } from '@/infrastructure/api/api-client'
import { Job } from '@/domain/entities/job.entity'
import { StartupBillingStatus } from '@/application/dto/billing.dto'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/presentation/components/ui/card'
import { Button } from '@/presentation/components/ui/button'
import { Crown, Zap, Loader2, CheckCircle2, Briefcase } from 'lucide-react'

export default function BillingPage() {
  return (
    <Suspense fallback={null}>
      <BillingPageContent />
    </Suspense>
  )
}

function BillingPageContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const highlightJobId = searchParams.get('boost_job_id')
  const { user, isAuthenticated } = useAuthStore()

  const [startups, setStartups] = useState<StartupBillingStatus[]>([])
  const [jobs, setJobs] = useState<Job[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [pendingKey, setPendingKey] = useState<string | null>(null)

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login')
      return
    }
    loadBillingData()
  }, [isAuthenticated, router])

  const loadBillingData = async () => {
    try {
      setLoading(true)
      setError(null)

      const statusRes = await apiClient.getBillingStatus()
      const myStartups = statusRes.data?.startups || []
      setStartups(myStartups)

      const jobLists = await Promise.all(
        myStartups.map((s) => apiClient.listJobs({ startup_id: s.startupId, page_size: 100 }))
      )
      setJobs(jobLists.flatMap((res) => res.data || []))
    } catch (err: any) {
      setError(err.message || 'Failed to load billing information')
    } finally {
      setLoading(false)
    }
  }

  const handleCheckout = async (product: 'job_boost' | 'startup_pro', ids: { jobId?: string; startupId?: string }) => {
    const key = product === 'job_boost' ? `job:${ids.jobId}` : `startup:${ids.startupId}`
    try {
      setPendingKey(key)
      const response = await apiClient.createCheckout({
        product,
        job_id: ids.jobId,
        startup_id: ids.startupId,
      })
      window.location.href = response.data.url
    } catch (err: any) {
      alert(err.message || 'Failed to start checkout')
      setPendingKey(null)
    }
  }

  const isBoosted = (job: Job) => !!job.boostedUntil && new Date(job.boostedUntil).getTime() > Date.now()
  const formatDate = (value?: string) => (value ? new Date(value).toLocaleDateString() : '')

  if (!isAuthenticated || !user) {
    return null
  }

  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 py-12">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-4xl">
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-secondary-900 mb-2">Billing &amp; Boosts</h1>
            <p className="text-secondary-600">Upgrade your startup or boost a job listing for more visibility</p>
          </div>

          {loading ? (
            <div className="flex items-center justify-center py-12">
              <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
            </div>
          ) : error ? (
            <Card>
              <CardContent className="pt-6">
                <p className="text-error-600">{error}</p>
                <Button onClick={loadBillingData} className="mt-4">
                  Retry
                </Button>
              </CardContent>
            </Card>
          ) : (
            <div className="space-y-10">
              {/* Startup Pro */}
              <section>
                <h2 className="text-2xl font-semibold text-secondary-900 mb-4">Startup Pro</h2>
                {startups.length === 0 ? (
                  <Card>
                    <CardContent className="pt-6 text-center py-8">
                      <p className="text-secondary-600 mb-4">Create a startup profile to unlock Startup Pro</p>
                    </CardContent>
                  </Card>
                ) : (
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {startups.map((startup) => {
                      const key = `startup:${startup.startupId}`
                      const isPro = startup.plan === 'pro'
                      return (
                        <Card key={startup.startupId} hover>
                          <CardHeader>
                            <div className="flex items-center justify-between">
                              <CardTitle className="text-lg">{startup.startupName}</CardTitle>
                              {isPro && (
                                <span className="inline-flex items-center gap-1 rounded-full bg-success-50 text-success-600 text-xs font-semibold px-2 py-0.5">
                                  <Crown className="h-3 w-3" />
                                  Pro
                                </span>
                              )}
                            </div>
                            <CardDescription>
                              {isPro
                                ? `Pro plan active${startup.planExpiresAt ? ` until ${formatDate(startup.planExpiresAt)}` : ''}`
                                : 'Free plan'}
                            </CardDescription>
                          </CardHeader>
                          <CardContent>
                            {isPro ? (
                              <div className="flex items-center text-success-600 text-sm font-medium">
                                <CheckCircle2 className="h-4 w-4 mr-2" />
                                You're all set
                              </div>
                            ) : (
                              <Button
                                className="w-full"
                                disabled={pendingKey === key}
                                onClick={() => handleCheckout('startup_pro', { startupId: startup.startupId })}
                              >
                                {pendingKey === key ? (
                                  <>
                                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                                    Redirecting...
                                  </>
                                ) : (
                                  <>
                                    <Crown className="h-4 w-4 mr-2" />
                                    Upgrade to Pro &mdash; €99/mo
                                  </>
                                )}
                              </Button>
                            )}
                          </CardContent>
                        </Card>
                      )
                    })}
                  </div>
                )}
              </section>

              {/* Job Boost */}
              <section>
                <h2 className="text-2xl font-semibold text-secondary-900 mb-4">Job Boost</h2>
                {jobs.length === 0 ? (
                  <Card>
                    <CardContent className="pt-6 text-center py-8">
                      <Briefcase className="h-10 w-10 mx-auto text-secondary-400 mb-3" />
                      <p className="text-secondary-600">Post a job to be able to boost it</p>
                    </CardContent>
                  </Card>
                ) : (
                  <div className="space-y-3">
                    {jobs.map((job) => {
                      const key = `job:${job.id}`
                      const boosted = isBoosted(job)
                      const highlighted = job.id === highlightJobId
                      return (
                        <Card key={job.id} className={highlighted ? 'ring-2 ring-warning-500' : ''}>
                          <CardContent className="pt-6 flex items-center justify-between flex-wrap gap-4">
                            <div>
                              <p className="font-semibold text-secondary-900">{job.title}</p>
                              <p className="text-sm text-secondary-600">{job.startupName}</p>
                            </div>
                            {boosted ? (
                              <span className="inline-flex items-center gap-1 rounded-full bg-warning-50 text-warning-600 text-sm font-semibold px-3 py-1">
                                <Zap className="h-4 w-4" />
                                Boosted until {formatDate(job.boostedUntil)}
                              </span>
                            ) : (
                              <Button
                                variant="outline"
                                disabled={pendingKey === key}
                                onClick={() => handleCheckout('job_boost', { jobId: job.id })}
                              >
                                {pendingKey === key ? (
                                  <>
                                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                                    Redirecting...
                                  </>
                                ) : (
                                  <>
                                    <Zap className="h-4 w-4 mr-2" />
                                    Boost this job &mdash; €49
                                  </>
                                )}
                              </Button>
                            )}
                          </CardContent>
                        </Card>
                      )
                    })}
                  </div>
                )}
              </section>
            </div>
          )}
        </div>
      </main>
      <Footer />
    </div>
  )
}
