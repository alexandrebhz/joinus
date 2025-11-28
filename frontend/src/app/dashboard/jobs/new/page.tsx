'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { useAuthStore } from '@/infrastructure/store/auth.store'
import { apiClient } from '@/infrastructure/api/api-client'
import { Startup } from '@/domain/entities/startup.entity'
import { CreateJobRequest } from '@/application/dto/job.dto'
import { JobType, LocationType } from '@/domain/entities/job.entity'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/presentation/components/ui/card'
import { Button } from '@/presentation/components/ui/button'
import { Input } from '@/presentation/components/ui/input'
import { Briefcase, Loader2, ArrowLeft } from 'lucide-react'
import Link from 'next/link'

export default function CreateJobPage() {
  const router = useRouter()
  const { user, isAuthenticated } = useAuthStore()
  const [startups, setStartups] = useState<Startup[]>([])
  const [loadingStartups, setLoadingStartups] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login')
      return
    }
    loadStartups()
  }, [isAuthenticated, router])

  const loadStartups = async () => {
    try {
      setLoadingStartups(true)
      const response = await apiClient.listStartups()
      setStartups(response.data || [])
    } catch (err: any) {
      setError(err.message || 'Failed to load startups')
    } finally {
      setLoadingStartups(false)
    }
  }

  const [formData, setFormData] = useState<CreateJobRequest>({
    startup_id: '',
    title: '',
    description: '',
    requirements: '',
    job_type: 'full_time',
    location_type: 'remote',
    city: '',
    country: '',
    salary_min: undefined,
    salary_max: undefined,
    currency: 'EUR',
    application_url: '',
    application_email: '',
    expires_at: '',
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    // Validation
    if (!formData.startup_id) {
      setError('Please select a startup')
      return
    }

    if (!formData.title || !formData.description || !formData.requirements) {
      setError('Please fill in all required fields')
      return
    }

    if (formData.location_type !== 'remote' && !formData.country) {
      setError('Please provide a country for non-remote positions')
      return
    }

    if (formData.salary_min && formData.salary_max && formData.salary_min > formData.salary_max) {
      setError('Minimum salary cannot be greater than maximum salary')
      return
    }

    if (formData.application_url && formData.application_email) {
      setError('Please provide either an application URL or email, not both')
      return
    }

    if (!formData.application_url && !formData.application_email) {
      setError('Please provide either an application URL or email')
      return
    }

    try {
      setSubmitting(true)
      
      // Clean up form data - remove empty optional fields
      const submitData: CreateJobRequest = {
        ...formData,
        city: formData.city || undefined,
        salary_min: formData.salary_min || undefined,
        salary_max: formData.salary_max || undefined,
        application_url: formData.application_url || undefined,
        application_email: formData.application_email || undefined,
        expires_at: formData.expires_at || undefined,
      }

      await apiClient.createJob(submitData)
      
      // Redirect to dashboard or jobs list
      router.push('/dashboard')
    } catch (err: any) {
      setError(err.message || 'Failed to create job')
    } finally {
      setSubmitting(false)
    }
  }

  if (!isAuthenticated || !user) {
    return null
  }

  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 py-12">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-4xl">
          <div className="mb-8">
            <Link href="/dashboard" className="inline-flex items-center text-primary-600 hover:text-primary-700 mb-4">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Dashboard
            </Link>
            <h1 className="text-4xl font-bold text-secondary-900 mb-2">Post a New Job</h1>
            <p className="text-secondary-600">Create a job listing for your startup</p>
          </div>

          {loadingStartups ? (
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center justify-center py-12">
                  <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
                </div>
              </CardContent>
            </Card>
          ) : startups.length === 0 ? (
            <Card>
              <CardContent className="pt-6 text-center py-12">
                <Briefcase className="h-12 w-12 mx-auto text-secondary-400 mb-4" />
                <h3 className="text-lg font-semibold text-secondary-900 mb-2">No Startups Found</h3>
                <p className="text-secondary-600 mb-4">
                  You need to create a startup profile before posting jobs
                </p>
                <Link href="/dashboard/startups">
                  <Button>
                    Create Startup
                  </Button>
                </Link>
              </CardContent>
            </Card>
          ) : (
            <Card>
              <CardHeader>
                <CardTitle>Job Details</CardTitle>
                <CardDescription>Fill in the information about the position you're hiring for</CardDescription>
              </CardHeader>
              <CardContent>
                {error && (
                  <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-md">
                    <p className="text-red-600 text-sm">{error}</p>
                  </div>
                )}

                <form onSubmit={handleSubmit} className="space-y-6">
                  {/* Startup Selection */}
                  <div>
                    <label className="block text-sm font-medium text-secondary-700 mb-1">
                      Startup *
                    </label>
                    <select
                      required
                      value={formData.startup_id}
                      onChange={(e) => setFormData({ ...formData, startup_id: e.target.value })}
                      className="w-full px-3 py-2 border border-secondary-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                    >
                      <option value="">Select a startup</option>
                      {startups.map((startup) => (
                        <option key={startup.id} value={startup.id}>
                          {startup.name}
                        </option>
                      ))}
                    </select>
                  </div>

                  {/* Title */}
                  <div>
                    <label className="block text-sm font-medium text-secondary-700 mb-1">
                      Job Title *
                    </label>
                    <Input
                      required
                      value={formData.title}
                      onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                      placeholder="e.g., Senior Software Engineer"
                    />
                  </div>

                  {/* Description */}
                  <div>
                    <label className="block text-sm font-medium text-secondary-700 mb-1">
                      Job Description *
                    </label>
                    <textarea
                      required
                      value={formData.description}
                      onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                      placeholder="Describe the role, responsibilities, and what makes it exciting..."
                      className="w-full px-3 py-2 border border-secondary-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                      rows={6}
                    />
                  </div>

                  {/* Requirements */}
                  <div>
                    <label className="block text-sm font-medium text-secondary-700 mb-1">
                      Requirements *
                    </label>
                    <textarea
                      required
                      value={formData.requirements}
                      onChange={(e) => setFormData({ ...formData, requirements: e.target.value })}
                      placeholder="List the skills, experience, and qualifications required..."
                      className="w-full px-3 py-2 border border-secondary-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                      rows={6}
                    />
                  </div>

                  {/* Job Type and Location Type */}
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-secondary-700 mb-1">
                        Job Type *
                      </label>
                      <select
                        required
                        value={formData.job_type}
                        onChange={(e) =>
                          setFormData({ ...formData, job_type: e.target.value as JobType })
                        }
                        className="w-full px-3 py-2 border border-secondary-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                      >
                        <option value="full_time">Full Time</option>
                        <option value="part_time">Part Time</option>
                        <option value="contract">Contract</option>
                        <option value="internship">Internship</option>
                      </select>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-secondary-700 mb-1">
                        Location Type *
                      </label>
                      <select
                        required
                        value={formData.location_type}
                        onChange={(e) =>
                          setFormData({ ...formData, location_type: e.target.value as LocationType })
                        }
                        className="w-full px-3 py-2 border border-secondary-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                      >
                        <option value="remote">Remote</option>
                        <option value="hybrid">Hybrid</option>
                        <option value="onsite">Onsite</option>
                      </select>
                    </div>
                  </div>

                  {/* Location Details */}
                  {formData.location_type !== 'remote' && (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-medium text-secondary-700 mb-1">
                          City
                        </label>
                        <Input
                          value={formData.city}
                          onChange={(e) => setFormData({ ...formData, city: e.target.value })}
                          placeholder="e.g., Dublin"
                        />
                      </div>

                      <div>
                        <label className="block text-sm font-medium text-secondary-700 mb-1">
                          Country *
                        </label>
                        <Input
                          required
                          value={formData.country}
                          onChange={(e) => setFormData({ ...formData, country: e.target.value })}
                          placeholder="e.g., Ireland"
                        />
                      </div>
                    </div>
                  )}

                  {formData.location_type === 'remote' && (
                    <div>
                      <label className="block text-sm font-medium text-secondary-700 mb-1">
                        Country (for timezone reference)
                      </label>
                      <Input
                        value={formData.country}
                        onChange={(e) => setFormData({ ...formData, country: e.target.value })}
                        placeholder="e.g., Ireland"
                      />
                    </div>
                  )}

                  {/* Salary */}
                  <div className="border-t pt-6">
                    <h3 className="text-lg font-semibold text-secondary-900 mb-4">Salary (Optional)</h3>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      <div>
                        <label className="block text-sm font-medium text-secondary-700 mb-1">
                          Currency
                        </label>
                        <select
                          value={formData.currency}
                          onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
                          className="w-full px-3 py-2 border border-secondary-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                        >
                          <option value="EUR">EUR (€)</option>
                          <option value="USD">USD ($)</option>
                          <option value="GBP">GBP (£)</option>
                          <option value="CAD">CAD (C$)</option>
                          <option value="AUD">AUD (A$)</option>
                        </select>
                      </div>

                      <div>
                        <label className="block text-sm font-medium text-secondary-700 mb-1">
                          Min Salary
                        </label>
                        <Input
                          type="number"
                          min="0"
                          value={formData.salary_min || ''}
                          onChange={(e) =>
                            setFormData({
                              ...formData,
                              salary_min: e.target.value ? parseInt(e.target.value) : undefined,
                            })
                          }
                          placeholder="e.g., 50000"
                        />
                      </div>

                      <div>
                        <label className="block text-sm font-medium text-secondary-700 mb-1">
                          Max Salary
                        </label>
                        <Input
                          type="number"
                          min="0"
                          value={formData.salary_max || ''}
                          onChange={(e) =>
                            setFormData({
                              ...formData,
                              salary_max: e.target.value ? parseInt(e.target.value) : undefined,
                            })
                          }
                          placeholder="e.g., 80000"
                        />
                      </div>
                    </div>
                  </div>

                  {/* Application Method */}
                  <div className="border-t pt-6">
                    <h3 className="text-lg font-semibold text-secondary-900 mb-4">Application Method *</h3>
                    <div className="space-y-4">
                      <div>
                        <label className="block text-sm font-medium text-secondary-700 mb-1">
                          Application URL
                        </label>
                        <Input
                          type="url"
                          value={formData.application_url}
                          onChange={(e) => setFormData({ ...formData, application_url: e.target.value })}
                          placeholder="https://your-startup.com/careers/apply"
                        />
                        <p className="text-xs text-secondary-500 mt-1">
                          Link to your application form or job posting
                        </p>
                      </div>

                      <div className="relative">
                        <div className="absolute inset-0 flex items-center">
                          <div className="w-full border-t border-secondary-300"></div>
                        </div>
                        <div className="relative flex justify-center text-sm">
                          <span className="px-2 bg-white text-secondary-500">OR</span>
                        </div>
                      </div>

                      <div>
                        <label className="block text-sm font-medium text-secondary-700 mb-1">
                          Application Email
                        </label>
                        <Input
                          type="email"
                          value={formData.application_email}
                          onChange={(e) => setFormData({ ...formData, application_email: e.target.value })}
                          placeholder="careers@your-startup.com"
                        />
                        <p className="text-xs text-secondary-500 mt-1">
                          Email address for applications
                        </p>
                      </div>
                    </div>
                  </div>

                  {/* Expiration Date */}
                  <div>
                    <label className="block text-sm font-medium text-secondary-700 mb-1">
                      Expiration Date (Optional)
                    </label>
                    <Input
                      type="date"
                      value={formData.expires_at}
                      onChange={(e) => setFormData({ ...formData, expires_at: e.target.value })}
                      min={new Date().toISOString().split('T')[0]}
                    />
                    <p className="text-xs text-secondary-500 mt-1">
                      When should this job posting expire? Leave empty for no expiration.
                    </p>
                  </div>

                  {/* Submit Buttons */}
                  <div className="flex gap-4 pt-4 border-t">
                    <Button type="submit" disabled={submitting} className="flex-1">
                      {submitting ? (
                        <>
                          <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                          Creating Job...
                        </>
                      ) : (
                        <>
                          <Briefcase className="h-4 w-4 mr-2" />
                          Post Job
                        </>
                      )}
                    </Button>
                    <Link href="/dashboard">
                      <Button type="button" variant="outline">
                        Cancel
                      </Button>
                    </Link>
                  </div>
                </form>
              </CardContent>
            </Card>
          )}
        </div>
      </main>
      <Footer />
    </div>
  )
}

