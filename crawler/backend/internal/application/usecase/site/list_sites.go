package site

import (
	"context"
	"encoding/json"

	"github.com/startup-job-board/crawler/internal/application/dto"
	"github.com/startup-job-board/crawler/internal/domain/entity"
	"github.com/startup-job-board/crawler/internal/domain/repository"
)

// ListSitesUseCase handles listing crawl sites
type ListSitesUseCase struct {
	siteRepo repository.SiteRepository
}

// NewListSitesUseCase creates a new ListSitesUseCase
func NewListSitesUseCase(siteRepo repository.SiteRepository) *ListSitesUseCase {
	return &ListSitesUseCase{
		siteRepo: siteRepo,
	}
}

// Execute lists all crawl sites
func (uc *ListSitesUseCase) Execute(ctx context.Context) ([]*dto.SiteOutput, error) {
	sites, err := uc.siteRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	output := make([]*dto.SiteOutput, len(sites))
	for i, site := range sites {
		output[i] = uc.toDTO(site)
	}

	return output, nil
}

func (uc *ListSitesUseCase) toDTO(site *entity.CrawlSite) *dto.SiteOutput {
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

