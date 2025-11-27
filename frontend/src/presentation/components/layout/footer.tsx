import Link from 'next/link'

export function Footer() {
  return (
    <footer className="border-t border-secondary-200 bg-secondary-50">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div>
            <div className="flex items-center space-x-2 mb-4">
              <div className="h-6 w-6 rounded bg-primary-500 flex items-center justify-center">
                <span className="text-white font-bold text-sm">J</span>
              </div>
              <span className="text-lg font-bold text-secondary-900">JoinUs</span>
            </div>
            <p className="text-sm text-secondary-600">
              Connecting talented professionals with innovative startups.
            </p>
          </div>

          <div>
            <h3 className="font-semibold text-secondary-900 mb-4">For Job Seekers</h3>
            <ul className="space-y-2 text-sm text-secondary-600">
              <li>
                <Link href="/jobs" className="hover:text-primary-600 transition-colors">
                  Browse Jobs
                </Link>
              </li>
              <li>
                <Link href="/startups" className="hover:text-primary-600 transition-colors">
                  Explore Startups
                </Link>
              </li>
            </ul>
          </div>

          <div>
            <h3 className="font-semibold text-secondary-900 mb-4">For Startups</h3>
            <ul className="space-y-2 text-sm text-secondary-600">
              <li>
                <Link href="/register" className="hover:text-primary-600 transition-colors">
                  Post a Job
                </Link>
              </li>
              <li>
                <Link href="/startups" className="hover:text-primary-600 transition-colors">
                  Create Profile
                </Link>
              </li>
            </ul>
          </div>

          <div>
            <h3 className="font-semibold text-secondary-900 mb-4">Company</h3>
            <ul className="space-y-2 text-sm text-secondary-600">
              <li>
                <Link href="/about" className="hover:text-primary-600 transition-colors">
                  About
                </Link>
              </li>
              <li>
                <Link href="/contact" className="hover:text-primary-600 transition-colors">
                  Contact
                </Link>
              </li>
            </ul>
          </div>
        </div>

        <div className="mt-8 pt-8 border-t border-secondary-200 text-center text-sm text-secondary-600">
          <p>&copy; {new Date().getFullYear()} JoinUs. All rights reserved.</p>
        </div>
      </div>
    </footer>
  )
}

