'use client'

import { useRouter, useSearchParams } from 'next/navigation'
import { useState, useEffect, useCallback } from 'react'
import { Input } from '../ui/input'
import { Button } from '../ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card'
import { Search, X, Filter, DollarSign, MapPin, ArrowUpDown } from 'lucide-react'
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
  const [country, setCountry] = useState(searchParams.get('country') || '')
  const [city, setCity] = useState(searchParams.get('city') || '')
  const [salaryMin, setSalaryMin] = useState(searchParams.get('salary_min') || '')
  const [salaryMax, setSalaryMax] = useState(searchParams.get('salary_max') || '')
  const [currency, setCurrency] = useState(searchParams.get('currency') || '')
  const [orderBy, setOrderBy] = useState(searchParams.get('order_by') || 'created_at')
  const [orderDir, setOrderDir] = useState<'ASC' | 'DESC'>(searchParams.get('order_dir') as 'ASC' | 'DESC' || 'DESC')

  // Debounced search to avoid too many URL updates
  const [searchDebounce, setSearchDebounce] = useState<NodeJS.Timeout | null>(null)
  const [locationDebounce, setLocationDebounce] = useState<NodeJS.Timeout | null>(null)

  const applyFilters = useCallback(() => {
    const params = new URLSearchParams()
    
    if (search.trim()) params.set('search', search.trim())
    if (jobType) params.set('job_type', jobType)
    if (locationType) params.set('location_type', locationType)
    if (country.trim()) params.set('country', country.trim())
    if (city.trim()) params.set('city', city.trim())
    if (salaryMin) params.set('salary_min', salaryMin)
    if (salaryMax) params.set('salary_max', salaryMax)
    if (currency) params.set('currency', currency)
    if (orderBy && orderBy !== 'created_at') params.set('order_by', orderBy)
    if (orderDir && orderDir !== 'DESC') params.set('order_dir', orderDir)
    
    const queryString = params.toString()
    router.push(queryString ? `/jobs?${queryString}` : '/jobs')
  }, [search, jobType, locationType, country, city, salaryMin, salaryMax, currency, orderBy, orderDir, router])

  // Auto-apply filters when dropdowns change
  useEffect(() => {
    applyFilters()
  }, [jobType, locationType, currency, orderBy, orderDir, applyFilters])

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

  // Debounced location filters (country, city)
  useEffect(() => {
    if (locationDebounce) {
      clearTimeout(locationDebounce)
    }
    
    const timer = setTimeout(() => {
      applyFilters()
    }, 500)
    
    setLocationDebounce(timer)
    
    return () => {
      if (timer) clearTimeout(timer)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [country, city, applyFilters])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchDebounce) {
      clearTimeout(searchDebounce)
    }
    if (locationDebounce) {
      clearTimeout(locationDebounce)
    }
    applyFilters()
  }

  const clearFilters = () => {
    setSearch('')
    setJobType('')
    setLocationType('')
    setCountry('')
    setCity('')
    setSalaryMin('')
    setSalaryMax('')
    setCurrency('')
    setOrderBy('created_at')
    setOrderDir('DESC')
    router.push('/jobs')
  }

  const activeFilters = [
    search,
    jobType,
    locationType,
    country,
    city,
    salaryMin,
    salaryMax,
    currency,
    orderBy !== 'created_at' ? orderBy : null,
    orderDir !== 'DESC' ? orderDir : null,
  ].filter(Boolean)
  
  const hasActiveFilters = activeFilters.length > 0

  return (
    <Card className={`${className} sticky top-4`}>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center text-lg">
            <Filter className="h-5 w-5 mr-2 text-primary-600" />
            Filters
            {hasActiveFilters && (
              <span className="ml-2 h-5 w-5 rounded-full bg-primary-600 text-white text-xs flex items-center justify-center">
                {activeFilters.length}
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
              Work Arrangement
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
                <span className="text-sm text-secondary-700">All Arrangements</span>
              </label>
            </div>
          </div>

          {/* Location Details */}
          <div className="border-t pt-4">
            <label className="block text-sm font-medium text-secondary-700 mb-2 flex items-center">
              <MapPin className="h-4 w-4 mr-1" />
              Location Details
            </label>
            <div className="space-y-3">
              <div>
                <Input
                  type="text"
                  placeholder="Country (e.g., Ireland, USA)"
                  value={country}
                  onChange={(e) => setCountry(e.target.value)}
                />
              </div>
              <div>
                <Input
                  type="text"
                  placeholder="City (e.g., Dublin, San Francisco)"
                  value={city}
                  onChange={(e) => setCity(e.target.value)}
                />
              </div>
            </div>
          </div>

          {/* Salary Range */}
          <div className="border-t pt-4">
            <label className="block text-sm font-medium text-secondary-700 mb-2 flex items-center">
              <DollarSign className="h-4 w-4 mr-1" />
              Salary Range
            </label>
            <div className="space-y-3">
              <div>
                <label className="block text-xs text-secondary-600 mb-1">Currency</label>
                <select
                  value={currency}
                  onChange={(e) => setCurrency(e.target.value)}
                  className="w-full px-3 py-2 border border-secondary-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 text-sm"
                >
                  <option value="">All Currencies</option>
                  <option value="EUR">EUR (€)</option>
                  <option value="USD">USD ($)</option>
                  <option value="GBP">GBP (£)</option>
                  <option value="CAD">CAD (C$)</option>
                  <option value="AUD">AUD (A$)</option>
                </select>
              </div>
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="block text-xs text-secondary-600 mb-1">Min</label>
                  <Input
                    type="number"
                    placeholder="Min"
                    value={salaryMin}
                    onChange={(e) => setSalaryMin(e.target.value)}
                    min="0"
                  />
                </div>
                <div>
                  <label className="block text-xs text-secondary-600 mb-1">Max</label>
                  <Input
                    type="number"
                    placeholder="Max"
                    value={salaryMax}
                    onChange={(e) => setSalaryMax(e.target.value)}
                    min="0"
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Sort Options */}
          <div className="border-t pt-4">
            <label className="block text-sm font-medium text-secondary-700 mb-2 flex items-center">
              <ArrowUpDown className="h-4 w-4 mr-1" />
              Sort By
            </label>
            <div className="space-y-2">
              <select
                value={orderBy}
                onChange={(e) => setOrderBy(e.target.value)}
                className="w-full px-3 py-2 border border-secondary-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 text-sm"
              >
                <option value="created_at">Date Posted</option>
                <option value="title">Job Title</option>
                <option value="salary_max">Salary (High to Low)</option>
                <option value="salary_min">Salary (Low to High)</option>
              </select>
              <div className="flex gap-2">
                <Button
                  type="button"
                  variant={orderDir === 'DESC' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setOrderDir('DESC')}
                  className="flex-1"
                >
                  Newest First
                </Button>
                <Button
                  type="button"
                  variant={orderDir === 'ASC' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setOrderDir('ASC')}
                  className="flex-1"
                >
                  Oldest First
                </Button>
              </div>
            </div>
          </div>

          {/* Apply Button */}
          <Button type="submit" className="w-full" size="lg">
            <Search className="h-4 w-4 mr-2" />
            Apply Filters
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}

