'use client'

import { useRouter, useSearchParams } from 'next/navigation'
import { useState, useEffect, useCallback } from 'react'
import { Input } from '../ui/input'
import { Button } from '../ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card'
import { Search, X, Filter } from 'lucide-react'
import { JobType, LocationType } from '@/domain/entities/job.entity'

interface JobFiltersProps {
  className?: string
}

export function JobFilters({ className }: JobFiltersProps) {
  const router = useRouter()
  const searchParams = useSearchParams()
  
  const [search, setSearch] = useState(searchParams.get('search') || '')
  const [jobType, setJobType] = useState<JobType | ''>(searchParams.get('job_type') as JobType || '')
  const [locationType, setLocationType] = useState<LocationType | ''>(searchParams.get('location_type') as LocationType || '')

  // Debounced search to avoid too many URL updates
  const [searchDebounce, setSearchDebounce] = useState<NodeJS.Timeout | null>(null)

  const applyFilters = useCallback(() => {
    const params = new URLSearchParams()
    
    if (search.trim()) params.set('search', search.trim())
    if (jobType) params.set('job_type', jobType)
    if (locationType) params.set('location_type', locationType)
    
    const queryString = params.toString()
    router.push(queryString ? `/jobs?${queryString}` : '/jobs')
  }, [search, jobType, locationType, router])

  // Auto-apply filters when job type or location type changes
  useEffect(() => {
    applyFilters()
  }, [jobType, locationType, applyFilters])

  // Debounced search
  useEffect(() => {
    if (searchDebounce) {
      clearTimeout(searchDebounce)
    }
    
    const timer = setTimeout(() => {
      applyFilters()
    }, 500) // Wait 500ms after user stops typing
    
    setSearchDebounce(timer)
    
    return () => {
      if (timer) clearTimeout(timer)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [search, applyFilters])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchDebounce) {
      clearTimeout(searchDebounce)
    }
    applyFilters()
  }

  const clearFilters = () => {
    setSearch('')
    setJobType('')
    setLocationType('')
    router.push('/jobs')
  }

  const hasActiveFilters = search || jobType || locationType

  return (
    <Card className={`${className} sticky top-4`}>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center text-lg">
            <Filter className="h-5 w-5 mr-2 text-primary-600" />
            Filters
            {hasActiveFilters && (
              <span className="ml-2 h-5 w-5 rounded-full bg-primary-600 text-white text-xs flex items-center justify-center">
                {[search, jobType, locationType].filter(Boolean).length}
              </span>
            )}
          </CardTitle>
          {hasActiveFilters && (
            <Button
              variant="ghost"
              size="sm"
              onClick={clearFilters}
              className="text-xs"
            >
              <X className="h-3 w-3 mr-1" />
              Clear
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Search */}
          <div>
            <label className="block text-sm font-medium text-secondary-700 mb-2">
              Search
            </label>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-secondary-400" />
              <Input
                type="text"
                placeholder="Job title, company, keywords..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-10"
              />
            </div>
          </div>

          {/* Job Type */}
          <div>
            <label className="block text-sm font-medium text-secondary-700 mb-2">
              Job Type
            </label>
            <div className="space-y-2">
              {(['full_time', 'part_time', 'contract', 'internship'] as JobType[]).map((type) => (
                <label
                  key={type}
                  className={`flex items-center space-x-2 cursor-pointer p-2 rounded-lg transition-colors ${
                    jobType === type
                      ? 'bg-primary-50 border border-primary-200'
                      : 'hover:bg-secondary-50'
                  }`}
                >
                  <input
                    type="radio"
                    name="job_type"
                    value={type}
                    checked={jobType === type}
                    onChange={(e) => setJobType(e.target.value as JobType)}
                    className="text-primary-600 focus:ring-primary-500"
                  />
                  <span className={`text-sm capitalize ${
                    jobType === type ? 'text-primary-700 font-medium' : 'text-secondary-700'
                  }`}>
                    {type.replace('_', ' ')}
                  </span>
                </label>
              ))}
              <label className="flex items-center space-x-2 cursor-pointer hover:bg-secondary-50 p-2 rounded-lg transition-colors">
                <input
                  type="radio"
                  name="job_type"
                  value=""
                  checked={jobType === ''}
                  onChange={() => setJobType('')}
                  className="text-primary-600 focus:ring-primary-500"
                />
                <span className="text-sm text-secondary-700">All Types</span>
              </label>
            </div>
          </div>

          {/* Location Type */}
          <div>
            <label className="block text-sm font-medium text-secondary-700 mb-2">
              Location
            </label>
            <div className="space-y-2">
              {(['remote', 'hybrid', 'onsite'] as LocationType[]).map((type) => (
                <label
                  key={type}
                  className={`flex items-center space-x-2 cursor-pointer p-2 rounded-lg transition-colors ${
                    locationType === type
                      ? 'bg-primary-50 border border-primary-200'
                      : 'hover:bg-secondary-50'
                  }`}
                >
                  <input
                    type="radio"
                    name="location_type"
                    value={type}
                    checked={locationType === type}
                    onChange={(e) => setLocationType(e.target.value as LocationType)}
                    className="text-primary-600 focus:ring-primary-500"
                  />
                  <span className={`text-sm capitalize ${
                    locationType === type ? 'text-primary-700 font-medium' : 'text-secondary-700'
                  }`}>
                    {type}
                  </span>
                </label>
              ))}
              <label className="flex items-center space-x-2 cursor-pointer hover:bg-secondary-50 p-2 rounded-lg transition-colors">
                <input
                  type="radio"
                  name="location_type"
                  value=""
                  checked={locationType === ''}
                  onChange={() => setLocationType('')}
                  className="text-primary-600 focus:ring-primary-500"
                />
                <span className="text-sm text-secondary-700">All Locations</span>
              </label>
            </div>
          </div>

          {/* Apply Button (for search) */}
          <Button type="submit" className="w-full" size="lg">
            <Search className="h-4 w-4 mr-2" />
            Search
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}

