'use client'

import Link from 'next/link'
import { Job } from '@/domain/entities/job.entity'
import { Card, CardContent, CardFooter } from '../ui/card'
import { Button } from '../ui/button'
import { MapPin, Briefcase, DollarSign, Clock } from 'lucide-react'
import { formatDistanceToNow } from 'date-fns'

interface JobCardProps {
  job: Job
}

export function JobCard({ job }: JobCardProps) {
  const formatSalary = () => {
    if (!job.salaryMin && !job.salaryMax) return null
    const min = job.salaryMin?.toLocaleString() || 'N/A'
    const max = job.salaryMax?.toLocaleString() || 'N/A'
    return `${job.currency} ${min} - ${max}`
  }

  const getLocationText = () => {
    if (job.locationType === 'remote') return 'Remote'
    if (job.locationType === 'hybrid') return `Hybrid • ${job.city || job.country}`
    return `${job.city ? `${job.city}, ` : ''}${job.country}`
  }

  return (
    <Card hover className="h-full flex flex-col">
      <CardContent className="flex-1 pt-6">
        <div className="flex items-start justify-between mb-4">
          <div className="flex-1">
            <Link href={`/jobs/${job.id}`}>
              <h3 className="text-lg font-semibold text-secondary-900 hover:text-primary-600 transition-colors mb-2">
                {job.title}
              </h3>
            </Link>
            {job.startupName && (
              <Link href={`/startups/${job.startupId}`}>
                <p className="text-sm text-primary-600 hover:text-primary-700 font-medium">
                  {job.startupName}
                </p>
              </Link>
            )}
          </div>
        </div>

        <p className="text-sm text-secondary-600 line-clamp-2 mb-4">{job.description}</p>

        <div className="flex flex-wrap gap-3 mb-4">
          <div className="flex items-center text-sm text-secondary-600">
            <MapPin className="h-4 w-4 mr-1" />
            {getLocationText()}
          </div>
          <div className="flex items-center text-sm text-secondary-600">
            <Briefcase className="h-4 w-4 mr-1" />
            {job.jobType.replace('_', ' ')}
          </div>
          {formatSalary() && (
            <div className="flex items-center text-sm text-secondary-600">
              <DollarSign className="h-4 w-4 mr-1" />
              {formatSalary()}
            </div>
          )}
        </div>
      </CardContent>

      <CardFooter className="flex items-center justify-between">
        <span className="text-xs text-secondary-500">
          <Clock className="h-3 w-3 inline mr-1" />
          {formatDistanceToNow(new Date(job.createdAt), { addSuffix: true })}
        </span>
        <Link href={`/jobs/${job.id}`}>
          <Button size="sm">View Details</Button>
        </Link>
      </CardFooter>
    </Card>
  )
}

