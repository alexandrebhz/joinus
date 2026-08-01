package admin

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/domain/service"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
	"github.com/startup-job-board/backend/pkg/utils"
)

type ListUsersUseCase struct {
	userRepo    repository.UserRepository
	authService *service.AuthorizationService
}

func NewListUsersUseCase(userRepo repository.UserRepository, authService *service.AuthorizationService) *ListUsersUseCase {
	return &ListUsersUseCase{userRepo: userRepo, authService: authService}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context, actorID string, page, pageSize int, search string) ([]dto.UserOutput, int64, error) {
	ok, err := uc.authService.IsPlatformAdmin(ctx, actorID)
	if err != nil || !ok {
		return nil, 0, errors.NewForbiddenError("platform admin required")
	}
	users, total, err := uc.userRepo.List(ctx, page, pageSize, search)
	if err != nil {
		return nil, 0, err
	}
	out := make([]dto.UserOutput, len(users))
	for i, u := range users {
		out[i] = dto.UserOutput{ID: u.ID, Email: u.Email, Name: u.Name, Role: string(u.Role)}
	}
	return out, total, nil
}

type UpdateUserUseCase struct {
	userRepo    repository.UserRepository
	authService *service.AuthorizationService
}

func NewUpdateUserUseCase(userRepo repository.UserRepository, authService *service.AuthorizationService) *UpdateUserUseCase {
	return &UpdateUserUseCase{userRepo: userRepo, authService: authService}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, actorID, targetID string, input dto.AdminUpdateUserInput) (*dto.UserOutput, error) {
	ok, err := uc.authService.IsPlatformAdmin(ctx, actorID)
	if err != nil || !ok {
		return nil, errors.NewForbiddenError("platform admin required")
	}
	user, err := uc.userRepo.FindByID(ctx, targetID)
	if err != nil {
		return nil, errors.NewNotFoundError("user")
	}
	if input.Role != "" {
		user.Role = entity.UserRole(input.Role)
	}
	if input.Status != "" {
		user.Status = entity.UserStatus(input.Status)
	}
	user.UpdatedAt = time.Now()
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return &dto.UserOutput{ID: user.ID, Email: user.Email, Name: user.Name, Role: string(user.Role)}, nil
}

type ListTeamsUseCase struct {
	teamRepo    repository.TeamRepository
	authService *service.AuthorizationService
}

func NewListTeamsUseCase(teamRepo repository.TeamRepository, authService *service.AuthorizationService) *ListTeamsUseCase {
	return &ListTeamsUseCase{teamRepo: teamRepo, authService: authService}
}

func (uc *ListTeamsUseCase) Execute(ctx context.Context, actorID string, page, pageSize int) ([]*dto.TeamOutput, int64, error) {
	ok, err := uc.authService.IsPlatformAdmin(ctx, actorID)
	if err != nil || !ok {
		return nil, 0, errors.NewForbiddenError("platform admin required")
	}
	teams, total, err := uc.teamRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	out := make([]*dto.TeamOutput, len(teams))
	for i, t := range teams {
		out[i] = &dto.TeamOutput{
			ID: t.ID, Name: t.Name, Slug: t.Slug,
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
			UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
		}
	}
	return out, total, nil
}

type CreateOrphanStartupUseCase struct {
	startupRepo repository.StartupRepository
	tokenGen    port.TokenService
	authService *service.AuthorizationService
	logger      logger.Logger
}

func NewCreateOrphanStartupUseCase(
	startupRepo repository.StartupRepository,
	tokenGen port.TokenService,
	authService *service.AuthorizationService,
	logger logger.Logger,
) *CreateOrphanStartupUseCase {
	return &CreateOrphanStartupUseCase{startupRepo: startupRepo, tokenGen: tokenGen, authService: authService, logger: logger}
}

func (uc *CreateOrphanStartupUseCase) Execute(ctx context.Context, actorID string, input dto.CreateStartupInput) (*dto.StartupOutput, error) {
	ok, err := uc.authService.IsPlatformAdmin(ctx, actorID)
	if err != nil || !ok {
		return nil, errors.NewForbiddenError("platform admin required")
	}
	tokenStr, err := uc.tokenGen.GenerateToken()
	if err != nil {
		return nil, err
	}
	baseSlug := utils.GenerateSlug(input.Name)
	slug := utils.GenerateUniqueSlug(baseSlug, func(s string) bool {
		_, e := uc.startupRepo.FindBySlug(ctx, s)
		return e == nil
	})
	now := time.Now()
	startup := &entity.Startup{
		ID: uuid.New().String(), Name: input.Name, Slug: slug, Description: input.Description,
		Website: input.Website, FoundedYear: input.FoundedYear, Industry: input.Industry,
		CompanySize: input.CompanySize, Location: input.Location, APIToken: utils.HashToken(tokenStr),
		AllowPublicJoin: input.AllowPublicJoin, Status: entity.StartupStatusActive,
		TeamID: nil, Plan: string(entity.StartupPlanFree), CreatedAt: now, UpdatedAt: now,
	}
	if err := uc.startupRepo.Create(ctx, startup); err != nil {
		return nil, err
	}
	return &dto.StartupOutput{
		ID: startup.ID, Name: startup.Name, Slug: startup.Slug, Description: startup.Description,
		Website: startup.Website, FoundedYear: startup.FoundedYear, Industry: startup.Industry,
		CompanySize: startup.CompanySize, Location: startup.Location,
		AllowPublicJoin: startup.AllowPublicJoin, Status: string(startup.Status), Plan: startup.Plan,
		APIToken: tokenStr, // plaintext shown once
		CreatedAt: startup.CreatedAt.Format(time.RFC3339), UpdatedAt: startup.UpdatedAt.Format(time.RFC3339),
	}, nil
}

type LinkStartupTeamUseCase struct {
	startupRepo repository.StartupRepository
	teamRepo    repository.TeamRepository
	authService *service.AuthorizationService
}

func NewLinkStartupTeamUseCase(
	startupRepo repository.StartupRepository,
	teamRepo repository.TeamRepository,
	authService *service.AuthorizationService,
) *LinkStartupTeamUseCase {
	return &LinkStartupTeamUseCase{startupRepo: startupRepo, teamRepo: teamRepo, authService: authService}
}

func (uc *LinkStartupTeamUseCase) Execute(ctx context.Context, actorID, startupID string, input dto.AdminLinkStartupTeamInput) error {
	ok, err := uc.authService.IsPlatformAdmin(ctx, actorID)
	if err != nil || !ok {
		return errors.NewForbiddenError("platform admin required")
	}
	startup, err := uc.startupRepo.FindByID(ctx, startupID)
	if err != nil {
		return errors.NewNotFoundError("startup")
	}
	if input.TeamID != nil && *input.TeamID != "" {
		if _, err := uc.teamRepo.FindByID(ctx, *input.TeamID); err != nil {
			return errors.NewNotFoundError("team")
		}
		startup.TeamID = input.TeamID
	} else {
		startup.TeamID = nil
	}
	startup.UpdatedAt = time.Now()
	return uc.startupRepo.Update(ctx, startup)
}
