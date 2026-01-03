import { MetadataRoute } from 'next'
import axios from 'axios'

const baseUrl = process.env.NEXT_PUBLIC_SITE_URL || process.env.NEXT_PUBLIC_APP_URL || 'https://joinus.com'

// Server-side API client for fetching data
async function fetchJobs() {
  try {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
    const response = await axios.get(`${apiUrl}/api/v1/jobs`, {
      params: {
        page_size: 1000, // Fetch a large number of jobs
        page: 1,
      },
      timeout: 10000, // 10 second timeout
    })
    return response.data?.data || []
  } catch (error) {
    console.error('Error fetching jobs for sitemap:', error)
    return []
  }
}

async function fetchStartups() {
  try {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
    const response = await axios.get(`${apiUrl}/api/v1/startups`, {
      params: {
        page_size: 1000, // Fetch a large number of startups
        page: 1,
      },
      timeout: 10000, // 10 second timeout
    })
    return response.data?.data || []
  } catch (error) {
    console.error('Error fetching startups for sitemap:', error)
    return []
  }
}

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  // Static pages
  const staticPages: MetadataRoute.Sitemap = [
    {
      url: baseUrl,
      lastModified: new Date(),
      changeFrequency: 'daily',
      priority: 1,
    },
    {
      url: `${baseUrl}/jobs`,
      lastModified: new Date(),
      changeFrequency: 'hourly',
      priority: 0.9,
    },
    {
      url: `${baseUrl}/startups`,
      lastModified: new Date(),
      changeFrequency: 'hourly',
      priority: 0.9,
    },
    {
      url: `${baseUrl}/about`,
      lastModified: new Date(),
      changeFrequency: 'monthly',
      priority: 0.7,
    },
    {
      url: `${baseUrl}/contact`,
      lastModified: new Date(),
      changeFrequency: 'monthly',
      priority: 0.7,
    },
  ]

  // Fetch dynamic content
  const [jobs, startups] = await Promise.all([
    fetchJobs(),
    fetchStartups(),
  ])

  // Job pages
  const jobPages: MetadataRoute.Sitemap = jobs.map((job: any) => ({
    url: `${baseUrl}/jobs/${job.id}`,
    lastModified: job.updated_at ? new Date(job.updated_at) : new Date(job.created_at || Date.now()),
    changeFrequency: 'weekly' as const,
    priority: 0.8,
  }))

  // Startup pages
  const startupPages: MetadataRoute.Sitemap = startups
    .filter((startup: any) => startup.slug) // Only include startups with slugs
    .map((startup: any) => ({
      url: `${baseUrl}/startups/${startup.slug}`,
      lastModified: startup.updated_at ? new Date(startup.updated_at) : new Date(startup.created_at || Date.now()),
      changeFrequency: 'weekly' as const,
      priority: 0.8,
    }))

  return [...staticPages, ...jobPages, ...startupPages]
}

