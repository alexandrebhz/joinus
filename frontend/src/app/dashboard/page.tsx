'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { useAuthStore } from '@/infrastructure/store/auth.store'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/presentation/components/ui/card'
import { Button } from '@/presentation/components/ui/button'
import { Briefcase, Building2, Plus } from 'lucide-react'
import Link from 'next/link'

export default function DashboardPage() {
  const router = useRouter()
  const { user, isAuthenticated } = useAuthStore()

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login')
    }
  }, [isAuthenticated, router])

  if (!isAuthenticated || !user) {
    return null
  }

  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1 py-12">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-secondary-900 mb-2">Dashboard</h1>
            <p className="text-secondary-600">Welcome back, {user.name}!</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <Card hover>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <Briefcase className="h-5 w-5 mr-2 text-primary-600" />
                  Post a Job
                </CardTitle>
                <CardDescription>Create a new job listing for your startup</CardDescription>
              </CardHeader>
              <CardContent>
                <Link href="/dashboard/jobs/new">
                  <Button className="w-full">
                    <Plus className="h-4 w-4 mr-2" />
                    Create Job
                  </Button>
                </Link>
              </CardContent>
            </Card>

            <Card hover>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <Building2 className="h-5 w-5 mr-2 text-primary-600" />
                  My Startups
                </CardTitle>
                <CardDescription>Manage your startup profiles</CardDescription>
              </CardHeader>
              <CardContent>
                <Link href="/dashboard/startups">
                  <Button variant="outline" className="w-full">
                    View Startups
                  </Button>
                </Link>
              </CardContent>
            </Card>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  )
}

