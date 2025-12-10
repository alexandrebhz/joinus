'use client'

import { useState, useEffect } from 'react'
import { useRouter, useParams } from 'next/navigation'
import { crawlerApiClient, UpdateSiteInput, PaginationConfig, ExtractionRules, CrawlSite } from '@/lib/api-client'
import { ArrowLeft, Save } from 'lucide-react'
import Link from 'next/link'

export default function EditSitePage() {
  const router = useRouter()
  const params = useParams()
  const siteId = params.id as string
  const [loading, setLoading] = useState(false)
  const [loadingSite, setLoadingSite] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [site, setSite] = useState<CrawlSite | null>(null)
  const [formData, setFormData] = useState<UpdateSiteInput>({})

  // Load site data
  useEffect(() => {
    const fetchSite = async () => {
      try {
        setLoadingSite(true)
        const response = await crawlerApiClient.getSite(siteId)
        const siteData = response.data
        setSite(siteData)
        // Pre-fill form with existing data
        setFormData({
          name: siteData.name,
          base_url: siteData.base_url,
          backend_startup_id: siteData.backend_startup_id,
          schedule: siteData.schedule,
          crawl_interval: siteData.crawl_interval,
          pagination_config: siteData.pagination_config as any,
          extraction_rules: siteData.extraction_rules as any,
          deduplication_key: siteData.deduplication_key,
          request_delay: siteData.request_delay,
          user_agent: siteData.user_agent,
          active: siteData.active,
        })
        setError(null)
      } catch (err: any) {
        setError(err.message || 'Failed to load site')
      } finally {
        setLoadingSite(false)
      }
    }

    if (siteId) {
      fetchSite()
    }
  }, [siteId])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      await crawlerApiClient.updateSite(siteId, formData)
      router.push('/sites')
    } catch (err: any) {
      setError(err.message || 'Failed to update site')
    } finally {
      setLoading(false)
    }
  }

  const updatePaginationConfig = (updates: Partial<PaginationConfig>) => {
    setFormData({
      ...formData,
      pagination_config: {
        ...(formData.pagination_config || {}),
        ...updates,
      } as any,
    })
  }

  const updateExtractionRules = (updates: Partial<ExtractionRules>) => {
    setFormData({
      ...formData,
      extraction_rules: {
        ...(formData.extraction_rules || {}),
        ...updates,
      } as any,
    })
  }

  const updateFieldRule = (fieldName: string, updates: Partial<any>) => {
    const currentRules = formData.extraction_rules?.fields || {}
    setFormData({
      ...formData,
      extraction_rules: {
        ...(formData.extraction_rules || {}),
        fields: {
          ...currentRules,
          [fieldName]: {
            ...currentRules[fieldName],
            ...updates,
          },
        },
      } as any,
    })
  }

  if (loadingSite) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading site...</p>
        </div>
      </div>
    )
  }

  if (!site) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-600 mb-4">Site not found</p>
          <Link href="/sites" className="text-blue-600 hover:text-blue-800">
            Back to Sites
          </Link>
        </div>
      </div>
    )
  }

  const paginationConfig = formData.pagination_config || site.pagination_config
  const extractionRules = formData.extraction_rules || site.extraction_rules

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Link
          href="/sites"
          className="inline-flex items-center text-gray-600 hover:text-gray-900 mb-6"
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Sites
        </Link>

        <h1 className="text-3xl font-bold text-gray-900 mb-8">Edit Crawler Site</h1>

        {error && (
          <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow p-6 space-y-8">
          {/* Basic Information */}
          <section>
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Basic Information</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Site Name *
                </label>
                <input
                  type="text"
                  required
                  value={formData.name ?? site.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="e.g., Railway Careers"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Base URL *
                </label>
                <input
                  type="url"
                  required
                  value={formData.base_url ?? site.base_url}
                  onChange={(e) => setFormData({ ...formData, base_url: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="https://example.com/careers"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Backend Startup ID *
                </label>
                <input
                  type="text"
                  required
                  value={formData.backend_startup_id ?? site.backend_startup_id}
                  onChange={(e) => setFormData({ ...formData, backend_startup_id: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="UUID of the startup in backend"
                />
              </div>
              <div className="flex items-center">
                <input
                  type="checkbox"
                  checked={formData.active ?? site.active}
                  onChange={(e) => setFormData({ ...formData, active: e.target.checked })}
                  className="mr-2"
                />
                <label className="text-sm font-medium text-gray-700">Active</label>
              </div>
            </div>
          </section>

          {/* Schedule */}
          <section>
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Schedule</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Crawl Interval *
                </label>
                <select
                  value={formData.crawl_interval ?? site.crawl_interval}
                  onChange={(e) => setFormData({ ...formData, crawl_interval: e.target.value as any })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="daily">Daily</option>
                  <option value="weekly">Weekly</option>
                  <option value="custom">Custom</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Cron Schedule *
                </label>
                <input
                  type="text"
                  required
                  value={formData.schedule ?? site.schedule}
                  onChange={(e) => setFormData({ ...formData, schedule: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="0 0 * * * (daily at midnight)"
                />
                <p className="mt-1 text-sm text-gray-500">
                  Cron expression (e.g., "0 0 * * *" for daily at midnight)
                </p>
              </div>
            </div>
          </section>

          {/* Pagination */}
          <section>
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Pagination</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Pagination Type *
                </label>
                <select
                  value={paginationConfig.type}
                  onChange={(e) => updatePaginationConfig({ type: e.target.value as any })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="query_param">Query Parameter</option>
                  <option value="url_pattern">URL Pattern</option>
                  <option value="link_follow">Follow Next Link</option>
                  <option value="api_pagination">API Pagination</option>
                </select>
              </div>

              {paginationConfig.type === 'query_param' && (
                <>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Parameter Name
                    </label>
                    <input
                      type="text"
                      value={paginationConfig.param_name || ''}
                      onChange={(e) => updatePaginationConfig({ param_name: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      placeholder="page"
                    />
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Start Page
                      </label>
                      <input
                        type="number"
                        value={paginationConfig.start_page || 1}
                        onChange={(e) => updatePaginationConfig({ start_page: parseInt(e.target.value) })}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Increment
                      </label>
                      <input
                        type="number"
                        value={paginationConfig.increment || 1}
                        onChange={(e) => updatePaginationConfig({ increment: parseInt(e.target.value) })}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                    </div>
                  </div>
                </>
              )}

              {paginationConfig.type === 'link_follow' && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Next Page Selector (CSS)
                  </label>
                  <input
                    type="text"
                    value={paginationConfig.next_page_selector || ''}
                    onChange={(e) => updatePaginationConfig({ next_page_selector: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="a.next-page"
                  />
                </div>
              )}
            </div>
          </section>

          {/* Extraction Rules */}
          <section>
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Extraction Rules</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Job List Selector (CSS) *
                </label>
                <input
                  type="text"
                  required
                  value={extractionRules.job_list_selector}
                  onChange={(e) => updateExtractionRules({ job_list_selector: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="div.job-listing"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Job Detail URL Selector (CSS) *
                </label>
                <input
                  type="text"
                  required
                  value={extractionRules.job_detail_url.selector}
                  onChange={(e) => updateExtractionRules({
                    job_detail_url: {
                      ...extractionRules.job_detail_url,
                      selector: e.target.value,
                    },
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="a.job-link"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  URL Attribute
                </label>
                <input
                  type="text"
                  value={extractionRules.job_detail_url.attribute || 'href'}
                  onChange={(e) => updateExtractionRules({
                    job_detail_url: {
                      ...extractionRules.job_detail_url,
                      attribute: e.target.value,
                    },
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="href"
                />
              </div>

              {/* Field Rules */}
              <div className="border-t pt-4">
                <h3 className="text-lg font-medium text-gray-900 mb-4">Field Extraction Rules</h3>
                {Object.entries(extractionRules.fields).map(([fieldName, rule]) => (
                  <div key={fieldName} className="mb-4 p-4 border rounded-lg">
                    <h4 className="font-medium text-gray-900 mb-3 capitalize">{fieldName}</h4>
                    <div className="space-y-3">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                          Selector (CSS)
                        </label>
                        <input
                          type="text"
                          value={rule.selector}
                          onChange={(e) => updateFieldRule(fieldName, { selector: e.target.value })}
                          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                          Type
                        </label>
                        <select
                          value={rule.type}
                          onChange={(e) => updateFieldRule(fieldName, { type: e.target.value })}
                          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        >
                          <option value="text">Text</option>
                          <option value="html">HTML</option>
                          <option value="attribute">Attribute</option>
                          <option value="regex">Regex</option>
                        </select>
                      </div>
                      {rule.type === 'attribute' && (
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-1">
                            Attribute Name
                          </label>
                          <input
                            type="text"
                            value={rule.attribute || ''}
                            onChange={(e) => updateFieldRule(fieldName, { attribute: e.target.value })}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                          />
                        </div>
                      )}
                      <div className="flex items-center">
                        <input
                          type="checkbox"
                          checked={rule.required}
                          onChange={(e) => updateFieldRule(fieldName, { required: e.target.checked })}
                          className="mr-2"
                        />
                        <label className="text-sm text-gray-700">Required field</label>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </section>

          {/* Advanced Settings */}
          <section>
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Advanced Settings</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Deduplication Key
                </label>
                <select
                  value={formData.deduplication_key ?? site.deduplication_key}
                  onChange={(e) => setFormData({ ...formData, deduplication_key: e.target.value as any })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="url">URL</option>
                  <option value="composite">Composite</option>
                  <option value="external_id">External ID</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Request Delay (seconds)
                </label>
                <input
                  type="number"
                  min="0"
                  max="60"
                  value={formData.request_delay ?? site.request_delay}
                  onChange={(e) => setFormData({ ...formData, request_delay: parseInt(e.target.value) })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  User Agent
                </label>
                <input
                  type="text"
                  value={formData.user_agent ?? site.user_agent}
                  onChange={(e) => setFormData({ ...formData, user_agent: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
            </div>
          </section>

          <div className="flex justify-end gap-4 pt-6 border-t">
            <Link
              href="/sites"
              className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
            >
              Cancel
            </Link>
            <button
              type="submit"
              disabled={loading}
              className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <Save className="h-5 w-5 mr-2" />
              {loading ? 'Updating...' : 'Update Site'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}






