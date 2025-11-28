package repository

import (
	"context"
	"github.com/startup-job-board/backend/internal/domain/entity"
)

type JobRepository interface {
	Create(ctx context.Context, job *entity.Job) error
	Update(ctx context.Context, job *entity.Job) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entity.Job, error)
	List(ctx context.Context, filter JobFilter) ([]*entity.Job, int64, error)
	FindByStartupID(ctx context.Context, startupID string, limit int) ([]*entity.Job, error)
}

type JobFilter struct {
	StartupID    string
	JobType      entity.JobType
	LocationType entity.LocationType
	Status       entity.JobStatus
	Search       string
	Country      string
	City         string
	SalaryMin    *int
	SalaryMax    *int
	Currency     string
	Pagination
}



