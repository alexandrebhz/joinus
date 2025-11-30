package dto

import "time"

// CrawledJobOutput represents a crawled job in API responses
type CrawledJobOutput struct {
	ID                string     `json:"id"`
	SiteID            string     `json:"site_id"`
	SiteName          string     `json:"site_name,omitempty"`
	ExternalID        string     `json:"external_id"`
	DetailURL         string     `json:"detail_url"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	Requirements      string     `json:"requirements"`
	Company           string     `json:"company"`
	Location          string     `json:"location"`
	City              string     `json:"city"`
	Country           string     `json:"country"`
	JobType           string     `json:"job_type"`
	LocationType      string     `json:"location_type"`
	SalaryMin         *int       `json:"salary_min"`
	SalaryMax         *int       `json:"salary_max"`
	Currency          string     `json:"currency"`
	ApplicationURL    *string    `json:"application_url"`
	ApplicationEmail  *string    `json:"application_email"`
	ExpiresAt         *time.Time `json:"expires_at"`
	Synced            bool       `json:"synced"`
	SyncedAt          *time.Time `json:"synced_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// SyncStatsOutput represents synchronization statistics
type SyncStatsOutput struct {
	TotalJobs      int64 `json:"total_jobs"`
	UnsyncedJobs   int64 `json:"unsynced_jobs"`
	SyncedJobs     int64 `json:"synced_jobs"`
	LastSyncTime   *time.Time `json:"last_sync_time"`
}

