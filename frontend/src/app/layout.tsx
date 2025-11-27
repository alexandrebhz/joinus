import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import { Providers } from './providers'
import { GoogleAnalytics } from '@/presentation/components/analytics/google-analytics'
import { defaultMetadata } from './metadata'

const inter = Inter({ subsets: ['latin'], variable: '--font-inter' })

export const metadata: Metadata = defaultMetadata

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" className={inter.variable}>
      <body>
        <GoogleAnalytics />
        <Providers>{children}</Providers>
      </body>
    </html>
  )
}

