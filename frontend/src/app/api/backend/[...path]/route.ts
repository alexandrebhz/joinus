import { NextRequest, NextResponse } from 'next/server'
import {
  ACCESS_COOKIE,
  REFRESH_COOKIE,
  cookieOptions,
  getAccessToken,
  getRefreshToken,
} from '@/lib/auth-cookies'

export const dynamic = 'force-dynamic'

function backendBase(): string {
  const raw = (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080').replace(/\/+$/, '')
  if (raw.startsWith('http://') || raw.startsWith('https://')) return raw
  return `https://${raw}`
}

function buildTarget(path: string[], search: string): string {
  const joined = path.map(encodeURIComponent).join('/')
  return `${backendBase()}/api/v1/${joined}${search}`
}

async function tryRefresh(): Promise<{ access: string; refresh: string } | null> {
  const refresh = await getRefreshToken()
  if (!refresh) return null
  const res = await fetch(`${backendBase()}/api/v1/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refresh }),
    cache: 'no-store',
  })
  if (!res.ok) return null
  const body = await res.json()
  const access = body?.data?.access_token as string | undefined
  const newRefresh = (body?.data?.refresh_token as string | undefined) || refresh
  if (!access) return null
  return { access, refresh: newRefresh }
}

function redactTokens(payload: unknown): unknown {
  if (!payload || typeof payload !== 'object') return payload
  const clone = JSON.parse(JSON.stringify(payload)) as Record<string, unknown>
  const data = clone.data as Record<string, unknown> | undefined
  if (data && typeof data === 'object') {
    delete data.access_token
    delete data.refresh_token
  }
  return clone
}

function attachAuthCookies(
  out: NextResponse,
  access?: string,
  refresh?: string
) {
  if (access) {
    out.cookies.set(ACCESS_COOKIE, access, cookieOptions(60 * 60 * 24))
  }
  if (refresh) {
    out.cookies.set(REFRESH_COOKIE, refresh, cookieOptions(60 * 60 * 24 * 7))
  }
}

function clearAuthCookiesOn(out: NextResponse) {
  out.cookies.set(ACCESS_COOKIE, '', { ...cookieOptions(0), maxAge: 0 })
  out.cookies.set(REFRESH_COOKIE, '', { ...cookieOptions(0), maxAge: 0 })
}

async function proxy(req: NextRequest, path: string[]): Promise<NextResponse> {
  const target = buildTarget(path, req.nextUrl.search)
  const headers = new Headers()
  const contentType = req.headers.get('content-type')
  if (contentType) headers.set('content-type', contentType)

  let access = await getAccessToken()
  if (access) headers.set('Authorization', `Bearer ${access}`)

  const body =
    req.method !== 'GET' && req.method !== 'HEAD'
      ? Buffer.from(await req.arrayBuffer())
      : undefined

  const doFetch = (authHeader?: string) => {
    const h = new Headers(headers)
    if (authHeader) h.set('Authorization', authHeader)
    return fetch(target, {
      method: req.method,
      headers: h,
      body,
      cache: 'no-store',
    })
  }

  let upstream = await doFetch(access ? `Bearer ${access}` : undefined)
  let refreshed: { access: string; refresh: string } | null = null

  const isAuthPath = path[0] === 'auth'
  if (upstream.status === 401 && !isAuthPath) {
    refreshed = await tryRefresh()
    if (refreshed) {
      upstream = await doFetch(`Bearer ${refreshed.access}`)
    }
  }

  const buf = await upstream.arrayBuffer()
  const text = new TextDecoder().decode(buf)

  let json: unknown
  try {
    json = JSON.parse(text)
  } catch {
    const out = new NextResponse(buf, {
      status: upstream.status,
      headers: {
        'Content-Type':
          upstream.headers.get('content-type') || 'application/octet-stream',
      },
    })
    if (refreshed) attachAuthCookies(out, refreshed.access, refreshed.refresh)
    return out
  }

  const data = (json as { data?: { access_token?: string; refresh_token?: string } })
    ?.data
  const out = NextResponse.json(redactTokens(json), { status: upstream.status })

  if (data?.access_token && data?.refresh_token) {
    attachAuthCookies(out, data.access_token, data.refresh_token)
  } else if (refreshed) {
    attachAuthCookies(out, refreshed.access, refreshed.refresh)
  }

  if (path[0] === 'auth' && path[1] === 'logout' && upstream.ok) {
    clearAuthCookiesOn(out)
  }

  return out
}

type Ctx = { params: Promise<{ path: string[] }> }

export async function GET(req: NextRequest, ctx: Ctx) {
  return proxy(req, (await ctx.params).path)
}
export async function POST(req: NextRequest, ctx: Ctx) {
  return proxy(req, (await ctx.params).path)
}
export async function PUT(req: NextRequest, ctx: Ctx) {
  return proxy(req, (await ctx.params).path)
}
export async function PATCH(req: NextRequest, ctx: Ctx) {
  return proxy(req, (await ctx.params).path)
}
export async function DELETE(req: NextRequest, ctx: Ctx) {
  return proxy(req, (await ctx.params).path)
}
