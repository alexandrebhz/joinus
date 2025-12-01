package gorm_model

import (
	"time"

	"gorm.io/gorm"
)

// CrawledJobModel represents the GORM model for crawled jobs
type CrawledJobModel struct {
	ID                string         `gorm:"primaryKey;type:uuid"`
	SiteID            string         `gorm:"type:uuid;not null;index"`
	ExternalID        string         `gorm:"type:varchar(255);index"`
	DetailURL         string         `gorm:"type:text;not null;index"`
	Title             string         `gorm:"type:varchar(500);not null"`
	Description       string         `gorm:"type:text"`
	Requirements      string         `gorm:"type:text"`
	Company           string         `gorm:"type:varchar(255)"`
	Location          string         `gorm:"type:varchar(255)"`
	City              string         `gorm:"type:varchar(100)"`
	Country           string         `gorm:"type:varchar(100)"`
	JobType           string         `gorm:"type:varchar(50)"`
	LocationType      string         `gorm:"type:varchar(50)"`
	SalaryMin         *int           `gorm:"type:integer"`
	SalaryMax         *int           `gorm:"type:integer"`
	Currency          string         `gorm:"type:varchar(3)"`
	ApplicationURL    *string        `gorm:"type:text"`
	ApplicationEmail  *string        `gorm:"type:varchar(255)"`
	ExpiresAt         *time.Time     `gorm:"type:timestamp"`
	RawHTML           string         `gorm:"type:text"`
	DeduplicationHash string         `gorm:"type:varchar(255);index"`
	Synced            bool           `gorm:"default:false;index"`
	SyncedAt          *time.Time     `gorm:"type:timestamp"`
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name
func (CrawledJobModel) TableName() string {
	return "crawled_jobs"
}

