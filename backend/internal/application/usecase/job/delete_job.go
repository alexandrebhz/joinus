package job

import (
	"context"

	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/domain/service"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
)

type DeleteJobUseCase struct {
	jobRepo     repository.JobRepository
	authService *service.AuthorizationService
	logger      logger.Logger
}

func NewDeleteJobUseCase(
	jobRepo repository.JobRepository,
	authService *service.AuthorizationService,
	logger logger.Logger,
) *DeleteJobUseCase {
	return &DeleteJobUseCase{
		jobRepo:     jobRepo,
		authService: authService,
		logger:      logger,
	}
}

func (uc *DeleteJobUseCase) Execute(ctx context.Context, jobID string, userID string) error {
	job, err := uc.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return errors.NewNotFoundError("job")
	}

	// Check permission
	canManage, err := uc.authService.CanManageJobs(ctx, userID, job.StartupID)
	if err != nil || !canManage {
		return errors.NewForbiddenError("you don't have permission to delete this job")
	}

	return uc.jobRepo.Delete(ctx, jobID)
}



