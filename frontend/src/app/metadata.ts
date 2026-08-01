import type { Metadata } from 'next'
import { absoluteUrl, getSiteUrl, truncateTitle } from '@/lib/seo'

export const siteConfig = {
  name: 'JoinUs',
  description:
    'Discover your next career opportunity at innovative startups. Browse remote, hybrid, and onsite tech jobs from fast-growing companies.',
  url: getSiteUrl(),
  ogImage: '/og-image.jpg',
  links: {
    twitter: 'https://twitter.com/joinus',
    github: 'https://github.com/joinus',
  },
}

const defaultTitle = truncateTitle('JoinUs – Startup Job Board | Tech Jobs')

export const defaultMetadata: Metadata = {
  metadataBase: new URL(siteConfig.url),
  title: {
    default: defaultTitle,
    template: '%s | JoinUs',
  },
  description: siteConfig.description,
  keywords: [
    'startup jobs',
    'tech jobs',
    'remote jobs',
    'software engineer jobs',
    'developer jobs',
    'startup careers',
    'tech careers',
    'job board',
    'startup hiring',
    'tech recruitment',
    'software jobs',
    'engineering jobs',
    'product jobs',
    'design jobs',
    'marketing jobs',
    'sales jobs',
  ],
  authors: [{ name: 'JoinUs' }],
  creator: 'JoinUs',
  alternates: {
    canonical: absoluteUrl('/'),
  },
  openGraph: {
    type: 'website',
    locale: 'en_IE',
    url: siteConfig.url,
    title: defaultTitle,
    description: siteConfig.description,
    siteName: 'JoinUs',
    images: [
      {
        url: siteConfig.ogImage,
        width: 1200,
        height: 630,
        alt: 'JoinUs - Startup Job Board',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: defaultTitle,
    description: siteConfig.description,
    images: [siteConfig.ogImage],
    creator: '@joinus',
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
    },
  },
}
