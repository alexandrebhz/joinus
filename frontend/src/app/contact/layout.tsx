import type { Metadata } from 'next'
import { buildPageMetadata } from '@/lib/seo'

export const metadata: Metadata = buildPageMetadata({
  title: 'Contact JoinUs',
  description: 'Get in touch with the JoinUs team about jobs, startups, or partnerships.',
  canonicalPath: '/contact',
})

export default function ContactLayout({ children }: { children: React.ReactNode }) {
  return children
}
