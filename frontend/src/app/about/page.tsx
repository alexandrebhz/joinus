import { Metadata } from 'next'
import Link from 'next/link'
import { Header } from '@/presentation/components/layout/header'
import { Footer } from '@/presentation/components/layout/footer'
import { Button } from '@/presentation/components/ui/button'
import { Briefcase, Building2, TrendingUp, Users, Zap, Target, ArrowRight, Heart, Shield, Rocket } from 'lucide-react'

export const metadata: Metadata = {
  title: 'About Us | JoinUs Job Board',
  description: 'Learn about JoinUs - the premier job board connecting talented professionals with innovative startups. Discover our mission, values, and commitment to helping you find your dream job.',
  openGraph: {
    title: 'About Us | JoinUs Job Board',
    description: 'Learn about JoinUs - connecting talented professionals with innovative startups.',
  },
}

export default function AboutPage() {
  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1">
        {/* Hero Section */}
        <section className="relative overflow-hidden bg-gradient-to-br from-primary-50 via-white to-secondary-50 py-20 sm:py-32">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="max-w-4xl mx-auto text-center">
              <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-secondary-900 mb-6 leading-tight">
                About <span className="text-primary-600">JoinUs</span>
              </h1>
              <p className="text-xl sm:text-2xl text-secondary-700 mb-4 font-medium">
                Connecting talented professionals with innovative startups
              </p>
              <p className="text-lg text-secondary-600 mb-8">
                We're on a mission to make it easier than ever to discover and land your dream job at fast-growing startups.
              </p>
            </div>
          </div>
        </section>

        {/* Mission Section */}
        <section className="py-20 bg-white">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="max-w-3xl mx-auto">
              <div className="text-center mb-12">
                <h2 className="text-3xl sm:text-4xl font-bold text-secondary-900 mb-4">
                  Our Mission
                </h2>
                <p className="text-lg text-secondary-600 leading-relaxed">
                  JoinUs was born from a simple belief: finding your next career opportunity shouldn't be complicated. 
                  We've created a platform that brings together the best of both worlds—talented professionals seeking 
                  meaningful work and innovative startups looking for exceptional team members.
                </p>
              </div>

              <div className="prose prose-lg max-w-none text-secondary-700 space-y-6">
                <p>
                  In today's fast-paced tech landscape, startups are reshaping industries and creating new opportunities 
                  every day. Yet, discovering these opportunities and connecting with the right companies can be overwhelming. 
                  That's where JoinUs comes in.
                </p>
                <p>
                  We curate job listings from verified startups, ensuring every opportunity is genuine and every company 
                  is actively hiring. We provide deep insights into company culture, team dynamics, and growth potential 
                  so you can make informed decisions about your next career move.
                </p>
                <p>
                  Whether you're a software engineer looking to build the next big thing, a product manager ready to 
                  shape innovative products, a designer wanting to create beautiful experiences, or a marketer eager to 
                  tell compelling stories—JoinUs is your gateway to startup opportunities.
                </p>
              </div>
            </div>
          </div>
        </section>

        {/* What Sets Us Apart */}
        <section className="py-20 bg-secondary-50">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="text-center mb-12">
              <h2 className="text-3xl sm:text-4xl font-bold text-secondary-900 mb-4">
                Why Job Seekers Choose JoinUs
              </h2>
              <p className="text-lg text-secondary-600 max-w-2xl mx-auto">
                We make it easy to discover and apply to the best startup jobs. Here's what sets us apart:
              </p>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-6xl mx-auto">
              <div className="bg-white p-8 rounded-xl shadow-sm hover:shadow-lg transition-shadow">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary-100 text-primary-600 mb-4">
                  <Briefcase className="h-8 w-8" />
                </div>
                <h3 className="text-xl font-semibold text-secondary-900 mb-3">Curated Job Listings</h3>
                <p className="text-secondary-600 leading-relaxed">
                  Every job is hand-verified from real startups. No spam, no scams—just genuine opportunities 
                  from companies actively hiring. We ensure quality over quantity, so you can focus on opportunities 
                  that truly matter.
                </p>
              </div>
              <div className="bg-white p-8 rounded-xl shadow-sm hover:shadow-lg transition-shadow">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary-100 text-primary-600 mb-4">
                  <Building2 className="h-8 w-8" />
                </div>
                <h3 className="text-xl font-semibold text-secondary-900 mb-3">Deep Company Insights</h3>
                <p className="text-secondary-600 leading-relaxed">
                  Get the full picture before you apply. Browse company profiles, team size, funding stage, and 
                  culture to find your perfect match. We believe informed decisions lead to better career outcomes.
                </p>
              </div>
              <div className="bg-white p-8 rounded-xl shadow-sm hover:shadow-lg transition-shadow">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary-100 text-primary-600 mb-4">
                  <TrendingUp className="h-8 w-8" />
                </div>
                <h3 className="text-xl font-semibold text-secondary-900 mb-3">Early-Stage Opportunities</h3>
                <p className="text-secondary-600 leading-relaxed">
                  Join startups at the ground floor. Get equity, make an impact, and accelerate your career at 
                  companies that are changing industries. Be part of something bigger from day one.
                </p>
              </div>
            </div>
          </div>
        </section>

        {/* Values Section */}
        <section className="py-20 bg-white">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="max-w-4xl mx-auto">
              <div className="text-center mb-12">
                <h2 className="text-3xl sm:text-4xl font-bold text-secondary-900 mb-4">
                  Our Values
                </h2>
                <p className="text-lg text-secondary-600">
                  The principles that guide everything we do
                </p>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                <div className="flex items-start space-x-4">
                  <div className="flex-shrink-0">
                    <div className="w-12 h-12 rounded-lg bg-primary-100 flex items-center justify-center">
                      <Shield className="h-6 w-6 text-primary-600" />
                    </div>
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold text-secondary-900 mb-2">Trust & Transparency</h3>
                    <p className="text-secondary-600">
                      We verify every startup and job listing to ensure authenticity. You can trust that every 
                      opportunity is real and every company is legitimate.
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-4">
                  <div className="flex-shrink-0">
                    <div className="w-12 h-12 rounded-lg bg-primary-100 flex items-center justify-center">
                      <Heart className="h-6 w-6 text-primary-600" />
                    </div>
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold text-secondary-900 mb-2">User-First Approach</h3>
                    <p className="text-secondary-600">
                      Your success is our success. We're constantly improving our platform based on feedback 
                      from job seekers and startups alike.
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-4">
                  <div className="flex-shrink-0">
                    <div className="w-12 h-12 rounded-lg bg-primary-100 flex items-center justify-center">
                      <Zap className="h-6 w-6 text-primary-600" />
                    </div>
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold text-secondary-900 mb-2">Quality Over Quantity</h3>
                    <p className="text-secondary-600">
                      We'd rather show you 100 great opportunities than 1000 mediocre ones. Every listing 
                      is carefully curated to match your career goals.
                    </p>
                  </div>
                </div>
                <div className="flex items-start space-x-4">
                  <div className="flex-shrink-0">
                    <div className="w-12 h-12 rounded-lg bg-primary-100 flex items-center justify-center">
                      <Rocket className="h-6 w-6 text-primary-600" />
                    </div>
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold text-secondary-900 mb-2">Innovation & Growth</h3>
                    <p className="text-secondary-600">
                      We're always evolving. As the startup ecosystem grows, so do we—introducing new features 
                      and tools to make your job search even better.
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* What We Offer */}
        <section className="py-20 bg-secondary-50">
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="max-w-4xl mx-auto">
              <div className="text-center mb-12">
                <h2 className="text-3xl sm:text-4xl font-bold text-secondary-900 mb-4">
                  What We Offer
                </h2>
                <p className="text-lg text-secondary-600">
                  Everything you need to find your next opportunity
                </p>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="bg-white p-6 rounded-lg border border-secondary-200">
                  <div className="flex items-center mb-3">
                    <Target className="h-5 w-5 text-primary-600 mr-2" />
                    <h3 className="font-semibold text-secondary-900">Verified Startups Only</h3>
                  </div>
                  <p className="text-sm text-secondary-600">
                    Every company on our platform is verified and actively hiring. No fake listings, no scams.
                  </p>
                </div>
                <div className="bg-white p-6 rounded-lg border border-secondary-200">
                  <div className="flex items-center mb-3">
                    <Zap className="h-5 w-5 text-primary-600 mr-2" />
                    <h3 className="font-semibold text-secondary-900">New Jobs Added Daily</h3>
                  </div>
                  <p className="text-sm text-secondary-600">
                    Fresh opportunities from innovative startups are added every day. Never miss a chance.
                  </p>
                </div>
                <div className="bg-white p-6 rounded-lg border border-secondary-200">
                  <div className="flex items-center mb-3">
                    <Building2 className="h-5 w-5 text-primary-600 mr-2" />
                    <h3 className="font-semibold text-secondary-900">100% Free to Apply</h3>
                  </div>
                  <p className="text-sm text-secondary-600">
                    Browse jobs, create your profile, and apply—all completely free. No hidden fees, ever.
                  </p>
                </div>
                <div className="bg-white p-6 rounded-lg border border-secondary-200">
                  <div className="flex items-center mb-3">
                    <Users className="h-5 w-5 text-primary-600 mr-2" />
                    <h3 className="font-semibold text-secondary-900">Direct Applications</h3>
                  </div>
                  <p className="text-sm text-secondary-600">
                    Apply directly to startups. No middlemen, no delays—just you and the company.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </section>

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
