'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { useAuthStore } from '@/infrastructure/store/auth.store'
import { apiClient } from '@/infrastructure/api/api-client'
import { Startup } from '@/domain/entities/startup.entity'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/presentation/components/ui/card'
import { Button } from '@/presentation/components/ui/button'
import { Input } from '@/presentation/components/ui/input'
import { Building2, Plus, Edit, ExternalLink, Loader2 } from 'lucide-react'
import Link from 'next/link'
import Image from 'next/image'

export default function DashboardStartupsPage() {
  const router = useRouter()
  const { user, isAuthenticated } = useAuthStore()
  const [startups, setStartups] = useState<Startup[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [editingStartup, setEditingStartup] = useState<Startup | null>(null)

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login')
      return
    }
    loadStartups()
  }, [isAuthenticated, router])

  const loadStartups = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await apiClient.listStartups()
      setStartups(response.data || [])
    } catch (err: any) {
      setError(err.message || 'Failed to load startups')
    } finally {
      setLoading(false)
    }
  }

  // Delete functionality not available yet - backend doesn't have DELETE endpoint
  // const handleDelete = async (id: string) => {
  //   if (!confirm('Are you sure you want to delete this startup?')) {
  //     return
  //   }
  //   // TODO: Implement when backend DELETE endpoint is available
  // }

  if (!isAuthenticated || !user) {
    return null
  }

  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 py-12">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="mb-8 flex items-center justify-between">
            <div>
              <h1 className="text-4xl font-bold text-secondary-900 mb-2">My Startups</h1>
              <p className="text-secondary-600">Manage your startup profiles</p>
            </div>
            <Button onClick={() => setShowCreateForm(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Create Startup
            </Button>
          </div>

          {showCreateForm && (
            <CreateStartupForm
              onSuccess={() => {
                setShowCreateForm(false)
                loadStartups()
              }}
              onCancel={() => setShowCreateForm(false)}
            />
          )}

          {editingStartup && (
            <EditStartupForm
              startup={editingStartup}
              onSuccess={() => {
                setEditingStartup(null)
                loadStartups()
              }}
              onCancel={() => setEditingStartup(null)}
            />
          )}

          {loading ? (
            <div className="flex items-center justify-center py-12">
              <Loader2 className="h-8 w-8 animate-spin text-primary-600" />
            </div>
          ) : error ? (
            <Card>
              <CardContent className="pt-6">
                <p className="text-red-600">{error}</p>
                <Button onClick={loadStartups} className="mt-4">
                  Retry
                </Button>
              </CardContent>
            </Card>
          ) : startups.length === 0 ? (
            <Card>
              <CardContent className="pt-6 text-center py-12">
                <Building2 className="h-12 w-12 mx-auto text-secondary-400 mb-4" />
                <h3 className="text-lg font-semibold text-secondary-900 mb-2">No startups yet</h3>
                <p className="text-secondary-600 mb-4">Create your first startup profile to get started</p>
                <Button onClick={() => setShowCreateForm(true)}>
                  <Plus className="h-4 w-4 mr-2" />
                  Create Startup
                </Button>
              </CardContent>
            </Card>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {startups.map((startup) => (
                <Card key={startup.id} hover>
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      <div className="flex items-start space-x-3">
                        {startup.logoUrl ? (
                          <div className="relative h-12 w-12 rounded-lg overflow-hidden flex-shrink-0 border border-secondary-200">
                            <Image
                              src={startup.logoUrl}
                              alt={startup.name}
                              fill
                              className="object-cover"
                            />
                          </div>
                        ) : (
                          <div className="h-12 w-12 rounded-lg bg-primary-100 flex items-center justify-center flex-shrink-0">
                            <Building2 className="h-6 w-6 text-primary-600" />
                          </div>
                        )}
                        <div className="flex-1 min-w-0">
                          <CardTitle className="text-lg mb-1">{startup.name}</CardTitle>
                          <CardDescription className="line-clamp-2">{startup.description}</CardDescription>
                        </div>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-2 text-sm text-secondary-600 mb-4">
                      <div className="flex items-center">
                        <Building2 className="h-4 w-4 mr-2 flex-shrink-0" />
                        {startup.industry}
                      </div>
                      <div className="flex items-center">
                        <ExternalLink className="h-4 w-4 mr-2 flex-shrink-0" />
                        <a
                          href={startup.website}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="hover:text-primary-600 truncate"
                        >
                          {startup.website.replace(/^https?:\/\//, '')}
                        </a>
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <Link href={`/startups/${startup.slug}`} className="flex-1">
                        <Button variant="outline" size="sm" className="w-full">
                          View Public
                        </Button>
                      </Link>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setEditingStartup(startup)}
                        className="flex-1"
                      >
                        <Edit className="h-4 w-4 mr-2" />
                        Edit
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>
      </main>
      <Footer />
    </div>
  )
}

function CreateStartupForm({
  onSuccess,
  onCancel,
}: {
  onSuccess: () => void
  onCancel: () => void
}) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    website: '',
    founded_year: new Date().getFullYear(),
    industry: '',
    company_size: '',
    location: '',
    allow_public_join: true,
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      setLoading(true)
      await apiClient.createStartup(formData)
      onSuccess()
    } catch (err: any) {
      alert(err.message || 'Failed to create startup')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card className="mb-8">
      <CardHeader>
        <CardTitle>Create New Startup</CardTitle>
        <CardDescription>Add a new startup profile to your account</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-secondary-700 mb-1">
              Startup Name *
            </label>
            <Input
              required
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              placeholder="Acme Inc."
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-secondary-700 mb-1">
              Description *
            </label>
            <textarea
              required
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              placeholder="Describe your startup..."
              className="w-full px-3 py-2 border border-secondary-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
              rows={4}
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-secondary-700 mb-1">
                Website *
              </label>
              <Input
                required
                type="url"
                value={formData.website}
                onChange={(e) => setFormData({ ...formData, website: e.target.value })}
                placeholder="https://example.com"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-secondary-700 mb-1">
                Founded Year *
              </label>
              <Input
                required
                type="number"
                value={formData.founded_year}
                onChange={(e) =>
                  setFormData({ ...formData, founded_year: parseInt(e.target.value) })
                }
                min="1900"
                max={new Date().getFullYear()}
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-secondary-700 mb-1">
                Industry *
              </label>
              <Input
                required
                value={formData.industry}
                onChange={(e) => setFormData({ ...formData, industry: e.target.value })}
                placeholder="Technology"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-secondary-700 mb-1">
                Company Size *
              </label>
              <Input
                required
                value={formData.company_size}
                onChange={(e) => setFormData({ ...formData, company_size: e.target.value })}
                placeholder="1-10"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-secondary-700 mb-1">
              Location *
            </label>
            <Input
              required
              value={formData.location}
              onChange={(e) => setFormData({ ...formData, location: e.target.value })}
              placeholder="Dublin, Ireland"
            />
          </div>

          <div className="flex items-center">
            <input
              type="checkbox"
              id="allow_public_join"
              checked={formData.allow_public_join}
              onChange={(e) =>
                setFormData({ ...formData, allow_public_join: e.target.checked })
              }
              className="mr-2"
            />
            <label htmlFor="allow_public_join" className="text-sm text-secondary-700">
              Allow public to view this startup profile
            </label>
          </div>

          <div className="flex gap-2">
            <Button type="submit" disabled={loading}>
              {loading ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  Creating...
                </>
              ) : (
                'Create Startup'
              )}
            </Button>
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancel
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}

function EditStartupForm({
  startup,
  onSuccess,
  onCancel,
}: {
  startup: Startup
  onSuccess: () => void
  onCancel: () => void
}) {
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState({
    name: startup.name,
    description: startup.description,
    website: startup.website,
    industry: startup.industry,
    company_size: startup.companySize,
    location: startup.location,
    allow_public_join: startup.allowPublicJoin,
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      setLoading(true)
      await apiClient.updateStartup({
        id: startup.id,
        ...formData,
      })
      onSuccess()
    } catch (err: any) {
      alert(err.message || 'Failed to update startup')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card className="mb-8">
      <CardHeader>
        <CardTitle>Edit Startup</CardTitle>
        <CardDescription>Update your startup profile information</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-secondary-700 mb-1">
              Startup Name *
            </label>
            <Input
              required
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-secondary-700 mb-1">
              Description *
            </label>
            <textarea
              required
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              className="w-full px-3 py-2 border border-secondary-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
              rows={4}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-secondary-700 mb-1">
              Website *
            </label>
            <Input
              required
              type="url"
              value={formData.website}
              onChange={(e) => setFormData({ ...formData, website: e.target.value })}
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-secondary-700 mb-1">
                Industry *
              </label>
              <Input
                required
                value={formData.industry}
                onChange={(e) => setFormData({ ...formData, industry: e.target.value })}
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-secondary-700 mb-1">
                Company Size *
              </label>
              <Input
                required
                value={formData.company_size}
                onChange={(e) => setFormData({ ...formData, company_size: e.target.value })}
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-secondary-700 mb-1">
              Location *
            </label>
            <Input
              required
              value={formData.location}
              onChange={(e) => setFormData({ ...formData, location: e.target.value })}
            />
          </div>

          <div className="flex items-center">
            <input
              type="checkbox"
              id="edit_allow_public_join"
              checked={formData.allow_public_join}
              onChange={(e) =>
                setFormData({ ...formData, allow_public_join: e.target.checked })
              }
              className="mr-2"
            />
            <label htmlFor="edit_allow_public_join" className="text-sm text-secondary-700">
              Allow public to view this startup profile
            </label>
          </div>

          <div className="flex gap-2">
            <Button type="submit" disabled={loading}>
              {loading ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  Updating...
                </>
              ) : (
                'Update Startup'
              )}
            </Button>
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancel
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}

