import { NextResponse } from 'next/server'
import { ACCESS_COOKIE, REFRESH_COOKIE, cookieOptions } from '@/lib/auth-cookies'

export const dynamic = 'force-dynamic'

/** DELETE clears HttpOnly auth cookies (used when logout API is unavailable). */
export async function DELETE() {
  const out = NextResponse.json({ ok: true })
  out.cookies.set(ACCESS_COOKIE, '', { ...cookieOptions(0), maxAge: 0 })
  out.cookies.set(REFRESH_COOKIE, '', { ...cookieOptions(0), maxAge: 0 })
  return out
}
