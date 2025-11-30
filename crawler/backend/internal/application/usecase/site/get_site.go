package site

import (
	"context"
	"encoding/json"

	"github.com/startup-job-board/crawler/internal/application/dto"
	"github.com/startup-job-board/crawler/internal/domain/entity"
	"github.com/startup-job-board/crawler/internal/domain/repository"
)

// GetSiteUseCase handles getting a single crawl site
type GetSiteUseCase struct {
	siteRepo repository.SiteRepository
}

// NewGetSiteUseCase creates a new GetSiteUseCase
func NewGetSiteUseCase(siteRepo repository.SiteRepository) *GetSiteUseCase {
	return &GetSiteUseCase{
		siteRepo: siteRepo,
	}
}

// Execute gets a crawl site by ID
func (uc *GetSiteUseCase) Execute(ctx context.Context, id string) (*dto.SiteOutput, error) {
	site, err := uc.siteRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return uc.toDTO(site), nil
}

func (uc *GetSiteUseCase) toDTO(site *entity.CrawlSite) *dto.SiteOutput {
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

