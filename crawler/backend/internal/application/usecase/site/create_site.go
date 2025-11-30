package site

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/crawler/internal/application/dto"
	"github.com/startup-job-board/crawler/internal/domain/entity"
	"github.com/startup-job-board/crawler/internal/domain/repository"
)

// CreateSiteUseCase handles creating a new crawl site
type CreateSiteUseCase struct {
	siteRepo repository.SiteRepository
}

// NewCreateSiteUseCase creates a new CreateSiteUseCase
func NewCreateSiteUseCase(siteRepo repository.SiteRepository) *CreateSiteUseCase {
	return &CreateSiteUseCase{
		siteRepo: siteRepo,
	}
}

// Execute creates a new crawl site
func (uc *CreateSiteUseCase) Execute(ctx context.Context, input dto.CreateSiteInput) (*dto.SiteOutput, error) {
	// Parse pagination config
	paginationConfig := entity.PaginationConfig{}
	paginationJSON, _ := json.Marshal(input.PaginationConfig)
	if err := paginationConfig.FromJSON(paginationJSON); err != nil {
		return nil, err
	}

	// Parse extraction rules
	extractionRules := entity.ExtractionRules{}
	extractionJSON, _ := json.Marshal(input.ExtractionRules)
	if err := extractionRules.FromJSON(extractionJSON); err != nil {
		return nil, err
	}

	// Set defaults
	deduplicationKey := input.DeduplicationKey
	if deduplicationKey == "" {
		deduplicationKey = "url"
	}
	requestDelay := input.RequestDelay
	if requestDelay == 0 {
		requestDelay = 2 // Default 2 seconds
	}
	userAgent := input.UserAgent
	if userAgent == "" {
		userAgent = "Mozilla/5.0 (compatible; JobCrawler/1.0)"
	}

	now := time.Now()
	site := &entity.CrawlSite{
		ID:               uuid.New().String(),
		Name:             input.Name,
		BaseURL:          input.BaseURL,
		BackendStartupID: input.BackendStartupID,
		Active:           true,
		Schedule:         input.Schedule,
		CrawlInterval:    input.CrawlInterval,
		PaginationConfig: paginationConfig,
		ExtractionRules:  extractionRules,
		DeduplicationKey: deduplicationKey,
		RequestDelay:     requestDelay,
		UserAgent:        userAgent,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := site.Validate(); err != nil {
		return nil, err
	}

	if err := uc.siteRepo.Create(ctx, site); err != nil {
		return nil, err
	}

	return uc.toDTO(site), nil
}

func (uc *CreateSiteUseCase) toDTO(site *entity.CrawlSite) *dto.SiteOutput {
	paginationJSON, _ := site.PaginationConfig.ToJSON()
	extractionJSON, _ := site.ExtractionRules.ToJSON()

	var paginationMap map[string]interface{}
	var extractionMap map[string]interface{}
	json.Unmarshal(paginationJSON, &paginationMap)
	json.Unmarshal(extractionJSON, &extractionMap)

	return &dto.SiteOutput{
		ID:               site.ID,
		Name:             site.Name,
		BaseURL:          site.BaseURL,
		BackendStartupID: site.BackendStartupID,
		Active:           site.Active,
		Schedule:         site.Schedule,
		LastCrawledAt:    site.LastCrawledAt,
		NextCrawlAt:      site.NextCrawlAt,
		CrawlInterval:    site.CrawlInterval,
		PaginationConfig: paginationMap,
		ExtractionRules:  extractionMap,
		DeduplicationKey: site.DeduplicationKey,
		RequestDelay:     site.RequestDelay,
		UserAgent:        site.UserAgent,
		CreatedAt:        site.CreatedAt,
		UpdatedAt:        site.UpdatedAt,
	}
}

