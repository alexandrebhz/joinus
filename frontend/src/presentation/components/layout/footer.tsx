import Link from 'next/link'

const jobSeekerLinks = [
  { href: '/jobs', label: 'Browse Jobs' },
  { href: '/jobs?location_type=remote', label: 'Remote Jobs' },
  { href: '/jobs?location_type=hybrid', label: 'Hybrid Jobs' },
  { href: '/jobs?job_type=full_time', label: 'Full-time Jobs' },
  { href: '/startups', label: 'Explore Startups' },
]

const startupLinks = [
  { href: '/register', label: 'Post a Job' },
  { href: '/startups', label: 'Browse Startups' },
  { href: '/about', label: 'Why JoinUs' },
]

const companyLinks = [
  { href: '/about', label: 'About' },
  { href: '/contact', label: 'Contact' },
  { href: '/jobs', label: 'Open roles' },
]

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
              {jobSeekerLinks.map((link) => (
                <li key={link.href + link.label}>
                  <Link href={link.href} className="hover:text-primary-600 transition-colors">
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h3 className="font-semibold text-secondary-900 mb-4">For Startups</h3>
            <ul className="space-y-2 text-sm text-secondary-600">
              {startupLinks.map((link) => (
                <li key={link.href + link.label}>
                  <Link href={link.href} className="hover:text-primary-600 transition-colors">
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h3 className="font-semibold text-secondary-900 mb-4">Company</h3>
            <ul className="space-y-2 text-sm text-secondary-600">
              {companyLinks.map((link) => (
                <li key={link.href + link.label}>
                  <Link href={link.href} className="hover:text-primary-600 transition-colors">
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <div className="mt-8 pt-8 border-t border-secondary-200 text-center text-sm text-secondary-600 space-y-2">
          <p>&copy; {new Date().getFullYear()} JoinUs. All rights reserved.</p>
          <p>
            Developed by{' '}
            <a
              href="https://zaptcode.com"
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary-600 hover:text-primary-700 transition-colors"
            >
              ZaptCode
            </a>
          </p>
        </div>
      </div>
    </footer>
  )
}
