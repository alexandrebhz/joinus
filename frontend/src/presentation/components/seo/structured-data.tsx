import { Job } from '@/domain/entities/job.entity'

export function JobPostingStructuredData({ jobs }: { jobs: Job[] }) {
  const structuredData = {
    '@context': 'https://schema.org',
    '@type': 'ItemList',
    itemListElement: jobs.map((job, index) => ({
      '@type': 'ListItem',
      position: index + 1,
      item: {
        '@type': 'JobPosting',
        title: job.title,
        description: job.description?.substring(0, 500) || '', // Limit description length
        identifier: {
          '@type': 'PropertyValue',
          name: 'JoinUs',
          value: job.id,
        },
        ...(job.createdAt && !isNaN(new Date(job.createdAt).getTime()) && {
          datePosted: job.createdAt,
        }),
        employmentType: job.jobType?.replace('_', ' ') || 'FULL_TIME',
        hiringOrganization: {
          '@type': 'Organization',
          name: job.startupName || 'Startup',
        },
        jobLocation: {
          '@type': 'Place',
          address: {
            '@type': 'PostalAddress',
            addressLocality: job.city || '',
            addressCountry: job.country,
          },
        },
        ...(job.salaryMin && job.salaryMax && {
          baseSalary: {
            '@type': 'MonetaryAmount',
            currency: job.currency,
            value: {
              '@type': 'QuantitativeValue',
              minValue: job.salaryMin,
              maxValue: job.salaryMax,
              unitText: 'YEAR',
            },
          },
        }),
        ...(job.applicationUrl && {
          url: job.applicationUrl,
        }),
      },
    })),
  }

  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(structuredData) }}
    />
  )
}

export function OrganizationStructuredData() {
  const structuredData = {
    '@context': 'https://schema.org',
    '@type': 'Organization',
    name: 'JoinUs',
    url: 'https://joinus.ie',
    logo: 'https://joinus.ie/logo.png',
    description: 'Startup job board connecting talented professionals with innovative startups',
    sameAs: [
      'https://twitter.com/joinus',
      'https://linkedin.com/company/joinus',
    ],
    contactPoint: {
      '@type': 'ContactPoint',
      contactType: 'Customer Service',
      email: 'hello@joinus.ie',
    },
  }

  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(structuredData) }}
    />
  )
}

export function WebSiteStructuredData() {
  const structuredData = {
    '@context': 'https://schema.org',
    '@type': 'WebSite',
    name: 'JoinUs',
    url: 'https://joinus.ie',
    potentialAction: {
      '@type': 'SearchAction',
      target: {
        '@type': 'EntryPoint',
        urlTemplate: 'https://joinus.ie/jobs?search={search_term_string}',
      },
      'query-input': 'required name=search_term_string',
    },
  }

  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(structuredData) }}
    />
  )
}

