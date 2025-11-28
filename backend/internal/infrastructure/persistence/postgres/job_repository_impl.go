package postgres

import (
	"context"
	"strings"

	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/gorm_model"
	"gorm.io/gorm"
)

type JobRepositoryImpl struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) repository.JobRepository {
	return &JobRepositoryImpl{db: db}
}

func (r *JobRepositoryImpl) Create(ctx context.Context, job *entity.Job) error {
	model := r.toModel(job)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *JobRepositoryImpl) Update(ctx context.Context, job *entity.Job) error {
	model := r.toModel(job)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *JobRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&gorm_model.Job{}, id).Error
}

func (r *JobRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Job, error) {
	var model gorm_model.Job
	if err := r.db.WithContext(ctx).Where(&gorm_model.Job{ID: id}).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *JobRepositoryImpl) List(ctx context.Context, filter repository.JobFilter) ([]*entity.Job, int64, error) {
	query := r.db.WithContext(ctx).Model(&gorm_model.Job{})

	if filter.StartupID != "" {
		query = query.Where(&gorm_model.Job{StartupID: filter.StartupID})
	}
	if filter.JobType != "" {
		query = query.Where(&gorm_model.Job{JobType: string(filter.JobType)})
	}
	if filter.LocationType != "" {
		query = query.Where(&gorm_model.Job{LocationType: string(filter.LocationType)})
	}
	if filter.Status != "" {
		query = query.Where(&gorm_model.Job{Status: string(filter.Status)})
	}
	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", search, search)
	}
	if filter.Country != "" {
		searchCountry := "%" + strings.ToLower(filter.Country) + "%"
		query = query.Where("LOWER(country) LIKE ?", searchCountry)
	}
	if filter.City != "" {
		searchCity := "%" + strings.ToLower(filter.City) + "%"
		query = query.Where("LOWER(city) LIKE ?", searchCity)
	}
	if filter.SalaryMin != nil {
		// Job must have salary_max >= min OR salary_min >= min (if no max specified)
		query = query.Where("salary_max >= ? OR (salary_min >= ? AND salary_max IS NULL)", *filter.SalaryMin, *filter.SalaryMin)
	}
	if filter.SalaryMax != nil {
		// Job must have salary_min <= max OR no min specified
		query = query.Where("salary_min <= ? OR salary_min IS NULL", *filter.SalaryMax)
	}
	if filter.Currency != "" {
		query = query.Where(&gorm_model.Job{Currency: filter.Currency})
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// Apply ordering using GORM's Order method
	orderBy := "created_at"
	if filter.OrderBy != "" {
		orderBy = filter.OrderBy
	}
	orderDir := "DESC"
	if filter.OrderDir != "" {
		orderDir = strings.ToUpper(filter.OrderDir)
	}
	if orderDir == "ASC" {
		query = query.Order(orderBy + " ASC")
	} else {
		query = query.Order(orderBy + " DESC")
	}

	var models []gorm_model.Job
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	jobs := make([]*entity.Job, len(models))
	for i, m := range models {
		jobs[i] = r.toDomain(&m)
	}

	return jobs, total, nil
}

func (r *JobRepositoryImpl) FindByStartupID(ctx context.Context, startupID string, limit int) ([]*entity.Job, error) {
	var models []gorm_model.Job
	query := r.db.WithContext(ctx).Where(&gorm_model.Job{StartupID: startupID})
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	jobs := make([]*entity.Job, len(models))
	for i, m := range models {
		jobs[i] = r.toDomain(&m)
	}
	return jobs, nil
}

func (r *JobRepositoryImpl) toModel(job *entity.Job) *gorm_model.Job {
	return &gorm_model.Job{
		ID:              job.ID,
		StartupID:       job.StartupID,
		Title:           job.Title,
		Description:     job.Description,
		Requirements:    job.Requirements,
		JobType:         string(job.JobType),
		LocationType:    string(job.LocationType),
		City:            job.City,
		Country:         job.Country,
		SalaryMin:       job.SalaryMin,
		SalaryMax:       job.SalaryMax,
		Currency:        job.Currency,
		ApplicationURL:  job.ApplicationURL,
		ApplicationEmail: job.ApplicationEmail,
		Status:          string(job.Status),
		ExpiresAt:       job.ExpiresAt,
		CreatedAt:       job.CreatedAt,
		UpdatedAt:       job.UpdatedAt,
	}
}

func (r *JobRepositoryImpl) toDomain(model *gorm_model.Job) *entity.Job {
	return &entity.Job{
		ID:              model.ID,
		StartupID:       model.StartupID,
		Title:           model.Title,
		Description:     model.Description,
		Requirements:    model.Requirements,
		JobType:         entity.JobType(model.JobType),
		LocationType:    entity.LocationType(model.LocationType),
		City:            model.City,
		Country:         model.Country,
		SalaryMin:       model.SalaryMin,
		SalaryMax:       model.SalaryMax,
		Currency:        model.Currency,
		ApplicationURL:  model.ApplicationURL,
		ApplicationEmail: model.ApplicationEmail,
		Status:          entity.JobStatus(model.Status),
		ExpiresAt:       model.ExpiresAt,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}



