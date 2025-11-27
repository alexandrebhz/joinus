'use client'

import { useState } from 'react'
import { Button } from '../ui/button'
import { JobFilters } from './job-filters'
import { Filter, X } from 'lucide-react'

export function JobFiltersMobile() {
  const [isOpen, setIsOpen] = useState(false)

  return (
    <>
      <div className="md:hidden mb-4">
        <Button
          variant="outline"
          onClick={() => setIsOpen(!isOpen)}
          className="w-full"
        >
          <Filter className="h-4 w-4 mr-2" />
          {isOpen ? 'Hide Filters' : 'Show Filters'}
        </Button>
      </div>

      {isOpen && (
        <div className="md:hidden mb-6">
          <div className="relative">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setIsOpen(false)}
              className="absolute top-2 right-2 z-10"
            >
              <X className="h-4 w-4" />
            </Button>
            <JobFilters />
          </div>
        </div>
      )}
    </>
  )
}

