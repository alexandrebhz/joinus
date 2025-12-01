'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { crawlerApiClient, CrawlSite, CrawlResult } from '@/lib/api-client'
import { Plus, Trash2, Play, Edit, ExternalLink, Clock, CheckCircle, XCircle, FileText } from 'lucide-react'

export default function SitesPage() {
  const [sites, setSites] = useState<CrawlSite[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchSites()
  }, [])

  const fetchSites = async () => {
    try {
      setLoading(true)
      const response = await crawlerApiClient.listSites()
      setSites(response.data || [])
      setError(null)
    } catch (err: any) {
      setError(err.message || 'Failed to fetch sites')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this site?')) return
    
    try {
      await crawlerApiClient.deleteSite(id)
      await fetchSites()
    } catch (err: any) {
      alert(err.message || 'Failed to delete site')
    }
  }

  const handleCrawl = async (id: string) => {
    try {
      const response = await crawlerApiClient.executeCrawl(id)
      const result = response.data
      const message = result 
        ? `Crawl completed! Found: ${result.jobs_found || 0}, Saved: ${result.jobs_saved || 0}, Skipped: ${result.jobs_skipped || 0}${result.errors?.length ? `, Errors: ${result.errors.length}` : ''}`
        : 'Crawl completed successfully'
      alert(message)
      await fetchSites()
    } catch (err: any) {
      alert(err.message || 'Failed to start crawl')
    }
  }

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'Never'
    return new Date(dateString).toLocaleString()
  }

  const getScheduleDisplay = (schedule: string) => {
    if (schedule.includes('daily')) return 'Daily'
    if (schedule.includes('weekly')) return 'Weekly'
    return schedule
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading sites...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8 flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Crawler Sites</h1>
            <p className="mt-2 text-gray-600">Manage websites to crawl for job listings</p>
          </div>
          <Link
            href="/sites/new"
            className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            <Plus className="h-5 w-5 mr-2" />
            Add New Site
          </Link>
        </div>

        {error && (
          <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
            {error}
          </div>
        )}

        {sites.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-12 text-center">
            <p className="text-gray-600 mb-4">No sites configured yet.</p>
            <Link
              href="/sites/new"
              className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              <Plus className="h-5 w-5 mr-2" />
              Add Your First Site
            </Link>
          </div>
        ) : (
          <div className="grid gap-6">
            {sites.map((site) => (
              <div key={site.id} className="bg-white rounded-lg shadow p-6">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-3 mb-2">
                      <h2 className="text-xl font-semibold text-gray-900">{site.name}</h2>
                      {site.active ? (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                          <CheckCircle className="h-3 w-3 mr-1" />
                          Active
                        </span>
                      ) : (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                          <XCircle className="h-3 w-3 mr-1" />
                          Inactive
                        </span>
                      )}
                    </div>
                    <a
                      href={site.base_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-blue-600 hover:text-blue-800 inline-flex items-center gap-1 mb-4"
                    >
                      {site.base_url}
                      <ExternalLink className="h-4 w-4" />
                    </a>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm text-gray-600">
                      <div>
                        <span className="font-medium">Schedule:</span> {getScheduleDisplay(site.schedule)}
                      </div>
                      <div>
                        <span className="font-medium">Interval:</span> {site.crawl_interval}
                      </div>
                      <div>
                        <span className="font-medium">Last Crawled:</span> {formatDate(site.last_crawled_at)}
                      </div>
                      <div>
                        <span className="font-medium">Next Crawl:</span> {formatDate(site.next_crawl_at)}
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-2 ml-4">
                    <button
                      onClick={() => handleCrawl(site.id)}
                      className="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                      title="Run crawl now"
                    >
                      <Play className="h-5 w-5" />
                    </button>
                    <Link
                      href={`/sites/${site.id}/edit`}
                      className="p-2 text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
                      title="Edit site"
                    >
                      <Edit className="h-5 w-5" />
                    </Link>
                    <Link
                      href={`/sites/${site.id}/logs`}
                      className="p-2 text-purple-600 hover:bg-purple-50 rounded-lg transition-colors"
                      title="View crawl logs"
                    >
                      <FileText className="h-5 w-5" />
                    </Link>
                    <button
                      onClick={() => handleDelete(site.id)}
                      className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                      title="Delete site"
                    >
                      <Trash2 className="h-5 w-5" />
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

