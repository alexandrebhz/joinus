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
  pageSize: number
  totalPages: number
  totalCount: number
}

