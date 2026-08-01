import { absoluteUrl, getSiteUrl } from '@/lib/seo'

/** Host this exact value at /{INDEXNOW_KEY}.txt (see public/). */
export const INDEXNOW_KEY =
  process.env.INDEXNOW_KEY || process.env.NEXT_PUBLIC_INDEXNOW_KEY || 'joinus-indexnow-key-7f3a9c2e'

const INDEXNOW_ENDPOINT = 'https://api.indexnow.org/indexnow'

export type IndexNowResult = {
  submitted: number
  ok: boolean
  status?: number
  error?: string
}

/**
 * Submit URLs to IndexNow (Bing, Yandex, Seznam, Naver, etc.).
 * Paths may be absolute or site-relative (e.g. `/jobs/abc`).
 */
export async function submitToIndexNow(pathsOrUrls: string[]): Promise<IndexNowResult> {
  const host = new URL(getSiteUrl()).host
  const urlList = [...new Set(pathsOrUrls.map((u) => (u.startsWith('http') ? u : absoluteUrl(u))))]

  if (urlList.length === 0) {
    return { submitted: 0, ok: true }
  }

  try {
    const res = await fetch(INDEXNOW_ENDPOINT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json; charset=utf-8' },
      body: JSON.stringify({
        host,
        key: INDEXNOW_KEY,
        keyLocation: absoluteUrl(`/${INDEXNOW_KEY}.txt`),
        urlList,
      }),
    })

    // IndexNow returns 200/202 on success; 422 if some URLs invalid.
    const ok = res.status === 200 || res.status === 202
    return {
      submitted: urlList.length,
      ok,
      status: res.status,
      error: ok ? undefined : await res.text().catch(() => res.statusText),
    }
  } catch (error) {
    return {
      submitted: urlList.length,
      ok: false,
      error: error instanceof Error ? error.message : 'IndexNow request failed',
    }
  }
}
