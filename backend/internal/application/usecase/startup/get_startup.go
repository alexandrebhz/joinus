package startup

import (
	"context"

	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
	"time"
)

type GetStartupUseCase struct {
	startupRepo repository.StartupRepository
	logger      logger.Logger
}

func NewGetStartupUseCase(
	startupRepo repository.StartupRepository,
	logger logger.Logger,
) *GetStartupUseCase {
	return &GetStartupUseCase{
		startupRepo: startupRepo,
		logger:      logger,
	}
}

func (uc *GetStartupUseCase) Execute(ctx context.Context, id string) (*dto.StartupOutput, error) {
	startup, err := uc.startupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("startup")
	}

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
	}, nil
}

func (uc *GetStartupUseCase) ExecuteBySlug(ctx context.Context, slug string) (*dto.StartupOutput, error) {
	startup, err := uc.startupRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.NewNotFoundError("startup")
	}

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
	}, nil
}



