package job

import (
	"context"
	"time"

	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/logger"
)

type ListJobsUseCase struct {
	jobRepo     repository.JobRepository
	startupRepo repository.StartupRepository
	logger      logger.Logger
}

func NewListJobsUseCase(
	jobRepo repository.JobRepository,
	startupRepo repository.StartupRepository,
	logger logger.Logger,
) *ListJobsUseCase {
	return &ListJobsUseCase{
		jobRepo:     jobRepo,
		startupRepo: startupRepo,
		logger:      logger,
	}
}

func (uc *ListJobsUseCase) Execute(ctx context.Context, filter repository.JobFilter) ([]*dto.JobOutput, int64, error) {
	jobs, total, err := uc.jobRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]*dto.JobOutput, len(jobs))
	for i, j := range jobs {
		startup, _ := uc.startupRepo.FindByID(ctx, j.StartupID)
		startupName := ""
		startupSlug := ""
		if startup != nil {
			startupName = startup.Name
			startupSlug = startup.Slug
		}
		outputs[i] = uc.toOutput(j, startupName, startupSlug)
	}

	return outputs, total, nil
}

func (uc *ListJobsUseCase) toOutput(job *entity.Job, startupName string, startupSlug string) *dto.JobOutput {
	output := &dto.JobOutput{
		ID:              job.ID,
		StartupID:       job.StartupID,
		StartupName:     startupName,
		StartupSlug:     startupSlug,
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
		CreatedAt:       job.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       job.UpdatedAt.Format(time.RFC3339),
	}

	if job.ExpiresAt != nil {
		expiresAtStr := job.ExpiresAt.Format(time.RFC3339)
		output.ExpiresAt = &expiresAtStr
	}

	return output
}

