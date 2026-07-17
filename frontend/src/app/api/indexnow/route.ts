import { NextRequest, NextResponse } from 'next/server'
import { INDEXNOW_KEY, submitToIndexNow } from '@/lib/indexnow'
import { absoluteUrl, getSiteUrl } from '@/lib/seo'
import { apiClient } from '@/infrastructure/api/api-client'

export const dynamic = 'force-dynamic'

type Body = {
  urls?: string[]
  /** When true, submit static pages + active jobs/startups from the API. */
  all?: boolean
}

async function collectEntityUrls(): Promise<string[]> {
  const urls: string[] = []
  const pageSize = 100
  const maxPages = 50

  for (let page = 1; page <= maxPages; page++) {
    const res = await apiClient.listJobs({ status: 'active', page, page_size: pageSize })
    const batch = res.data || []
    for (const job of batch) {
      urls.push(absoluteUrl(`/jobs/${job.id}`))
    }
    if (batch.length < pageSize) break
    if (res.meta && page >= res.meta.totalPages) break
  }

  for (let page = 1; page <= maxPages; page++) {
    const res = await apiClient.listStartups({ status: 'active', page, page_size: pageSize })
    const batch = res.data || []
    for (const startup of batch) {
      if (startup.slug) urls.push(absoluteUrl(`/startups/${startup.slug}`))
    }
    if (batch.length < pageSize) break
    if (res.meta && page >= res.meta.totalPages) break
  }

  return urls
}

async function collectPublicUrls(): Promise<string[]> {
  const urls = [
    absoluteUrl('/'),
    absoluteUrl('/jobs'),
    absoluteUrl('/startups'),
    absoluteUrl('/about'),
    absoluteUrl('/contact'),
  ]

  try {
    urls.push(...(await collectEntityUrls()))
  } catch (error) {
    console.error('IndexNow: failed to load entities', error)
  }

  return urls
}

/**
 * POST /api/indexnow
 * Body: { urls?: string[], all?: boolean }
 * Optional header: x-indexnow-secret matching INDEXNOW_SECRET (recommended in prod).
 */
export async function POST(request: NextRequest) {
  const secret = process.env.INDEXNOW_SECRET
  if (secret) {
    const provided = request.headers.get('x-indexnow-secret')
    if (provided !== secret) {
      return NextResponse.json({ error: 'Unauthorized' }, { status: 401 })
    }
  }

  let body: Body = {}
  try {
    body = (await request.json()) as Body
  } catch {
    body = {}
  }

  const urls = body.all || !body.urls?.length ? await collectPublicUrls() : body.urls

  // IndexNow accepts up to 10,000 URLs per request; chunk to stay safe.
  const chunkSize = 1000
  let submitted = 0
  let lastStatus: number | undefined
  let error: string | undefined

  for (let i = 0; i < urls.length; i += chunkSize) {
    const chunk = urls.slice(i, i + chunkSize)
    const result = await submitToIndexNow(chunk)
    submitted += result.submitted
    lastStatus = result.status
    if (!result.ok) {
      error = result.error
      break
    }
  }

  return NextResponse.json({
    host: new URL(getSiteUrl()).host,
    key: INDEXNOW_KEY,
    submitted,
    ok: !error,
    status: lastStatus,
    error,
  })
}

export async function GET() {
  return NextResponse.json({
    endpoint: '/api/indexnow',
    method: 'POST',
    keyLocation: absoluteUrl(`/${INDEXNOW_KEY}.txt`),
    usage: {
      all: 'POST { "all": true }',
      urls: 'POST { "urls": ["/jobs/id", "https://joinus.ie/startups/slug"] }',
    },
  })
}
