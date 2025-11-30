package site

import (
	"context"
	"encoding/json"
	"time"

	"github.com/startup-job-board/crawler/internal/application/dto"
	"github.com/startup-job-board/crawler/internal/domain/entity"
	"github.com/startup-job-board/crawler/internal/domain/repository"
)

// UpdateSiteUseCase handles updating a crawl site
type UpdateSiteUseCase struct {
	siteRepo repository.SiteRepository
}

// NewUpdateSiteUseCase creates a new UpdateSiteUseCase
func NewUpdateSiteUseCase(siteRepo repository.SiteRepository) *UpdateSiteUseCase {
	return &UpdateSiteUseCase{
		siteRepo: siteRepo,
	}
}

// Execute updates an existing crawl site
func (uc *UpdateSiteUseCase) Execute(ctx context.Context, id string, input dto.UpdateSiteInput) (*dto.SiteOutput, error) {
	site, err := uc.siteRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if input.Name != nil {
		site.Name = *input.Name
	}
	if input.BaseURL != nil {
		site.BaseURL = *input.BaseURL
	}
	if input.BackendStartupID != nil {
		site.BackendStartupID = *input.BackendStartupID
	}
	if input.Active != nil {
		site.Active = *input.Active
	}
	if input.Schedule != nil {
		site.Schedule = *input.Schedule
	}
	if input.CrawlInterval != nil {
		site.CrawlInterval = *input.CrawlInterval
	}
	if input.PaginationConfig != nil {
		paginationJSON, _ := json.Marshal(input.PaginationConfig)
		site.PaginationConfig.FromJSON(paginationJSON)
	}
	if input.ExtractionRules != nil {
		extractionJSON, _ := json.Marshal(input.ExtractionRules)
		site.ExtractionRules.FromJSON(extractionJSON)
	}
	if input.DeduplicationKey != nil {
		site.DeduplicationKey = *input.DeduplicationKey
	}
	if input.RequestDelay != nil {
		site.RequestDelay = *input.RequestDelay
	}
	if input.UserAgent != nil {
		site.UserAgent = *input.UserAgent
	}

	site.UpdatedAt = time.Now()

	if err := site.Validate(); err != nil {
		return nil, err
	}

	if err := uc.siteRepo.Update(ctx, site); err != nil {
		return nil, err
	}

	return uc.toDTO(site), nil
}

func (uc *UpdateSiteUseCase) toDTO(site *entity.CrawlSite) *dto.SiteOutput {
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

