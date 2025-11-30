package entity

import "errors"

var (
	ErrInvalidSiteName    = errors.New("site name is required")
	ErrInvalidBaseURL     = errors.New("base URL is required")
	ErrInvalidSchedule    = errors.New("schedule is required")
	ErrSiteNotFound       = errors.New("site not found")
	ErrJobNotFound        = errors.New("job not found")
	ErrInvalidExtraction  = errors.New("invalid extraction rules")
	ErrInvalidPagination  = errors.New("invalid pagination config")
)

