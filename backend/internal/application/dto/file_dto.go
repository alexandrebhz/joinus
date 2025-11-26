package dto

type FileOutput struct {
	ID        string `json:"id"`
	FileName  string `json:"file_name"`
	FileSize  int64  `json:"file_size"`
	MimeType  string `json:"mime_type"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}



