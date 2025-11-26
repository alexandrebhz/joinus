package gorm_model

import (
	"time"
)

type File struct {
	ID         string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FileName   string    `gorm:"type:varchar(255);not null"`
	FileSize   int64     `gorm:"type:bigint;not null"`
	MimeType   string    `gorm:"type:varchar(100);not null"`
	StorageKey string    `gorm:"type:varchar(500);not null"`
	URL        string    `gorm:"type:varchar(500);not null"`
	UploadedBy string    `gorm:"type:uuid;not null"`
	CreatedAt  time.Time
}

func (File) TableName() string {
	return "files"
}



