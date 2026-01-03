import type { Metadata } from 'next'

export const siteConfig = {
  name: 'JoinUs',
  description: 'Discover your next career opportunity at innovative startups. Browse remote, hybrid, and onsite tech jobs from fast-growing companies. Join thousands of professionals finding their dream roles.',
  url: 'https://joinus.ie',
  ogImage: '/og-image.jpg', // You'll need to create this
  links: {
    twitter: 'https://twitter.com/joinus',
    github: 'https://github.com/joinus',
  },
}

export const defaultMetadata: Metadata = {
  metadataBase: new URL(siteConfig.url),
  title: {
    default: 'JoinUs - Startup Job Board | Find Tech Jobs at Innovative Startups',
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
  authors: [
    {
      name: 'JoinUs',
    },
  ],
  creator: 'JoinUs',
  openGraph: {
    type: 'website',
    locale: 'en_US',
    url: siteConfig.url,
    title: 'JoinUs - Startup Job Board | Find Tech Jobs at Innovative Startups',
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
    title: 'JoinUs - Startup Job Board | Find Tech Jobs at Innovative Startups',
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
  verification: {
    google: 'your-google-verification-code', // Add your Google Search Console verification code
  },
}

