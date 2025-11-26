package startup

import (
	"context"
	"time"

	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
)

type UpdateStartupUseCase struct {
	startupRepo repository.StartupRepository
	logger      logger.Logger
}

func NewUpdateStartupUseCase(
	startupRepo repository.StartupRepository,
	logger logger.Logger,
) *UpdateStartupUseCase {
	return &UpdateStartupUseCase{
		startupRepo: startupRepo,
		logger:      logger,
	}
}

func (uc *UpdateStartupUseCase) Execute(ctx context.Context, input dto.UpdateStartupInput) (*dto.StartupOutput, error) {
	startup, err := uc.startupRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, errors.NewNotFoundError("startup")
	}

	// Update fields
	if input.Name != nil {
		startup.Name = *input.Name
	}
	if input.Description != nil {
		startup.Description = *input.Description
	}
	if input.LogoURL != nil {
		startup.LogoURL = input.LogoURL
	}
	if input.Website != nil {
		startup.Website = *input.Website
	}
	if input.AllowPublicJoin != nil {
		startup.AllowPublicJoin = *input.AllowPublicJoin
	}
	if input.JoinCode != nil {
		startup.JoinCode = input.JoinCode
	}
	if input.Industry != nil {
		startup.Industry = *input.Industry
	}
	if input.CompanySize != nil {
		startup.CompanySize = *input.CompanySize
	}
	if input.Location != nil {
		startup.Location = *input.Location
	}

	startup.UpdatedAt = time.Now()

	if err := uc.startupRepo.Update(ctx, startup); err != nil {
		return nil, err
	}

	return uc.toOutput(startup), nil
}

func (uc *UpdateStartupUseCase) toOutput(startup *entity.Startup) *dto.StartupOutput {
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



