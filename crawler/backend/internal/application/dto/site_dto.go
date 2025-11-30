package dto

import "time"

// CreateSiteInput represents input for creating a crawl site
type CreateSiteInput struct {
	Name             string                 `json:"name" validate:"required,min=3,max=100"`
	BaseURL          string                 `json:"base_url" validate:"required,url"`
	BackendStartupID string                 `json:"backend_startup_id" validate:"required"`
	Schedule         string                 `json:"schedule" validate:"required"`
	CrawlInterval    string                 `json:"crawl_interval" validate:"required,oneof=daily weekly custom"`
	PaginationConfig map[string]interface{} `json:"pagination_config" validate:"required"`
	ExtractionRules  map[string]interface{} `json:"extraction_rules" validate:"required"`
	DeduplicationKey string                 `json:"deduplication_key" validate:"omitempty,oneof=url composite external_id"`
	RequestDelay     int                    `json:"request_delay" validate:"omitempty,min=0,max=60"`
	UserAgent        string                 `json:"user_agent"`
}

// UpdateSiteInput represents input for updating a crawl site
type UpdateSiteInput struct {
	Name             *string                `json:"name" validate:"omitempty,min=3,max=100"`
	BaseURL          *string                `json:"base_url" validate:"omitempty,url"`
	BackendStartupID *string                `json:"backend_startup_id"`
	Active           *bool                  `json:"active"`
	Schedule         *string                `json:"schedule"`
	CrawlInterval    *string                `json:"crawl_interval" validate:"omitempty,oneof=daily weekly custom"`
	PaginationConfig map[string]interface{} `json:"pagination_config"`
	ExtractionRules  map[string]interface{} `json:"extraction_rules"`
	DeduplicationKey *string                `json:"deduplication_key" validate:"omitempty,oneof=url composite external_id"`
	RequestDelay     *int                   `json:"request_delay" validate:"omitempty,min=0,max=60"`
	UserAgent        *string                `json:"user_agent"`
}

// SiteOutput represents a crawl site in API responses
type SiteOutput struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	BaseURL          string                 `json:"base_url"`
	BackendStartupID string                 `json:"backend_startup_id"`
	Active           bool                   `json:"active"`
	Schedule         string                 `json:"schedule"`
	LastCrawledAt    *time.Time             `json:"last_crawled_at"`
	NextCrawlAt      *time.Time             `json:"next_crawl_at"`
	CrawlInterval    string                 `json:"crawl_interval"`
	PaginationConfig map[string]interface{} `json:"pagination_config"`
	ExtractionRules  map[string]interface{} `json:"extraction_rules"`
	DeduplicationKey string                 `json:"deduplication_key"`
	RequestDelay     int                    `json:"request_delay"`
	UserAgent        string                 `json:"user_agent"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

