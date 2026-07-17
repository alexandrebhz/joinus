import type { Metadata } from 'next'
import { buildPageMetadata } from '@/lib/seo'

export const metadata: Metadata = buildPageMetadata({
  title: 'Login',
  canonicalPath: '/login',
  index: false,
})

export default function LoginLayout({ children }: { children: React.ReactNode }) {
  return children
}
