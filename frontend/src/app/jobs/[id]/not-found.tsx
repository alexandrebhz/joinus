import Link from 'next/link'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { Button } from '@/presentation/components/ui/button'
import { Card, CardContent } from '@/presentation/components/ui/card'
import { Briefcase, ArrowLeft } from 'lucide-react'

export default function JobNotFound() {
  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 py-12">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-4xl">
          <Card>
            <CardContent className="pt-6 text-center py-12">
              <Briefcase className="h-16 w-16 mx-auto text-secondary-400 mb-4" />
              <h1 className="text-3xl font-bold text-secondary-900 mb-2">Job Not Found</h1>
              <p className="text-secondary-600 mb-6">
                The job you're looking for doesn't exist or may have been removed.
              </p>
              <div className="flex gap-4 justify-center">
                <Link href="/jobs">
                  <Button>
                    <ArrowLeft className="h-4 w-4 mr-2" />
                    Browse All Jobs
                  </Button>
                </Link>
                <Link href="/">
                  <Button variant="outline">
                    Go Home
                  </Button>
                </Link>
              </div>
            </CardContent>
          </Card>
        </div>
      </main>
      <Footer />
    </div>
  )
}

