package site

import (
	"context"
	"github.com/startup-job-board/crawler/internal/domain/repository"
)

// DeleteSiteUseCase handles deleting a crawl site
type DeleteSiteUseCase struct {
	siteRepo repository.SiteRepository
}

// NewDeleteSiteUseCase creates a new DeleteSiteUseCase
func NewDeleteSiteUseCase(siteRepo repository.SiteRepository) *DeleteSiteUseCase {
	return &DeleteSiteUseCase{
		siteRepo: siteRepo,
	}
}

// Execute deletes a crawl site
func (uc *DeleteSiteUseCase) Execute(ctx context.Context, id string) error {
	return uc.siteRepo.Delete(ctx, id)
}

