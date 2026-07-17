'use client'

import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { clsx } from 'clsx'
import { PaginationMeta } from '@/domain/value-objects/api-response.vo'

interface PaginationProps {
  meta?: PaginationMeta
  className?: string
  basePath?: string
}

function pageHref(
  basePath: string,
  searchParams: URLSearchParams,
  page: number,
): string {
  const params = new URLSearchParams(searchParams.toString())
  if (page <= 1) {
    params.delete('page')
  } else {
    params.set('page', page.toString())
  }
  const qs = params.toString()
  return qs ? `${basePath}?${qs}` : basePath
}

const pageLinkClass = (active: boolean, disabled?: boolean) =>
  clsx(
    'inline-flex items-center justify-center rounded-lg font-medium transition-all duration-200',
    'focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500',
    'px-3 py-1.5 text-sm min-w-[2.5rem]',
    disabled && 'pointer-events-none opacity-50',
    active
      ? 'bg-primary-500 text-white'
      : 'border-2 border-primary-500 text-primary-500 hover:bg-primary-50',
  )

export function Pagination({ meta, className, basePath = '/jobs' }: PaginationProps) {
  const searchParams = useSearchParams()

  if (!meta || meta.totalPages <= 1) {
    return null
  }

  const currentPage = meta.page || 1
  const totalPages = meta.totalPages || 1

  const getPageNumbers = () => {
    const pages: (number | string)[] = []
    const maxVisible = 5

    if (totalPages <= maxVisible) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i)
      }
    } else {
      pages.push(1)

      let start = Math.max(2, currentPage - 1)
      let end = Math.min(totalPages - 1, currentPage + 1)

      if (currentPage <= 3) {
        end = 4
      }

      if (currentPage >= totalPages - 2) {
        start = totalPages - 3
      }

      if (start > 2) {
        pages.push('...')
      }

      for (let i = start; i <= end; i++) {
        pages.push(i)
      }

      if (end < totalPages - 1) {
        pages.push('...')
      }

      pages.push(totalPages)
    }

    return pages
  }

  const startItem = (currentPage - 1) * meta.pageSize + 1
  const endItem = Math.min(currentPage * meta.pageSize, meta.totalCount)

  return (
    <div className={`flex flex-col sm:flex-row items-center justify-between gap-4 ${className}`}>
      <div className="text-sm text-secondary-600">
        Showing <span className="font-medium">{startItem}</span> to{' '}
        <span className="font-medium">{endItem}</span> of{' '}
        <span className="font-medium">{meta.totalCount}</span> results
      </div>

      <nav aria-label="Pagination" className="flex items-center gap-2">
        {currentPage === 1 ? (
          <span className={pageLinkClass(false, true)} aria-disabled="true">
            <ChevronLeft className="h-4 w-4 mr-1" />
            Previous
          </span>
        ) : (
          <Link
            href={pageHref(basePath, searchParams, currentPage - 1)}
            className={pageLinkClass(false)}
            rel="prev"
          >
            <ChevronLeft className="h-4 w-4 mr-1" />
            Previous
          </Link>
        )}

        <div className="flex items-center gap-1">
          {getPageNumbers().map((page, index) => {
            if (page === '...') {
              return (
                <span key={`ellipsis-${index}`} className="px-2 text-secondary-500">
                  ...
                </span>
              )
            }

            const pageNum = page as number
            const isActive = pageNum === currentPage

            return (
              <Link
                key={pageNum}
                href={pageHref(basePath, searchParams, pageNum)}
                className={pageLinkClass(isActive)}
                aria-current={isActive ? 'page' : undefined}
              >
                {pageNum}
              </Link>
            )
          })}
        </div>

        {currentPage === totalPages ? (
          <span className={pageLinkClass(false, true)} aria-disabled="true">
            Next
            <ChevronRight className="h-4 w-4 ml-1" />
          </span>
        ) : (
          <Link
            href={pageHref(basePath, searchParams, currentPage + 1)}
            className={pageLinkClass(false)}
            rel="next"
          >
            Next
            <ChevronRight className="h-4 w-4 ml-1" />
          </Link>
        )}
      </nav>
    </div>
  )
}
