import Link from 'next/link'
import { Plus, Settings, Play, List } from 'lucide-react'

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="text-center mb-12">
          <h1 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4">
            Job Crawler Management
          </h1>
          <p className="text-xl text-gray-600 max-w-2xl mx-auto">
            Configure and manage web crawlers to automatically collect job listings from career pages
          </p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
          <Link
            href="/sites"
            className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow border border-gray-200"
          >
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-blue-100 rounded-lg">
                <List className="h-6 w-6 text-blue-600" />
              </div>
            </div>
            <h2 className="text-xl font-semibold text-gray-900 mb-2">Manage Sites</h2>
            <p className="text-gray-600 text-sm">
              View and manage all configured crawler sites
            </p>
          </Link>

          <Link
            href="/sites/new"
            className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow border border-gray-200"
          >
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-green-100 rounded-lg">
                <Plus className="h-6 w-6 text-green-600" />
              </div>
            </div>
            <h2 className="text-xl font-semibold text-gray-900 mb-2">Create New Site</h2>
            <p className="text-gray-600 text-sm">
              Configure a new website to crawl for job listings
            </p>
          </Link>

          <Link
            href="/sites/railway"
            className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow border-2 border-blue-300"
          >
            <div className="flex items-center justify-between mb-4">
              <div className="p-3 bg-purple-100 rounded-lg">
                <Settings className="h-6 w-6 text-purple-600" />
              </div>
            </div>
            <h2 className="text-xl font-semibold text-gray-900 mb-2">Railway Reference</h2>
            <p className="text-gray-600 text-sm">
              Use Railway's careers page structure as a template to configure crawlers
            </p>
          </Link>
        </div>

        <div className="bg-white rounded-lg shadow-md p-8 border border-gray-200">
          <h2 className="text-2xl font-semibold text-gray-900 mb-4">How It Works</h2>
          <div className="grid md:grid-cols-3 gap-6">
            <div>
              <div className="flex items-center justify-center w-12 h-12 bg-blue-100 rounded-full mb-4">
                <span className="text-blue-600 font-bold text-xl">1</span>
              </div>
              <h3 className="font-semibold text-gray-900 mb-2">Configure Site</h3>
              <p className="text-gray-600 text-sm">
                Add a website URL and configure extraction rules using CSS selectors
              </p>
            </div>
            <div>
              <div className="flex items-center justify-center w-12 h-12 bg-green-100 rounded-full mb-4">
                <span className="text-green-600 font-bold text-xl">2</span>
              </div>
              <h3 className="font-semibold text-gray-900 mb-2">Schedule Crawls</h3>
              <p className="text-gray-600 text-sm">
                Set up automatic crawls using cron expressions (daily, weekly, etc.)
              </p>
            </div>
            <div>
              <div className="flex items-center justify-center w-12 h-12 bg-purple-100 rounded-full mb-4">
                <span className="text-purple-600 font-bold text-xl">3</span>
              </div>
              <h3 className="font-semibold text-gray-900 mb-2">Sync Jobs</h3>
              <p className="text-gray-600 text-sm">
                Crawled jobs are automatically synced to your backend API
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
