import Link from 'next/link'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { Card, CardContent } from '@/presentation/components/ui/card'
import { Button } from '@/presentation/components/ui/button'
import { XCircle } from 'lucide-react'

export default function BillingCancelPage() {
  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 flex items-center justify-center py-20">
        <div className="container mx-auto px-4 max-w-lg">
          <Card>
            <CardContent className="pt-10 pb-10 text-center">
              <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-secondary-100">
                <XCircle className="h-9 w-9 text-secondary-500" />
              </div>
              <h1 className="text-2xl font-bold text-secondary-900 mb-2">Checkout canceled</h1>
              <p className="text-secondary-600 mb-8">
                No charge was made. You can restart checkout anytime from your billing page.
              </p>
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
