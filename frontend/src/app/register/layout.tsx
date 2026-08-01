import type { Metadata } from 'next'
import { buildPageMetadata } from '@/lib/seo'

export const metadata: Metadata = buildPageMetadata({
  title: 'Sign Up',
  canonicalPath: '/register',
  index: false,
})

export default function RegisterLayout({ children }: { children: React.ReactNode }) {
  return children
}
