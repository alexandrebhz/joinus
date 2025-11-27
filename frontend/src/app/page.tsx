import Link from 'next/link'
import { Metadata } from 'next'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { Button } from '@/presentation/components/ui/button'
import { apiClient } from '@/infrastructure/api/api-client'
import { JobCard } from '@/presentation/components/job/job-card'
import { StartupCard } from '@/presentation/components/startup/startup-card'
import { JobPostingStructuredData, OrganizationStructuredData, WebSiteStructuredData } from '@/presentation/components/seo/structured-data'
import { ArrowRight, Search, Briefcase, Building2, TrendingUp, Users, Zap, Target } from 'lucide-react'

export const metadata: Metadata = {
  title: 'Find Your Dream Job at Innovative Startups | JoinUs Job Board',
  description: 'Discover 1000+ tech jobs at fast-growing startups. Browse remote, hybrid, and onsite software engineering, product, design, and marketing roles. Join thousands of professionals finding their next career opportunity.',
  openGraph: {
    title: 'Find Your Dream Job at Innovative Startups | JoinUs Job Board',
    description: 'Discover 1000+ tech jobs at fast-growing startups. Browse remote, hybrid, and onsite software engineering, product, design, and marketing roles.',
  },
}

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
    <>
      <OrganizationStructuredData />
      <WebSiteStructuredData />
      {jobs.length > 0 && <JobPostingStructuredData jobs={jobs} />}
      
      <div className="min-h-screen flex flex-col">
        <Header />
        <main className="flex-1">
          {/* Hero Section */}
          <section className="relative overflow-hidden bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-20 sm:py-32">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8">
              <div className="max-w-4xl mx-auto text-center">
                <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-secondary-900 mb-6 leading-tight">
                  Land Your Dream Job at
                  <span className="text-primary-600"> Fast-Growing Startups</span>
                </h1>
                <p className="text-xl sm:text-2xl text-secondary-700 mb-4 font-medium">
                  Discover 1000+ tech jobs from innovative companies. Remote, hybrid, and onsite opportunities for software engineers, product managers, designers, and more.
                </p>
                <p className="text-lg text-secondary-600 mb-8">
                  Join thousands of professionals who found their next career move through our curated startup job board.
                </p>
                <div className="flex flex-col sm:flex-row gap-4 justify-center mb-8">
                  <Link href="/jobs">
                    <Button size="lg" className="w-full sm:w-auto text-lg px-8 py-6">
                      <Search className="mr-2 h-5 w-5" />
                      Browse All Jobs
                    </Button>
                  </Link>
                  <Link href="/register">
                    <Button variant="outline" size="lg" className="w-full sm:w-auto text-lg px-8 py-6">
                      <Users className="mr-2 h-5 w-5" />
                      Create Free Account
                    </Button>
                  </Link>
                </div>
                <div className="flex flex-wrap justify-center gap-6 text-sm text-secondary-600">
                  <div className="flex items-center">
                    <Zap className="h-4 w-4 mr-2 text-primary-600" />
                    <span>New jobs added daily</span>
                  </div>
                  <div className="flex items-center">
                    <Target className="h-4 w-4 mr-2 text-primary-600" />
                    <span>Verified startups only</span>
                  </div>
                  <div className="flex items-center">
                    <Building2 className="h-4 w-4 mr-2 text-primary-600" />
                    <span>100% free to apply</span>
                  </div>
                </div>
              </div>
            </div>
          </section>

        {/* Features Section */}
        <section className="py-20 bg-white">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="text-center mb-12">
              <h2 className="text-3xl sm:text-4xl font-bold text-secondary-900 mb-4">
                Why Job Seekers Choose JoinUs
              </h2>
              <p className="text-lg text-secondary-600 max-w-2xl mx-auto">
                We make it easy to discover and apply to the best startup jobs. Here's what sets us apart:
              </p>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <div className="text-center p-6 rounded-xl hover:bg-secondary-50 transition-colors">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary-100 text-primary-600 mb-4">
                  <Briefcase className="h-8 w-8" />
                </div>
                <h3 className="text-xl font-semibold text-secondary-900 mb-3">Curated Job Listings</h3>
                <p className="text-secondary-600 leading-relaxed">
                  Every job is hand-verified from real startups. No spam, no scams—just genuine opportunities from companies actively hiring.
                </p>
              </div>
              <div className="text-center p-6 rounded-xl hover:bg-secondary-50 transition-colors">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary-100 text-primary-600 mb-4">
                  <Building2 className="h-8 w-8" />
                </div>
                <h3 className="text-xl font-semibold text-secondary-900 mb-3">Deep Company Insights</h3>
                <p className="text-secondary-600 leading-relaxed">
                  Get the full picture before you apply. Browse company profiles, team size, funding stage, and culture to find your perfect match.
                </p>
              </div>
              <div className="text-center p-6 rounded-xl hover:bg-secondary-50 transition-colors">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary-100 text-primary-600 mb-4">
                  <TrendingUp className="h-8 w-8" />
                </div>
                <h3 className="text-xl font-semibold text-secondary-900 mb-3">Early-Stage Opportunities</h3>
                <p className="text-secondary-600 leading-relaxed">
                  Join startups at the ground floor. Get equity, make an impact, and accelerate your career at companies that are changing industries.
                </p>
              </div>
            </div>
          </div>
        </section>

        {/* Featured Jobs */}
        {jobs.length > 0 && (
          <section className="py-20 bg-secondary-50">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8">
              <div className="flex items-center justify-between mb-10">
                <div>
                  <h2 className="text-3xl sm:text-4xl font-bold text-secondary-900 mb-3">
                    Latest Startup Jobs
                  </h2>
                  <p className="text-lg text-secondary-600">
                    Fresh opportunities from innovative companies. Apply today and start your next chapter.
                  </p>
                </div>
                <Link href="/jobs">
                  <Button variant="ghost" size="lg">
                    View All Jobs
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Button>
                </Link>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
                {jobs.map((job) => (
                  <JobCard key={job.id} job={job} />
                ))}
              </div>
              <div className="text-center">
                <Link href="/jobs">
                  <Button size="lg" variant="outline">
                    Explore All {jobs.length > 0 ? '1000+' : ''} Jobs
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </div>
          </section>
        )}

        {/* Featured Startups */}
        {startups.length > 0 && (
          <section className="py-20 bg-white">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8">
              <div className="flex items-center justify-between mb-10">
                <div>
                  <h2 className="text-3xl sm:text-4xl font-bold text-secondary-900 mb-3">
                    Top Startups Hiring Now
                  </h2>
                  <p className="text-lg text-secondary-600">
                    Discover innovative companies that are reshaping industries and building the future of technology.
                  </p>
                </div>
                <Link href="/startups">
                  <Button variant="ghost" size="lg">
                    Browse All Startups
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
        <section className="py-20 bg-gradient-to-br from-primary-600 to-primary-700">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8 text-center">
            <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
              Ready to Find Your Next Opportunity?
            </h2>
            <p className="text-xl text-primary-100 mb-6 max-w-2xl mx-auto">
              Join thousands of software engineers, product managers, designers, and marketers who found their dream jobs through JoinUs.
            </p>
            <p className="text-primary-200 mb-10 max-w-xl mx-auto">
              Create your free account in 30 seconds. No credit card required. Start applying to top startup jobs today.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center mb-6">
              <Link href="/register">
                <Button size="lg" variant="secondary" className="w-full sm:w-auto text-lg px-8 py-6">
                  Create Free Account
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
              </Link>
              <Link href="/jobs">
                <Button size="lg" variant="outline" className="w-full sm:w-auto text-lg px-8 py-6 border-2 border-white text-white hover:bg-white/10">
                  Browse All Jobs
                </Button>
              </Link>
            </div>
            <p className="text-sm text-primary-200">
              ✓ Free forever • ✓ No spam • ✓ Direct applications
            </p>
          </div>
        </section>
      </main>
      <Footer />
    </div>
  )
}

