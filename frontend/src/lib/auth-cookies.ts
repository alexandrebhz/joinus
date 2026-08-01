import { cookies } from 'next/headers'

export const ACCESS_COOKIE = 'ju_access_token'
export const REFRESH_COOKIE = 'ju_refresh_token'

const isProd = process.env.NODE_ENV === 'production'

export function cookieOptions(maxAgeSeconds: number) {
  return {
    httpOnly: true as const,
    secure: isProd,
    sameSite: 'lax' as const,
    path: '/',
    maxAge: maxAgeSeconds,
  }
}

export async function setAuthCookies(accessToken: string, refreshToken: string) {
  const jar = await cookies()
  jar.set(ACCESS_COOKIE, accessToken, cookieOptions(60 * 60 * 24)) // 24h
  jar.set(REFRESH_COOKIE, refreshToken, cookieOptions(60 * 60 * 24 * 7)) // 7d
}

export async function clearAuthCookies() {
  const jar = await cookies()
  jar.delete(ACCESS_COOKIE)
  jar.delete(REFRESH_COOKIE)
}

export async function getAccessToken(): Promise<string | undefined> {
  const jar = await cookies()
  return jar.get(ACCESS_COOKIE)?.value
}

export async function getRefreshToken(): Promise<string | undefined> {
  const jar = await cookies()
  return jar.get(REFRESH_COOKIE)?.value
}
