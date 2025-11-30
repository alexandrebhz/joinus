package postgres

import (
	"context"
	"time"

	"github.com/startup-job-board/crawler/internal/domain/entity"
	"github.com/startup-job-board/crawler/internal/domain/repository"
	"github.com/startup-job-board/crawler/internal/infrastructure/persistence/gorm_model"
	"gorm.io/gorm"
)

// JobRepositoryImpl implements JobRepository using PostgreSQL
type JobRepositoryImpl struct {
	db *gorm.DB
}

// NewJobRepository creates a new JobRepositoryImpl
func NewJobRepository(db *gorm.DB) repository.JobRepository {
	return &JobRepositoryImpl{db: db}
}

// Create creates a new crawled job
func (r *JobRepositoryImpl) Create(ctx context.Context, job *entity.CrawledJob) error {
	model := r.toModel(job)
	return r.db.WithContext(ctx).Create(model).Error
}

// Update updates an existing crawled job
func (r *JobRepositoryImpl) Update(ctx context.Context, job *entity.CrawledJob) error {
	model := r.toModel(job)
	return r.db.WithContext(ctx).Save(model).Error
}

// FindByID finds a crawled job by ID
func (r *JobRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.CrawledJob, error) {
	var model gorm_model.CrawledJobModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return r.toEntity(&model), nil
}

// FindBySiteID finds all jobs for a specific site
func (r *JobRepositoryImpl) FindBySiteID(ctx context.Context, siteID string) ([]*entity.CrawledJob, error) {
	var models []gorm_model.CrawledJobModel
	if err := r.db.WithContext(ctx).Where("site_id = ?", siteID).Find(&models).Error; err != nil {
		return nil, err
	}

	jobs := make([]*entity.CrawledJob, len(models))
	for i, model := range models {
		jobs[i] = r.toEntity(&model)
	}
	return jobs, nil
}

// FindUnsynced finds unsynced jobs
func (r *JobRepositoryImpl) FindUnsynced(ctx context.Context, limit int) ([]*entity.CrawledJob, error) {
	var models []gorm_model.CrawledJobModel
	query := r.db.WithContext(ctx).Where("synced = ?", false)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	jobs := make([]*entity.CrawledJob, len(models))
	for i, model := range models {
		jobs[i] = r.toEntity(&model)
	}
	return jobs, nil
}

// ExistsByURL checks if a job with the given URL exists
func (r *JobRepositoryImpl) ExistsByURL(ctx context.Context, url string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&gorm_model.CrawledJobModel{}).
		Where("detail_url = ?", url).
		Count(&count).Error
	return count > 0, err
}

// ExistsByHash checks if a job with the given hash exists
func (r *JobRepositoryImpl) ExistsByHash(ctx context.Context, hash string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&gorm_model.CrawledJobModel{}).
		Where("deduplication_hash = ?", hash).
		Count(&count).Error
	return count > 0, err
}

// ExistsByExternalID checks if a job with the given external ID exists
func (r *JobRepositoryImpl) ExistsByExternalID(ctx context.Context, externalID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&gorm_model.CrawledJobModel{}).
		Where("external_id = ?", externalID).
		Count(&count).Error
	return count > 0, err
}

// MarkAsSynced marks a job as synced
func (r *JobRepositoryImpl) MarkAsSynced(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&gorm_model.CrawledJobModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"synced":    true,
			"synced_at": now,
		}).Error
}

// CountUnsynced counts unsynced jobs
func (r *JobRepositoryImpl) CountUnsynced(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&gorm_model.CrawledJobModel{}).
		Where("synced = ?", false).
		Count(&count).Error
	return count, err
}

// toModel converts entity to GORM model
func (r *JobRepositoryImpl) toModel(job *entity.CrawledJob) *gorm_model.CrawledJobModel {
	customFields := gorm_model.JSONBString{}
	if job.CustomFields != nil {
		customFields = gorm_model.JSONBString(job.CustomFields)
	}

	return &gorm_model.CrawledJobModel{
		ID:                job.ID,
		SiteID:            job.SiteID,
		ExternalID:        job.ExternalID,
		DetailURL:         job.DetailURL,
		Title:             job.Title,
		Description:       job.Description,
		Requirements:      job.Requirements,
		Company:           job.Company,
		Location:          job.Location,
		City:              job.City,
		Country:           job.Country,
		JobType:           job.JobType,
		LocationType:      job.LocationType,
		SalaryMin:         job.SalaryMin,
		SalaryMax:         job.SalaryMax,
		Currency:          job.Currency,
		ApplicationURL:    job.ApplicationURL,
		ApplicationEmail:  job.ApplicationEmail,
		ExpiresAt:         job.ExpiresAt,
		CustomFields:      customFields,
		RawHTML:           job.RawHTML,
		DeduplicationHash: job.DeduplicationHash,
		Synced:            job.Synced,
		SyncedAt:          job.SyncedAt,
		CreatedAt:         job.CreatedAt,
		UpdatedAt:         job.UpdatedAt,
	}
}

// toEntity converts GORM model to entity
func (r *JobRepositoryImpl) toEntity(model *gorm_model.CrawledJobModel) *entity.CrawledJob {
	customFields := make(map[string]string)
	if model.CustomFields != nil {
		customFields = map[string]string(model.CustomFields)
	}

	return &entity.CrawledJob{
		ID:                model.ID,
		SiteID:            model.SiteID,
		ExternalID:        model.ExternalID,
		DetailURL:         model.DetailURL,
		Title:             model.Title,
		Description:       model.Description,
		Requirements:      model.Requirements,
		Company:           model.Company,
		Location:          model.Location,
		City:              model.City,
		Country:           model.Country,
		JobType:           model.JobType,
		LocationType:      model.LocationType,
		SalaryMin:         model.SalaryMin,
		SalaryMax:         model.SalaryMax,
		Currency:          model.Currency,
		ApplicationURL:    model.ApplicationURL,
		ApplicationEmail:  model.ApplicationEmail,
		ExpiresAt:         model.ExpiresAt,
		CustomFields:      customFields,
		RawHTML:           model.RawHTML,
		DeduplicationHash: model.DeduplicationHash,
		Synced:            model.Synced,
		SyncedAt:          model.SyncedAt,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
	}
}

