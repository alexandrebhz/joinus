package entity

import "time"

// CrawledJob represents a job that has been crawled from an external site
type CrawledJob struct {
	ID                string     `json:"id"`
	SiteID            string     `json:"site_id"`
	ExternalID        string     `json:"external_id"`
	DetailURL         string     `json:"detail_url"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	Requirements      string     `json:"requirements"`
	Company           string     `json:"company"`
	Location          string     `json:"location"`
	City              string     `json:"city"`
	Country           string     `json:"country"`
	JobType           string     `json:"job_type"`        // "full_time", "part_time", "contract", "internship"
	LocationType      string     `json:"location_type"`    // "remote", "hybrid", "onsite"
	SalaryMin         *int       `json:"salary_min"`
	SalaryMax         *int       `json:"salary_max"`
	Currency          string     `json:"currency"`
	ApplicationURL    *string    `json:"application_url"`
	ApplicationEmail  *string    `json:"application_email"`
	ExpiresAt         *time.Time `json:"expires_at"`
	RawHTML           string     `json:"raw_html"`
	DeduplicationHash string     `json:"deduplication_hash"`
	Synced            bool       `json:"synced"`
	SyncedAt          *time.Time `json:"synced_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// GenerateDeduplicationHash generates a hash for deduplication based on the strategy
func (cj *CrawledJob) GenerateDeduplicationHash(strategy string) string {
	switch strategy {
	case "url":
		return cj.DetailURL
	case "composite":
		return generateCompositeHash(cj.Title, cj.Company, cj.Location)
	case "external_id":
		return cj.ExternalID
	default:
		return cj.DetailURL
	}
}

// generateCompositeHash creates a hash from multiple fields
func generateCompositeHash(fields ...string) string {
	// Simple hash implementation - in production, use crypto/sha256
	hash := ""
	for _, field := range fields {
		hash += field + "|"
	}
	return hash
}

