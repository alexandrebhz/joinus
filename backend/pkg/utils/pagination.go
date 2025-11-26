package utils

type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
	TotalCount int64 `json:"total_count"`
}

func CalculateTotalPages(totalCount int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}
	return totalPages
}



