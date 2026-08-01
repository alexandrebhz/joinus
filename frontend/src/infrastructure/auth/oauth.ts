export function getApiBaseUrl(): string {
  const raw = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
  return raw.trim().replace(/\/+$/, '')
}

/** Full URL to start Google OAuth on the backend. */
export function getGoogleOAuthUrl(): string {
  return `${getApiBaseUrl()}/api/v1/auth/oauth/google`
}
