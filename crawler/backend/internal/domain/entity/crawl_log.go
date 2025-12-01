package entity

import (
	"time"
)

// CrawlLog represents a log entry for a crawl execution
type CrawlLog struct {
	ID          string    `json:"id"`
	SiteID      string    `json:"site_id"`
	Status      string    `json:"status"` // "running", "completed", "failed"
	StartedAt   time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	DurationMs  int        `json:"duration_ms"`
	JobsFound   int        `json:"jobs_found"`
	JobsSaved   int        `json:"jobs_saved"`
	JobsSkipped int        `json:"jobs_skipped"`
	PagesCrawled int       `json:"pages_crawled"`
	Errors      []string   `json:"errors"`
	Logs        []LogEntry `json:"logs"`
	CreatedAt   time.Time  `json:"created_at"`
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"` // "info", "warning", "error"
	Message   string    `json:"message"`
}

// AddLog adds a log entry
func (cl *CrawlLog) AddLog(level string, message string) {
	cl.Logs = append(cl.Logs, LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	})
}

// AddError adds an error message
func (cl *CrawlLog) AddError(err error) {
	if err != nil {
		cl.Errors = append(cl.Errors, err.Error())
		cl.AddLog("error", err.Error())
	}
}

// Complete marks the crawl as completed
func (cl *CrawlLog) Complete() {
	now := time.Now()
	cl.CompletedAt = &now
	cl.Status = "completed"
	cl.DurationMs = int(now.Sub(cl.StartedAt).Milliseconds())
}

// Fail marks the crawl as failed
func (cl *CrawlLog) Fail(err error) {
	now := time.Now()
	cl.CompletedAt = &now
	cl.Status = "failed"
	cl.DurationMs = int(now.Sub(cl.StartedAt).Milliseconds())
	if err != nil {
		cl.AddError(err)
	}
}

