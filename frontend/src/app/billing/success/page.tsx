'use client'

import { Suspense } from 'react'
import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { Card, CardContent } from '@/presentation/components/ui/card'
import { Button } from '@/presentation/components/ui/button'
import { CheckCircle2 } from 'lucide-react'

export default function BillingSuccessPage() {
  return (
    <Suspense fallback={null}>
      <BillingSuccessContent />
    </Suspense>
  )
}

function BillingSuccessContent() {
  const searchParams = useSearchParams()
  const sessionId = searchParams.get('session_id')

  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 flex items-center justify-center py-20">
        <div className="container mx-auto px-4 max-w-lg">
          <Card>
            <CardContent className="pt-10 pb-10 text-center">
              <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-success-50">
                <CheckCircle2 className="h-9 w-9 text-success-600" />
              </div>
              <h1 className="text-2xl font-bold text-secondary-900 mb-2">Payment successful!</h1>
              <p className="text-secondary-600 mb-8">
                Thanks for your purchase. It may take a few seconds for your account to update once Stripe confirms
                the payment.
              </p>
              {sessionId && (
                <p className="text-xs text-secondary-400 mb-6 break-all">Reference: {sessionId}</p>
              )}
              <Link href="/dashboard/billing">
                <Button className="w-full">Back to Billing</Button>
              </Link>
            </CardContent>
          </Card>
        </div>
      </main>
      <Footer />
    </div>
  )
}
