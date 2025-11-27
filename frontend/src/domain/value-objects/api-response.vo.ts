export interface ApiResponse<T> {
  success: boolean
  data: T
  meta?: PaginationMeta
}

export interface ApiError {
  success: false
  error: string
  code?: string
}

export interface PaginationMeta {
  page: number
  page_size: number
  total_pages: number
  total_count: number
}

