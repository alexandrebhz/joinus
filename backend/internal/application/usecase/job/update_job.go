package job

import (
	"context"
	"time"

	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/domain/service"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
)

type UpdateJobUseCase struct {
	jobRepo     repository.JobRepository
	startupRepo repository.StartupRepository
	authService *service.AuthorizationService
	logger      logger.Logger
}

func NewUpdateJobUseCase(
	jobRepo repository.JobRepository,
	startupRepo repository.StartupRepository,
	authService *service.AuthorizationService,
	logger logger.Logger,
) *UpdateJobUseCase {
	return &UpdateJobUseCase{
		jobRepo:     jobRepo,
		startupRepo: startupRepo,
		authService: authService,
		logger:      logger,
	}
}

func (uc *UpdateJobUseCase) Execute(ctx context.Context, input dto.UpdateJobInput, userID string) (*dto.JobOutput, error) {
	job, err := uc.jobRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, errors.NewNotFoundError("job")
	}

	// Check permission
	canManage, err := uc.authService.CanManageJobs(ctx, userID, job.StartupID)
	if err != nil || !canManage {
		return nil, errors.NewForbiddenError("you don't have permission to update this job")
	}

	// Update fields
	if input.Title != nil {
		job.Title = *input.Title
	}
	if input.Description != nil {
		job.Description = *input.Description
	}
	if input.Requirements != nil {
		job.Requirements = *input.Requirements
	}
	if input.JobType != nil {
		job.JobType = entity.JobType(*input.JobType)
	}
	if input.LocationType != nil {
		job.LocationType = entity.LocationType(*input.LocationType)
	}
	if input.City != nil {
		job.City = *input.City
	}
	if input.Country != nil {
		job.Country = *input.Country
	}
	if input.SalaryMin != nil {
		job.SalaryMin = input.SalaryMin
	}
	if input.SalaryMax != nil {
		job.SalaryMax = input.SalaryMax
	}
	if input.Currency != nil {
		job.Currency = *input.Currency
	}
	if input.ApplicationURL != nil {
		job.ApplicationURL = input.ApplicationURL
	}
	if input.ApplicationEmail != nil {
		job.ApplicationEmail = input.ApplicationEmail
	}
	if input.Status != nil {
		job.Status = entity.JobStatus(*input.Status)
	}
	if input.ExpiresAt != nil {
		if *input.ExpiresAt != "" {
			parsed, err := time.Parse(time.RFC3339, *input.ExpiresAt)
			if err != nil {
				return nil, errors.NewBadRequestError("invalid expires_at format")
			}
			job.ExpiresAt = &parsed
		} else {
			job.ExpiresAt = nil
		}
	}

	job.UpdatedAt = time.Now()

	if err := uc.jobRepo.Update(ctx, job); err != nil {
		return nil, err
	}

	startup, _ := uc.startupRepo.FindByID(ctx, job.StartupID)
	startupName := ""
	if startup != nil {
		startupName = startup.Name
	}

	return uc.toOutput(job, startupName), nil
}

func (uc *UpdateJobUseCase) toOutput(job *entity.Job, startupName string) *dto.JobOutput {
	output := &dto.JobOutput{
		ID:              job.ID,
		StartupID:       job.StartupID,
		StartupName:     startupName,
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



