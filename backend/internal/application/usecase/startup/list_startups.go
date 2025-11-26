package startup

import (
	"context"
	"time"

	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/logger"
)

type ListStartupsUseCase struct {
	startupRepo repository.StartupRepository
	logger      logger.Logger
}

func NewListStartupsUseCase(
	startupRepo repository.StartupRepository,
	logger logger.Logger,
) *ListStartupsUseCase {
	return &ListStartupsUseCase{
		startupRepo: startupRepo,
		logger:      logger,
	}
}

func (uc *ListStartupsUseCase) Execute(ctx context.Context, filter repository.StartupFilter) ([]*dto.StartupOutput, int64, error) {
	startups, total, err := uc.startupRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	outputs := make([]*dto.StartupOutput, len(startups))
	for i, s := range startups {
		outputs[i] = uc.toOutput(s)
	}

	return outputs, total, nil
}

func (uc *ListStartupsUseCase) toOutput(startup *entity.Startup) *dto.StartupOutput {
	return &dto.StartupOutput{
		ID:              startup.ID,
		Name:            startup.Name,
		Slug:            startup.Slug,
		Description:     startup.Description,
		LogoURL:         startup.LogoURL,
		Website:         startup.Website,
		FoundedYear:     startup.FoundedYear,
		Industry:        startup.Industry,
		CompanySize:     startup.CompanySize,
		Location:        startup.Location,
		AllowPublicJoin: startup.AllowPublicJoin,
		Status:          string(startup.Status),
		CreatedAt:       startup.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       startup.UpdatedAt.Format(time.RFC3339),
	}
}



