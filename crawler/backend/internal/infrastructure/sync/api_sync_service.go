package sync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/startup-job-board/crawler/internal/domain/entity"
	"github.com/startup-job-board/crawler/internal/domain/repository"
)

// SyncService handles synchronization with the backend API
type SyncService struct {
	backendURL    string
	apiToken      string
	client        *http.Client
	jobRepo       repository.JobRepository
	siteRepo      repository.SiteRepository
	batchSize     int
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	SuccessCount int      `json:"success_count"`
	FailureCount int      `json:"failure_count"`
	Errors       []string `json:"errors"`
}

// NewSyncService creates a new sync service
func NewSyncService(backendURL, apiToken string, jobRepo repository.JobRepository, siteRepo repository.SiteRepository) *SyncService {
	return &SyncService{
		backendURL: backendURL,
		apiToken:   apiToken,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		jobRepo:   jobRepo,
		siteRepo: siteRepo,
		batchSize: 50, // Default batch size
	}
}

// SyncJobs syncs unsynced jobs to the backend API
func (s *SyncService) SyncJobs(ctx context.Context) (*SyncResult, error) {
	result := &SyncResult{
		Errors: []string{},
	}

	// Get unsynced jobs
	jobs, err := s.jobRepo.FindUnsynced(ctx, 0) // 0 = no limit
	if err != nil {
		return nil, fmt.Errorf("failed to fetch unsynced jobs: %w", err)
	}

	if len(jobs) == 0 {
		return result, nil
	}

	// Process in batches
	for i := 0; i < len(jobs); i += s.batchSize {
		end := i + s.batchSize
		if end > len(jobs) {
			end = len(jobs)
		}

		batch := jobs[i:end]
		if err := s.syncBatch(ctx, batch, result); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("batch sync error: %v", err))
			result.FailureCount += len(batch)
		}
	}

	return result, nil
}

// syncBatch syncs a batch of jobs
func (s *SyncService) syncBatch(ctx context.Context, jobs []*entity.CrawledJob, result *SyncResult) error {
	for _, job := range jobs {
		if err := s.SyncJob(ctx, job); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("job %s: %v", job.ID, err))
			result.FailureCount++
			continue
		}
		result.SuccessCount++
	}
	return nil
}

// SyncJob syncs a single job to the backend API
func (s *SyncService) SyncJob(ctx context.Context, job *entity.CrawledJob) error {
	// Get site to get backend startup ID
	site, err := s.siteRepo.FindByID(ctx, job.SiteID)
	if err != nil {
		return fmt.Errorf("failed to get site: %w", err)
	}

	// Map to backend DTO
	backendJob := s.mapToBackendJob(job, site.BackendStartupID)

	// Create request
	reqBody, err := json.Marshal(backendJob)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.backendURL+"/api/v1/token/jobs", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiToken)

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("backend API error: %d - %s", resp.StatusCode, string(body))
	}

	// Mark as synced
	if err := s.jobRepo.MarkAsSynced(ctx, job.ID); err != nil {
		return fmt.Errorf("failed to mark job as synced: %w", err)
	}

	return nil
}

// mapToBackendJob maps a crawled job to backend job DTO
func (s *SyncService) mapToBackendJob(job *entity.CrawledJob, startupID string) map[string]interface{} {
	backendJob := map[string]interface{}{
		"startup_id":      startupID,
		"title":           job.Title,
		"description":     job.Description,
		"requirements":    job.Requirements,
		"job_type":        s.normalizeJobType(job.JobType),
		"location_type":   s.normalizeLocationType(job.LocationType),
		"city":            job.City,
		"country":         job.Country,
		"currency":        s.normalizeCurrency(job.Currency),
		"application_url": job.ApplicationURL,
	}

	if job.SalaryMin != nil {
		backendJob["salary_min"] = *job.SalaryMin
	}
	if job.SalaryMax != nil {
		backendJob["salary_max"] = *job.SalaryMax
	}
	if job.ApplicationEmail != nil {
		backendJob["application_email"] = *job.ApplicationEmail
	}
	if job.ExpiresAt != nil {
		backendJob["expires_at"] = job.ExpiresAt.Format(time.RFC3339)
	}

	return backendJob
}

// normalizeJobType normalizes job type to backend format
func (s *SyncService) normalizeJobType(jobType string) string {
	normalized := map[string]string{
		"full-time":  "full_time",
		"fulltime":   "full_time",
		"part-time":  "part_time",
		"parttime":   "part_time",
		"contract":   "contract",
		"internship": "internship",
		"freelance":  "contract",
	}

	if val, ok := normalized[jobType]; ok {
		return val
	}
	return "full_time" // Default
}

// normalizeLocationType normalizes location type to backend format
func (s *SyncService) normalizeLocationType(locationType string) string {
	normalized := map[string]string{
		"remote":  "remote",
		"hybrid":  "hybrid",
		"onsite":  "onsite",
		"on-site": "onsite",
		"office":  "onsite",
	}

	if val, ok := normalized[locationType]; ok {
		return val
	}
	return "remote" // Default
}

// normalizeCurrency normalizes currency code
func (s *SyncService) normalizeCurrency(currency string) string {
	if len(currency) == 3 {
		return currency
	}
	return "USD" // Default
}

