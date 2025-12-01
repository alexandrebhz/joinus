package postgres

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/startup-job-board/crawler/internal/domain/entity"
	"github.com/startup-job-board/crawler/internal/domain/repository"
	"github.com/startup-job-board/crawler/internal/infrastructure/persistence/gorm_model"
	"gorm.io/gorm"
)

// CrawlLogRepositoryImpl implements CrawlLogRepository using PostgreSQL
type CrawlLogRepositoryImpl struct {
	db *gorm.DB
}

// NewCrawlLogRepository creates a new CrawlLogRepositoryImpl
func NewCrawlLogRepository(db *gorm.DB) repository.CrawlLogRepository {
	return &CrawlLogRepositoryImpl{db: db}
}

// Create creates a new crawl log
func (r *CrawlLogRepositoryImpl) Create(ctx context.Context, log *entity.CrawlLog) error {
	model := r.toModel(log)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update updates an existing crawl log
func (r *CrawlLogRepositoryImpl) Update(ctx context.Context, log *entity.CrawlLog) error {
	model := r.toModel(log)
	return r.db.WithContext(ctx).Save(model).Error
}

// FindByID finds a crawl log by ID
func (r *CrawlLogRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.CrawlLog, error) {
	var model gorm_model.CrawlLogModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

// FindBySiteID finds crawl logs for a site
func (r *CrawlLogRepositoryImpl) FindBySiteID(ctx context.Context, siteID string, limit int) ([]*entity.CrawlLog, error) {
	var models []gorm_model.CrawlLogModel
	query := r.db.WithContext(ctx).Where("site_id = ?", siteID).Order("started_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	logs := make([]*entity.CrawlLog, len(models))
	for i, model := range models {
		logs[i] = r.toDomain(&model)
	}
	return logs, nil
}

// FindLatestBySiteID finds the latest crawl log for a site
func (r *CrawlLogRepositoryImpl) FindLatestBySiteID(ctx context.Context, siteID string) (*entity.CrawlLog, error) {
	var model gorm_model.CrawlLogModel
	if err := r.db.WithContext(ctx).
		Where("site_id = ?", siteID).
		Order("started_at DESC").
		First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

// toModel converts domain entity to GORM model
func (r *CrawlLogRepositoryImpl) toModel(log *entity.CrawlLog) *gorm_model.CrawlLogModel {
	if log.ID == "" {
		log.ID = uuid.New().String()
	}

	var errorsJSONB gorm_model.JSONB
	var logsJSONB gorm_model.JSONB
	
	errorsJSON, _ := json.Marshal(log.Errors)
	logsJSON, _ := json.Marshal(log.Logs)
	
	var errorsInterface interface{}
	var logsInterface interface{}
	json.Unmarshal(errorsJSON, &errorsInterface)
	json.Unmarshal(logsJSON, &logsInterface)
	
	// Convert to JSONB (map[string]interface{})
	// Since JSONB is map[string]interface{}, we need to wrap arrays in a map or use a different approach
	// For arrays, we'll store them as {"data": array}
	if errorsMap, ok := errorsInterface.([]interface{}); ok {
		errorsJSONB = gorm_model.JSONB{"data": errorsMap}
	} else if errorsMap, ok := errorsInterface.(map[string]interface{}); ok {
		errorsJSONB = gorm_model.JSONB(errorsMap)
	}
	
	if logsMap, ok := logsInterface.([]interface{}); ok {
		logsJSONB = gorm_model.JSONB{"data": logsMap}
	} else if logsMap, ok := logsInterface.(map[string]interface{}); ok {
		logsJSONB = gorm_model.JSONB(logsMap)
	}

	return &gorm_model.CrawlLogModel{
		ID:           log.ID,
		SiteID:       log.SiteID,
		Status:       log.Status,
		StartedAt:    log.StartedAt,
		CompletedAt:  log.CompletedAt,
		DurationMs:   log.DurationMs,
		JobsFound:    log.JobsFound,
		JobsSaved:    log.JobsSaved,
		JobsSkipped:  log.JobsSkipped,
		PagesCrawled: log.PagesCrawled,
		Errors:       errorsJSONB,
		Logs:         logsJSONB,
		CreatedAt:    log.CreatedAt,
	}
}

// toDomain converts GORM model to domain entity
func (r *CrawlLogRepositoryImpl) toDomain(model *gorm_model.CrawlLogModel) *entity.CrawlLog {
	log := &entity.CrawlLog{
		ID:           model.ID,
		SiteID:       model.SiteID,
		Status:       model.Status,
		StartedAt:    model.StartedAt,
		CompletedAt:  model.CompletedAt,
		DurationMs:   model.DurationMs,
		JobsFound:    model.JobsFound,
		JobsSaved:    model.JobsSaved,
		JobsSkipped:  model.JobsSkipped,
		PagesCrawled: model.PagesCrawled,
		CreatedAt:    model.CreatedAt,
	}

	// Parse errors JSON
	if model.Errors != nil {
		errorsJSON, _ := json.Marshal(model.Errors)
		var errorsData interface{}
		json.Unmarshal(errorsJSON, &errorsData)
		// If wrapped in {"data": array}, unwrap it
		if errorsMap, ok := errorsData.(map[string]interface{}); ok {
			if data, exists := errorsMap["data"]; exists {
				errorsJSON, _ = json.Marshal(data)
			}
		}
		json.Unmarshal(errorsJSON, &log.Errors)
	}

	// Parse logs JSON
	if model.Logs != nil {
		logsJSON, _ := json.Marshal(model.Logs)
		var logsData interface{}
		json.Unmarshal(logsJSON, &logsData)
		// If wrapped in {"data": array}, unwrap it
		if logsMap, ok := logsData.(map[string]interface{}); ok {
			if data, exists := logsMap["data"]; exists {
				logsJSON, _ = json.Marshal(data)
			}
		}
		json.Unmarshal(logsJSON, &log.Logs)
	}

	return log
}

