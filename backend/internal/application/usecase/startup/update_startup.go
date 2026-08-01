package startup

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

type UpdateStartupUseCase struct {
	startupRepo repository.StartupRepository
	authService *service.AuthorizationService
	logger      logger.Logger
}

func NewUpdateStartupUseCase(
	startupRepo repository.StartupRepository,
	authService *service.AuthorizationService,
	logger logger.Logger,
) *UpdateStartupUseCase {
	return &UpdateStartupUseCase{
		startupRepo: startupRepo,
		authService: authService,
		logger:      logger,
	}
}

func (uc *UpdateStartupUseCase) Execute(ctx context.Context, input dto.UpdateStartupInput, userID string) (*dto.StartupOutput, error) {
	startup, err := uc.startupRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, errors.NewNotFoundError("startup")
	}

	// Security: deny with NOT_FOUND to avoid leaking existence of private startups.
	ok, err := uc.authService.CanAccessStartup(ctx, userID, input.ID, entity.ScopeStartupManage)
	if err != nil || !ok {
		return nil, errors.NewNotFoundError("startup")
	}

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
	output := &dto.StartupOutput{
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
		Plan:            startup.Plan,
		CreatedAt:       startup.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       startup.UpdatedAt.Format(time.RFC3339),
	}

	if startup.PlanExpiresAt != nil {
		planExpiresAtStr := startup.PlanExpiresAt.Format(time.RFC3339)
		output.PlanExpiresAt = &planExpiresAtStr
	}

	return output
}
