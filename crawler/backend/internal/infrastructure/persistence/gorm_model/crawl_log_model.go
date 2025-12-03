package gorm_model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// CrawlLogModel represents the GORM model for crawl logs
type CrawlLogModel struct {
	ID           string         `gorm:"primaryKey;type:uuid"`
	SiteID       string         `gorm:"type:uuid;not null;index"`
	Status       string         `gorm:"type:varchar(50);not null;default:'running';index"`
	StartedAt    time.Time      `gorm:"type:timestamp;not null;index:idx_crawl_logs_started_at"`
	CompletedAt  *time.Time     `gorm:"type:timestamp"`
	DurationMs   int            `gorm:"type:integer"`
	JobsFound    int            `gorm:"type:integer;default:0"`
	JobsSaved    int            `gorm:"type:integer;default:0"`
	JobsSkipped  int            `gorm:"type:integer;default:0"`
	PagesCrawled int            `gorm:"type:integer;default:0"`
	Errors       JSONBArray     `gorm:"type:jsonb;default:'[]'::jsonb"`
	Logs         JSONBArray     `gorm:"type:jsonb;default:'[]'::jsonb"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
}

// TableName specifies the table name
func (CrawlLogModel) TableName() string {
	return "crawl_logs"
}

// JSONBArray is a custom type for JSONB array columns
type JSONBArray []interface{}

// Value implements the driver.Valuer interface
func (j JSONBArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONBArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}

