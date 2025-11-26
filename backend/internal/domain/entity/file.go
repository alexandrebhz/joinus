package entity

import "time"

type File struct {
	ID         string
	FileName   string
	FileSize   int64
	MimeType   string
	StorageKey string
	URL        string
	UploadedBy string
	CreatedAt  time.Time
}



