import Link from 'next/link'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { Button } from '@/presentation/components/ui/button'
import { apiClient } from '@/infrastructure/api/api-client'
import { JobCard } from '@/presentation/components/job/job-card'
import { StartupCard } from '@/presentation/components/startup/startup-card'
import { ArrowRight, Search, Briefcase, Building2, TrendingUp } from 'lucide-react'

async function getFeaturedJobs() {
  try {
    const response = await apiClient.listJobs({ page: 1, page_size: 6 })
    return response.data || []
  } catch {
    return []
  }
}

async function getFeaturedStartups() {
  try {
    const response = await apiClient.listStartups({ page: 1, page_size: 6 })
    return response.data || []
  } catch {
    return []
  }
}

export default async function HomePage() {
  const [jobs, startups] = await Promise.all([getFeaturedJobs(), getFeaturedStartups()])

  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1">
        {/* Hero Section */}
        <section className="relative overflow-hidden bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-20 sm:py-32">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="max-w-3xl mx-auto text-center">
              <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-secondary-900 mb-6">
                Find Your Next Opportunity at
                <span className="text-primary-600"> Innovative Startups</span>
              </h1>
              <p className="text-lg sm:text-xl text-secondary-600 mb-8">
                Connect with fast-growing startups and discover roles that match your skills and passion.
              </p>
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <Link href="/jobs">
                  <Button size="lg" className="w-full sm:w-auto">
                    <Search className="mr-2 h-5 w-5" />
                    Browse Jobs
                  </Button>
                </Link>
                <Link href="/startups">
                  <Button variant="outline" size="lg" className="w-full sm:w-auto">
                    <Building2 className="mr-2 h-5 w-5" />
                    Explore Startups
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section className="py-16 bg-white">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <div className="text-center">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary-100 text-primary-600 mb-4">
                  <Briefcase className="h-8 w-8" />
                </div>
                <h3 className="text-xl font-semibold text-secondary-900 mb-2">Curated Jobs</h3>
                <p className="text-secondary-600">
                  Hand-picked opportunities from verified startups looking for top talent.
                </p>
              </div>
              <div className="text-center">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary-100 text-primary-600 mb-4">
                  <Building2 className="h-8 w-8" />
                </div>
                <h3 className="text-xl font-semibold text-secondary-900 mb-2">Startup Profiles</h3>
                <p className="text-secondary-600">
                  Learn about company culture, team, and mission before applying.
                </p>
              </div>
              <div className="text-center">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary-100 text-primary-600 mb-4">
                  <TrendingUp className="h-8 w-8" />
                </div>
                <h3 className="text-xl font-semibold text-secondary-900 mb-2">Fast Growth</h3>
                <p className="text-secondary-600">
                  Join companies at the forefront of innovation and rapid expansion.
                </p>
              </div>
            </div>
          </div>
        </section>

        {/* Featured Jobs */}
        {jobs.length > 0 && (
          <section className="py-16 bg-secondary-50">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8">
              <div className="flex items-center justify-between mb-8">
                <div>
                  <h2 className="text-3xl font-bold text-secondary-900 mb-2">Featured Jobs</h2>
                  <p className="text-secondary-600">Latest opportunities from top startups</p>
                </div>
                <Link href="/jobs">
                  <Button variant="ghost">
                    View All
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Button>
                </Link>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {jobs.map((job) => (
                  <JobCard key={job.id} job={job} />
                ))}
              </div>
            </div>
          </section>
        )}

        {/* Featured Startups */}
        {startups.length > 0 && (
          <section className="py-16 bg-white">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8">
              <div className="flex items-center justify-between mb-8">
                <div>
                  <h2 className="text-3xl font-bold text-secondary-900 mb-2">Featured Startups</h2>
                  <p className="text-secondary-600">Innovative companies building the future</p>
                </div>
                <Link href="/startups">
                  <Button variant="ghost">
                    View All
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Button>
                </Link>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {startups.map((startup) => (
                  <StartupCard key={startup.id} startup={startup} />
                ))}
              </div>
            </div>
          </section>
        )}

        {/* CTA Section */}
        <section className="py-16 bg-primary-600">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8 text-center">
            <h2 className="text-3xl font-bold text-white mb-4">Ready to get started?</h2>
            <p className="text-primary-100 mb-8 max-w-2xl mx-auto">
              Join thousands of professionals finding their dream roles at innovative startups.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Link href="/register">
                <Button size="lg" variant="secondary" className="w-full sm:w-auto">
                  Create Account
                </Button>
              </Link>
              <Link href="/jobs">
                <Button size="lg" variant="outline" className="w-full sm:w-auto border-white text-white hover:bg-white/10">
                  Browse Jobs
                </Button>
              </Link>
            </div>
          </div>
        </section>
      </main>
      <Footer />
    </div>
  )
}

