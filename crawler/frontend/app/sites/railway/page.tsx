'use client'

import { useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { ExternalLink, ArrowRight } from 'lucide-react'
import Link from 'next/link'

// Railway careers page reference configuration
// This is used as a template/reference to help configure crawlers
const RAILWAY_REFERENCE = {
  name: '',
  base_url: 'https://railway.com/careers#open-positions',
  backend_startup_id: '',
  schedule: '0 0 * * *', // Daily at midnight
  crawl_interval: 'daily',
  pagination_config: {
    type: 'link_follow',
    next_page_selector: 'a[aria-label="Next page"], a.next-page, button[aria-label="Next"]',
  },
  extraction_rules: {
    // Railway careers page structure - jobs are listed in sections by category
    job_list_selector: 'section h3, div[class*="position"], article, li[class*="job"]',
    job_detail_url: {
      type: 'relative',
      selector: 'a[href*="/careers/"], a[href*="/jobs/"], a[href*="careers"], a',
      attribute: 'href',
      base_url: 'https://railway.com',
    },
    fields: {
      title: {
        selector: 'h3, h2, [class*="title"], [class*="Title"]',
        type: 'text',
        required: true,
      },
      description: {
        selector: 'p, [class*="description"], [class*="Description"]',
        type: 'text',
        required: false,
      },
      location: {
        selector: 'span:contains("Anywhere"), [class*="location"], [class*="Location"]',
        type: 'text',
        required: false,
        default_value: 'Anywhere',
      },
      category: {
        selector: 'h2, h3[class*="category"], [class*="Category"]',
        type: 'text',
        required: false,
      },
      department: {
        selector: '[class*="department"], [class*="Department"], [class*="team"]',
        type: 'text',
        required: false,
      },
    },
  },
  deduplication_key: 'url',
  request_delay: 3,
  user_agent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
}

export default function RailwayReferencePage() {
  const router = useRouter()
  const searchParams = useSearchParams()

  // If coming from the form, redirect with Railway template data
  useEffect(() => {
    const useTemplate = searchParams.get('use')
    if (useTemplate === 'railway') {
      // Store Railway config in sessionStorage and redirect to new site form
      sessionStorage.setItem('crawlerTemplate', JSON.stringify(RAILWAY_REFERENCE))
      router.push('/sites/new?template=railway')
    }
  }, [searchParams, router])

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Link
          href="/sites"
          className="inline-flex items-center text-gray-600 hover:text-gray-900 mb-6"
        >
          ← Back to Sites
        </Link>

        <div className="bg-white rounded-lg shadow p-8">
          <div className="mb-6">
            <h1 className="text-3xl font-bold text-gray-900 mb-2">Railway Careers Page Reference</h1>
            <p className="text-gray-600 mb-4">
              Use Railway's careers page structure as a reference to help configure your crawler. 
              This page shows the recommended selectors and configuration based on Railway's page structure.
            </p>
            <a
              href="https://railway.com/careers#open-positions"
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 hover:text-blue-800 inline-flex items-center gap-1"
            >
              View Railway Careers Page
              <ExternalLink className="h-4 w-4" />
            </a>
          </div>

          <div className="mb-8 p-4 bg-blue-50 border border-blue-200 rounded-lg">
            <h3 className="font-semibold text-blue-900 mb-3">Reference Configuration</h3>
            <div className="space-y-2 text-sm text-blue-800">
              <div>
                <strong>URL:</strong> <code className="bg-white px-2 py-1 rounded">{RAILWAY_REFERENCE.base_url}</code>
              </div>
              <div>
                <strong>Job List Selector:</strong> <code className="bg-white px-2 py-1 rounded text-xs">{RAILWAY_REFERENCE.extraction_rules.job_list_selector}</code>
              </div>
              <div>
                <strong>Job URL Selector:</strong> <code className="bg-white px-2 py-1 rounded text-xs">{RAILWAY_REFERENCE.extraction_rules.job_detail_url.selector}</code>
              </div>
              <div>
                <strong>Title Selector:</strong> <code className="bg-white px-2 py-1 rounded text-xs">{RAILWAY_REFERENCE.extraction_rules.fields.title.selector}</code>
              </div>
              <div>
                <strong>Location Selector:</strong> <code className="bg-white px-2 py-1 rounded text-xs">{RAILWAY_REFERENCE.extraction_rules.fields.location.selector || 'Default: Anywhere'}</code>
              </div>
            </div>
          </div>

          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-6">
            <h4 className="font-semibold text-yellow-900 mb-2">📝 How to Use This Reference</h4>
            <ul className="text-sm text-yellow-800 space-y-2">
              <li>• Click the button below to open the crawler configuration form</li>
              <li>• The form will be pre-filled with Railway's page structure as a starting point</li>
              <li>• Modify the selectors and configuration to match your target website</li>
              <li>• Use Railway's structure as a guide for similar career pages</li>
            </ul>
          </div>

          <div className="flex justify-end gap-4 pt-6 border-t">
            <Link
              href="/sites"
              className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
            >
              Cancel
            </Link>
            <Link
              href="/sites/new?template=railway"
              className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
              onClick={() => {
                sessionStorage.setItem('crawlerTemplate', JSON.stringify(RAILWAY_REFERENCE))
              }}
            >
              Use Railway as Template
              <ArrowRight className="h-5 w-5 ml-2" />
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}

