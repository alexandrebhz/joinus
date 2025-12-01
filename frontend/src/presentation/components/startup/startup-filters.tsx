'use client'

import { useRouter, useSearchParams } from 'next/navigation'
import { useState, useEffect, useCallback } from 'react'
import { Input } from '../ui/input'
import { Button } from '../ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card'
import { Search, X, Filter, Building2, MapPin, Users } from 'lucide-react'
import { StartupStatus } from '@/domain/entities/startup.entity'

interface StartupFiltersProps {
  className?: string
}

export function StartupFilters({ className }: StartupFiltersProps) {
  const router = useRouter()
  const searchParams = useSearchParams()
  
  const [search, setSearch] = useState(searchParams.get('search') || '')
  const [industry, setIndustry] = useState(searchParams.get('industry') || '')
  const [status, setStatus] = useState<StartupStatus | ''>(searchParams.get('status') as StartupStatus || '')
  const [companySize, setCompanySize] = useState(searchParams.get('company_size') || '')
  const [location, setLocation] = useState(searchParams.get('location') || '')

  // Debounced search to avoid too many URL updates
  const [searchDebounce, setSearchDebounce] = useState<NodeJS.Timeout | null>(null)
  const [textDebounce, setTextDebounce] = useState<NodeJS.Timeout | null>(null)

  const applyFilters = useCallback(() => {
    const params = new URLSearchParams()
    
    if (search.trim()) params.set('search', search.trim())
    if (industry.trim()) params.set('industry', industry.trim())
    if (status) params.set('status', status)
    if (companySize.trim()) params.set('company_size', companySize.trim())
    if (location.trim()) params.set('location', location.trim())
    
    const queryString = params.toString()
    router.push(queryString ? `/startups?${queryString}` : '/startups')
  }, [search, industry, status, companySize, location, router])

  // Auto-apply filters when dropdowns change
  useEffect(() => {
    applyFilters()
  }, [status, applyFilters])

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

  // Debounced text filters (industry, companySize, location)
  useEffect(() => {
    if (textDebounce) {
      clearTimeout(textDebounce)
    }
    
    const timer = setTimeout(() => {
      applyFilters()
    }, 500)
    
    setTextDebounce(timer)
    
    return () => {
      if (timer) clearTimeout(timer)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [industry, companySize, location, applyFilters])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchDebounce) {
      clearTimeout(searchDebounce)
    }
    if (textDebounce) {
      clearTimeout(textDebounce)
    }
    applyFilters()
  }

  const clearFilters = () => {
    setSearch('')
    setIndustry('')
    setStatus('')
    setCompanySize('')
    setLocation('')
    router.push('/startups')
  }

  const activeFilters = [
    search,
    industry,
    status,
    companySize,
    location,
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
                placeholder="Startup name, description, keywords..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-10"
              />
            </div>
          </div>

          {/* Status */}
          <div>
            <label className="block text-sm font-medium text-secondary-700 mb-2">
              Status
            </label>
            <div className="space-y-2">
              {(['active', 'inactive'] as StartupStatus[]).map((s) => (
                <label
                  key={s}
                  className={`flex items-center space-x-2 cursor-pointer p-2 rounded-lg transition-colors ${
                    status === s
                      ? 'bg-primary-50 border border-primary-200'
                      : 'hover:bg-secondary-50'
                  }`}
                >
                  <input
                    type="radio"
                    name="status"
                    value={s}
                    checked={status === s}
                    onChange={(e) => setStatus(e.target.value as StartupStatus)}
                    className="text-primary-600 focus:ring-primary-500"
                  />
                  <span className={`text-sm capitalize ${
                    status === s ? 'text-primary-700 font-medium' : 'text-secondary-700'
                  }`}>
                    {s}
                  </span>
                </label>
              ))}
              <label className="flex items-center space-x-2 cursor-pointer hover:bg-secondary-50 p-2 rounded-lg transition-colors">
                <input
                  type="radio"
                  name="status"
                  value=""
                  checked={status === ''}
                  onChange={() => setStatus('')}
                  className="text-primary-600 focus:ring-primary-500"
                />
                <span className="text-sm text-secondary-700">All Statuses</span>
              </label>
            </div>
          </div>

          {/* Industry */}
          <div className="border-t pt-4">
            <label className="block text-sm font-medium text-secondary-700 mb-2 flex items-center">
              <Building2 className="h-4 w-4 mr-1" />
              Industry
            </label>
            <Input
              type="text"
              placeholder="e.g., Technology, Healthcare, Finance"
              value={industry}
              onChange={(e) => setIndustry(e.target.value)}
            />
          </div>

          {/* Company Size */}
          <div>
            <label className="block text-sm font-medium text-secondary-700 mb-2 flex items-center">
              <Users className="h-4 w-4 mr-1" />
              Company Size
            </label>
            <Input
              type="text"
              placeholder="e.g., 1-10, 11-50, 51-200"
              value={companySize}
              onChange={(e) => setCompanySize(e.target.value)}
            />
          </div>

          {/* Location */}
          <div>
            <label className="block text-sm font-medium text-secondary-700 mb-2 flex items-center">
              <MapPin className="h-4 w-4 mr-1" />
              Location
            </label>
            <Input
              type="text"
              placeholder="e.g., Dublin, Ireland or San Francisco, USA"
              value={location}
              onChange={(e) => setLocation(e.target.value)}
            />
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


