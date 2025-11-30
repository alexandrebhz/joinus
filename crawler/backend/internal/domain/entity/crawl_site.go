package entity

import (
	"encoding/json"
	"time"
)

// CrawlSite represents a website to crawl for job listings
type CrawlSite struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	BaseURL          string           `json:"base_url"`
	BackendStartupID string           `json:"backend_startup_id"` // Maps to backend startup
	Active           bool             `json:"active"`
	Schedule         string           `json:"schedule"` // Cron expression
	LastCrawledAt    *time.Time       `json:"last_crawled_at"`
	NextCrawlAt      *time.Time       `json:"next_crawl_at"`
	CrawlInterval    string           `json:"crawl_interval"` // "daily", "weekly", "custom"
	PaginationConfig PaginationConfig `json:"pagination_config"`
	ExtractionRules  ExtractionRules  `json:"extraction_rules"`
	DeduplicationKey string           `json:"deduplication_key"` // "url", "composite", "external_id"
	RequestDelay     int              `json:"request_delay"`      // Delay in seconds between requests
	UserAgent        string           `json:"user_agent"`        // Custom user agent
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

// PaginationConfig defines how to paginate through job listings
type PaginationConfig struct {
	Type             string              `json:"type"` // "query_param", "url_pattern", "link_follow", "api_pagination"
	ParamName        string              `json:"param_name,omitempty"`
	StartPage        int                 `json:"start_page,omitempty"`
	Increment        int                 `json:"increment,omitempty"`
	MaxPages         int                 `json:"max_pages,omitempty"`
	NextPageSelector string             `json:"next_page_selector,omitempty"`
	URLPattern       string             `json:"url_pattern,omitempty"`
	APIConfig        *APIPaginationConfig `json:"api_config,omitempty"`
}

// APIPaginationConfig for API-based pagination
type APIPaginationConfig struct {
	Endpoint  string `json:"endpoint"`
	PageParam string `json:"page_param"`
	PageSize  int    `json:"page_size"`
	MaxPages  int    `json:"max_pages"`
}

// ExtractionRules defines how to extract job data from HTML
type ExtractionRules struct {
	JobListSelector string                `json:"job_list_selector"`
	JobDetailURL    JobURLRule            `json:"job_detail_url"`
	Fields          map[string]FieldRule  `json:"fields"`
}

// JobURLRule defines how to extract job detail page URL
type JobURLRule struct {
	Type      string `json:"type"`       // "relative", "absolute", "attribute"
	Selector  string `json:"selector"`
	Attribute string `json:"attribute,omitempty"`
	BaseURL   string `json:"base_url,omitempty"`
}

// FieldRule defines how to extract a specific field from HTML
type FieldRule struct {
	Selector        string   `json:"selector"`
	Type            string   `json:"type"` // "text", "html", "attribute", "regex"
	Attribute       string   `json:"attribute,omitempty"`
	RegexPattern    string   `json:"regex_pattern,omitempty"`
	Required        bool     `json:"required"`
	DefaultValue    string   `json:"default_value,omitempty"`
	Transformations []string `json:"transformations,omitempty"`
}

// Validate validates the crawl site configuration
func (cs *CrawlSite) Validate() error {
	if cs.Name == "" {
		return ErrInvalidSiteName
	}
	if cs.BaseURL == "" {
		return ErrInvalidBaseURL
	}
	if cs.Schedule == "" {
		return ErrInvalidSchedule
	}
	if cs.DeduplicationKey == "" {
		cs.DeduplicationKey = "url" // Default to URL-based
	}
	return nil
}

// ToJSON converts PaginationConfig to JSON
func (pc *PaginationConfig) ToJSON() ([]byte, error) {
	return json.Marshal(pc)
}

// FromJSON parses PaginationConfig from JSON
func (pc *PaginationConfig) FromJSON(data []byte) error {
	return json.Unmarshal(data, pc)
}

// ToJSON converts ExtractionRules to JSON
func (er *ExtractionRules) ToJSON() ([]byte, error) {
	return json.Marshal(er)
}

// FromJSON parses ExtractionRules from JSON
func (er *ExtractionRules) FromJSON(data []byte) error {
	return json.Unmarshal(data, er)
}

