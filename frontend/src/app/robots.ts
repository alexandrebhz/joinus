import { MetadataRoute } from 'next'
import { getSiteUrl } from '@/lib/seo'

export default function robots(): MetadataRoute.Robots {
  const baseUrl = getSiteUrl()

  return {
    rules: [
      {
        userAgent: '*',
        allow: '/',
        disallow: ['/api/', '/dashboard/', '/login', '/register'],
      },
    ],
    sitemap: `${baseUrl}/sitemap.xml`,
    host: baseUrl.replace(/^https?:\/\//, ''),
  }
}
