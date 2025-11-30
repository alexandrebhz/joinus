package gorm_model

import (
	"time"

	"gorm.io/gorm"
)

// CrawlSiteModel represents the GORM model for crawl sites
type CrawlSiteModel struct {
	ID               string         `gorm:"primaryKey;type:uuid"`
	Name             string         `gorm:"type:varchar(255);not null"`
	BaseURL          string         `gorm:"type:text;not null"`
	BackendStartupID string         `gorm:"type:uuid;not null"`
	Active           bool           `gorm:"default:true"`
	Schedule         string         `gorm:"type:varchar(100);not null"`
	LastCrawledAt    *time.Time     `gorm:"type:timestamp"`
	NextCrawlAt      *time.Time     `gorm:"type:timestamp"`
	CrawlInterval    string         `gorm:"type:varchar(50);not null"`
	PaginationConfig JSONB          `gorm:"type:jsonb;not null"`
	ExtractionRules  JSONB          `gorm:"type:jsonb;not null"`
	DeduplicationKey string         `gorm:"type:varchar(50);default:'url'"`
	RequestDelay     int            `gorm:"default:2"`
	UserAgent        string         `gorm:"type:text"`
	CreatedAt        time.Time      `gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name
func (CrawlSiteModel) TableName() string {
	return "crawl_sites"
}

