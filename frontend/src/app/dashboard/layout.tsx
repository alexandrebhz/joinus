import type { Metadata } from 'next'
import { buildPageMetadata } from '@/lib/seo'

export const metadata: Metadata = buildPageMetadata({
  title: 'Dashboard',
  canonicalPath: '/dashboard',
  index: false,
})

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return children
}
