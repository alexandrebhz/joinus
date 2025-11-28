'use client'

import { useRouter, useSearchParams } from 'next/navigation'
import { ArrowUpDown } from 'lucide-react'

export function JobSortControls() {
  const router = useRouter()
  const searchParams = useSearchParams()

  const orderBy = searchParams.get('order_by') || 'created_at'
  const orderDir = (searchParams.get('order_dir') as 'ASC' | 'DESC') || 'DESC'

  const updateSort = (newOrderBy: string, newOrderDir: 'ASC' | 'DESC' = 'DESC') => {
    const params = new URLSearchParams(searchParams.toString())
    
    if (newOrderBy === 'created_at' && newOrderDir === 'DESC') {
      // Default sort - remove params
      params.delete('order_by')
      params.delete('order_dir')
    } else {
      params.set('order_by', newOrderBy)
      params.set('order_dir', newOrderDir)
    }
    
    // Reset to page 1 when sorting changes
    params.delete('page')
    
    router.push(`/jobs?${params.toString()}`)
  }

  const sortOptions = [
    { value: 'created_at', label: 'Date Posted', dir: 'DESC' },
    { value: 'title', label: 'Job Title', dir: 'ASC' },
    { value: 'salary_max', label: 'Salary (High to Low)', dir: 'DESC' },
    { value: 'salary_min', label: 'Salary (Low to High)', dir: 'ASC' },
  ]

  const currentSort = sortOptions.find(opt => opt.value === orderBy) || sortOptions[0]

  return (
    <div className="flex items-center justify-between bg-secondary-50 rounded-lg p-4">
      <div className="flex items-center gap-2">
        <ArrowUpDown className="h-4 w-4 text-secondary-600" />
        <span className="text-sm font-medium text-secondary-700">Sort by:</span>
      </div>
      
      <div className="flex items-center gap-2">
        {sortOptions.map((option) => {
          const isActive = orderBy === option.value && orderDir === option.dir
          
          return (
            <button
              key={option.value}
              onClick={() => updateSort(option.value, option.dir)}
              className={`px-3 py-1.5 text-sm rounded-md transition-colors ${
                isActive
                  ? 'bg-primary-500 text-white font-medium'
                  : 'bg-white text-secondary-700 hover:bg-secondary-100 border border-secondary-300'
              }`}
            >
              {option.label}
            </button>
          )
        })}
      </div>
    </div>
  )
}

