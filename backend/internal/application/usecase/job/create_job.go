package job

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/domain/service"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
)

type CreateJobUseCase struct {
	jobRepo      repository.JobRepository
	startupRepo  repository.StartupRepository
	memberRepo   repository.StartupMemberRepository
	authService  *service.AuthorizationService
	logger       logger.Logger
}

func NewCreateJobUseCase(
	jobRepo repository.JobRepository,
	startupRepo repository.StartupRepository,
	memberRepo repository.StartupMemberRepository,
	authService *service.AuthorizationService,
	logger logger.Logger,
) *CreateJobUseCase {
	return &CreateJobUseCase{
		jobRepo:     jobRepo,
		startupRepo: startupRepo,
		memberRepo:  memberRepo,
		authService: authService,
		logger:      logger,
	}
}

func (uc *CreateJobUseCase) Execute(ctx context.Context, input dto.CreateJobInput, userID string) (*dto.JobOutput, error) {
	// Verify startup exists
	startup, err := uc.startupRepo.FindByID(ctx, input.StartupID)
	if err != nil {
		return nil, errors.NewNotFoundError("startup")
	}

	// Check permission
	canManage, err := uc.authService.CanManageJobs(ctx, userID, input.StartupID)
	if err != nil || !canManage {
		return nil, errors.NewForbiddenError("you don't have permission to create jobs for this startup")
	}

	// Parse expires_at if provided
	var expiresAt *time.Time
	if input.ExpiresAt != nil && *input.ExpiresAt != "" {
		parsed, err := time.Parse(time.RFC3339, *input.ExpiresAt)
		if err != nil {
			return nil, errors.NewBadRequestError("invalid expires_at format")
		}
		expiresAt = &parsed
	}

	// Create job
	job := &entity.Job{
		ID:              uuid.New().String(),
		StartupID:       input.StartupID,
		Title:           input.Title,
		Description:     input.Description,
		Requirements:    input.Requirements,
		JobType:         entity.JobType(input.JobType),
		LocationType:    entity.LocationType(input.LocationType),
		City:            input.City,
		Country:         input.Country,
		SalaryMin:       input.SalaryMin,
		SalaryMax:       input.SalaryMax,
		Currency:        input.Currency,
		ApplicationURL:  input.ApplicationURL,
		ApplicationEmail: input.ApplicationEmail,
		Status:          entity.JobStatusActive,
		ExpiresAt:       expiresAt,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := uc.jobRepo.Create(ctx, job); err != nil {
		return nil, err
	}

	return uc.toOutput(job, startup.Name), nil
}

func (uc *CreateJobUseCase) toOutput(job *entity.Job, startupName string) *dto.JobOutput {
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



