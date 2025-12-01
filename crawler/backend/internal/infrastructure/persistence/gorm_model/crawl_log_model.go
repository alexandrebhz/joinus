package gorm_model

import (
	"time"
)

// CrawlLogModel represents the GORM model for crawl logs
type CrawlLogModel struct {
	ID           string         `gorm:"primaryKey;type:uuid"`
	SiteID       string         `gorm:"type:uuid;not null;index"`
	Status       string         `gorm:"type:varchar(50);not null;default:'running'"`
	StartedAt    time.Time      `gorm:"type:timestamp;not null"`
	CompletedAt  *time.Time     `gorm:"type:timestamp"`
	DurationMs   int            `gorm:"type:integer"`
	JobsFound    int            `gorm:"type:integer;default:0"`
	JobsSaved    int            `gorm:"type:integer;default:0"`
	JobsSkipped  int            `gorm:"type:integer;default:0"`
	PagesCrawled int            `gorm:"type:integer;default:0"`
	Errors       JSONB          `gorm:"type:jsonb;default:'[]'::jsonb"`
	Logs         JSONB          `gorm:"type:jsonb;default:'[]'::jsonb"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
}

// TableName specifies the table name
func (CrawlLogModel) TableName() string {
	return "crawl_logs"
}

