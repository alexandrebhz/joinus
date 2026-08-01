package utils

const (
	DefaultPage       = 1
	DefaultPageSize   = 20
	MaxPage           = 500
	MaxPageSizePublic = 50
	MaxPageSizeAuth   = 100
)

// ClampPagination normalizes page/page_size. maxPageSize should be
// MaxPageSizePublic or MaxPageSizeAuth depending on caller trust.
func ClampPagination(page, pageSize, maxPageSize int) (int, int) {
	if page < 1 {
		page = DefaultPage
	}
	if page > MaxPage {
		page = MaxPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if maxPageSize < 1 {
		maxPageSize = MaxPageSizePublic
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}
