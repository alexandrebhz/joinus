'use client'

import { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { crawlerApiClient, CrawlLog, LogEntry } from '@/lib/api-client'
import { ArrowLeft, RefreshCw, Clock, CheckCircle, XCircle, AlertCircle } from 'lucide-react'

export default function CrawlLogsPage() {
  const params = useParams()
  const router = useRouter()
  const siteId = params.id as string
  const [logs, setLogs] = useState<CrawlLog[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchLogs()
  }, [siteId])

  const fetchLogs = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await crawlerApiClient.getCrawlLogs(siteId, 50)
      setLogs(response.data || [])
    } catch (err: any) {
      setError(err.message || 'Failed to fetch logs')
    } finally {
      setLoading(false)
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <CheckCircle className="h-5 w-5 text-green-600" />
      case 'failed':
        return <XCircle className="h-5 w-5 text-red-600" />
      case 'running':
        return <Clock className="h-5 w-5 text-blue-600 animate-spin" />
      default:
        return <AlertCircle className="h-5 w-5 text-gray-600" />
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-green-100 text-green-800'
      case 'failed':
        return 'bg-red-100 text-red-800'
      case 'running':
        return 'bg-blue-100 text-blue-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const formatDuration = (ms: number) => {
    if (ms < 1000) return `${ms}ms`
    if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`
    return `${(ms / 60000).toFixed(1)}m`
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString()
  }

  const getLogLevelColor = (level: string) => {
    switch (level) {
      case 'error':
        return 'text-red-600 bg-red-50'
      case 'warning':
        return 'text-yellow-600 bg-yellow-50'
      case 'info':
        return 'text-blue-600 bg-blue-50'
      default:
        return 'text-gray-600 bg-gray-50'
    }
  }

  if (loading && logs.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading logs...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-6 flex items-center justify-between">
          <Link
            href="/sites"
            className="inline-flex items-center text-gray-600 hover:text-gray-900"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Sites
          </Link>
          <button
            onClick={fetchLogs}
            disabled={loading}
            className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
          >
            <RefreshCw className={`h-4 w-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </button>
        </div>

        <h1 className="text-3xl font-bold text-gray-900 mb-8">Crawl Execution Logs</h1>

        {error && (
          <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
            {error}
          </div>
        )}

        {logs.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-8 text-center">
            <p className="text-gray-600">No crawl logs found for this site.</p>
            <p className="text-sm text-gray-500 mt-2">Run a crawl to see execution logs here.</p>
          </div>
        ) : (
          <div className="space-y-6">
            {logs.map((log) => (
              <div key={log.id} className="bg-white rounded-lg shadow p-6">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center gap-3">
                    {getStatusIcon(log.status)}
                    <div>
                      <div className="flex items-center gap-2">
                        <span className={`px-2 py-1 rounded text-sm font-medium ${getStatusColor(log.status)}`}>
                          {log.status.toUpperCase()}
                        </span>
                        {log.status === 'running' && (
                          <span className="text-sm text-gray-500">Running...</span>
                        )}
                      </div>
                      <p className="text-sm text-gray-500 mt-1">
                        Started: {formatDate(log.started_at)}
                        {log.completed_at && ` • Completed: ${formatDate(log.completed_at)}`}
                      </p>
                    </div>
                  </div>
                  {log.completed_at && (
                    <div className="text-right">
                      <p className="text-sm text-gray-500">Duration</p>
                      <p className="font-semibold">{formatDuration(log.duration_ms)}</p>
                    </div>
                  )}
                </div>

                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
                  <div className="bg-gray-50 rounded-lg p-3">
                    <p className="text-sm text-gray-500">Pages Crawled</p>
                    <p className="text-xl font-semibold">{log.pages_crawled}</p>
                  </div>
                  <div className="bg-blue-50 rounded-lg p-3">
                    <p className="text-sm text-blue-600">Jobs Found</p>
                    <p className="text-xl font-semibold text-blue-700">{log.jobs_found}</p>
                  </div>
                  <div className="bg-green-50 rounded-lg p-3">
                    <p className="text-sm text-green-600">Jobs Saved</p>
                    <p className="text-xl font-semibold text-green-700">{log.jobs_saved}</p>
                  </div>
                  <div className="bg-yellow-50 rounded-lg p-3">
                    <p className="text-sm text-yellow-600">Jobs Skipped</p>
                    <p className="text-xl font-semibold text-yellow-700">{log.jobs_skipped}</p>
                  </div>
                </div>

                {log.errors.length > 0 && (
                  <div className="mb-4">
                    <h3 className="text-sm font-semibold text-red-700 mb-2">Errors ({log.errors.length})</h3>
                    <div className="bg-red-50 border border-red-200 rounded-lg p-3">
                      <ul className="space-y-1">
                        {log.errors.map((err, idx) => (
                          <li key={idx} className="text-sm text-red-700">• {err}</li>
                        ))}
                      </ul>
                    </div>
                  </div>
                )}

                {log.logs.length > 0 && (
                  <div>
                    <h3 className="text-sm font-semibold text-gray-700 mb-2">Execution Logs ({log.logs.length})</h3>
                    <div className="bg-gray-50 rounded-lg p-3 max-h-64 overflow-y-auto">
                      <div className="space-y-1 font-mono text-xs">
                        {log.logs.map((entry: LogEntry, idx: number) => (
                          <div key={idx} className={`px-2 py-1 rounded ${getLogLevelColor(entry.level)}`}>
                            <span className="text-gray-500">
                              {new Date(entry.timestamp).toLocaleTimeString()}
                            </span>
                            <span className="ml-2 font-semibold uppercase">{entry.level}</span>
                            <span className="ml-2">{entry.message}</span>
                          </div>
                        ))}
                      </div>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}



