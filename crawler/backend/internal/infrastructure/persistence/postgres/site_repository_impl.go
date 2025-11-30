package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/startup-job-board/crawler/internal/domain/entity"
	"github.com/startup-job-board/crawler/internal/domain/repository"
	"github.com/startup-job-board/crawler/internal/infrastructure/persistence/gorm_model"
	"gorm.io/gorm"
)

// SiteRepositoryImpl implements SiteRepository using PostgreSQL
type SiteRepositoryImpl struct {
	db *gorm.DB
}

// NewSiteRepository creates a new SiteRepositoryImpl
func NewSiteRepository(db *gorm.DB) repository.SiteRepository {
	return &SiteRepositoryImpl{db: db}
}

// Create creates a new crawl site
func (r *SiteRepositoryImpl) Create(ctx context.Context, site *entity.CrawlSite) error {
	model := r.toModel(site)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update updates an existing crawl site
func (r *SiteRepositoryImpl) Update(ctx context.Context, site *entity.CrawlSite) error {
	model := r.toModel(site)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete deletes a crawl site
func (r *SiteRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&gorm_model.CrawlSiteModel{}, "id = ?", id).Error
}

// FindByID finds a crawl site by ID
func (r *SiteRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.CrawlSite, error) {
	var model gorm_model.CrawlSiteModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return r.toEntity(&model), nil
}

// FindAll finds all crawl sites
func (r *SiteRepositoryImpl) FindAll(ctx context.Context) ([]*entity.CrawlSite, error) {
	var models []gorm_model.CrawlSiteModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	sites := make([]*entity.CrawlSite, len(models))
	for i, model := range models {
		sites[i] = r.toEntity(&model)
	}
	return sites, nil
}

// FindActive finds all active crawl sites
func (r *SiteRepositoryImpl) FindActive(ctx context.Context) ([]*entity.CrawlSite, error) {
	var models []gorm_model.CrawlSiteModel
	if err := r.db.WithContext(ctx).Where("active = ?", true).Find(&models).Error; err != nil {
		return nil, err
	}

	sites := make([]*entity.CrawlSite, len(models))
	for i, model := range models {
		sites[i] = r.toEntity(&model)
	}
	return sites, nil
}

// UpdateLastCrawledAt updates the last crawled timestamp
func (r *SiteRepositoryImpl) UpdateLastCrawledAt(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&gorm_model.CrawlSiteModel{}).
		Where("id = ?", id).
		Update("last_crawled_at", now).Error
}

// toModel converts entity to GORM model
func (r *SiteRepositoryImpl) toModel(site *entity.CrawlSite) *gorm_model.CrawlSiteModel {
	paginationJSON, _ := site.PaginationConfig.ToJSON()
	extractionJSON, _ := site.ExtractionRules.ToJSON()

	var paginationMap map[string]interface{}
	var extractionMap map[string]interface{}
	json.Unmarshal(paginationJSON, &paginationMap)
	json.Unmarshal(extractionJSON, &extractionMap)

	return &gorm_model.CrawlSiteModel{
		ID:               site.ID,
		Name:             site.Name,
		BaseURL:          site.BaseURL,
		BackendStartupID: site.BackendStartupID,
		Active:           site.Active,
		Schedule:         site.Schedule,
		LastCrawledAt:    site.LastCrawledAt,
		NextCrawlAt:      site.NextCrawlAt,
		CrawlInterval:    site.CrawlInterval,
		PaginationConfig: gorm_model.JSONB(paginationMap),
		ExtractionRules:  gorm_model.JSONB(extractionMap),
		DeduplicationKey: site.DeduplicationKey,
		RequestDelay:     site.RequestDelay,
		UserAgent:        site.UserAgent,
		CreatedAt:        site.CreatedAt,
		UpdatedAt:        site.UpdatedAt,
	}
}

// toEntity converts GORM model to entity
func (r *SiteRepositoryImpl) toEntity(model *gorm_model.CrawlSiteModel) *entity.CrawlSite {
	paginationJSON, _ := json.Marshal(model.PaginationConfig)
	extractionJSON, _ := json.Marshal(model.ExtractionRules)

	var paginationConfig entity.PaginationConfig
	var extractionRules entity.ExtractionRules
	paginationConfig.FromJSON(paginationJSON)
	extractionRules.FromJSON(extractionJSON)

	return &entity.CrawlSite{
		ID:               model.ID,
		Name:             model.Name,
		BaseURL:          model.BaseURL,
		BackendStartupID: model.BackendStartupID,
		Active:           model.Active,
		Schedule:         model.Schedule,
		LastCrawledAt:    model.LastCrawledAt,
		NextCrawlAt:      model.NextCrawlAt,
		CrawlInterval:    model.CrawlInterval,
		PaginationConfig: paginationConfig,
		ExtractionRules:  extractionRules,
		DeduplicationKey: model.DeduplicationKey,
		RequestDelay:     model.RequestDelay,
		UserAgent:        model.UserAgent,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}
}

