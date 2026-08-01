package team

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
	"github.com/startup-job-board/backend/pkg/utils"
)

type CreateTeamUseCase struct {
	teamRepo       repository.TeamRepository
	teamMemberRepo repository.TeamMemberRepository
	roleRepo       repository.RoleRepository
	logger         logger.Logger
}

func NewCreateTeamUseCase(
	teamRepo repository.TeamRepository,
	teamMemberRepo repository.TeamMemberRepository,
	roleRepo repository.RoleRepository,
	logger logger.Logger,
) *CreateTeamUseCase {
	return &CreateTeamUseCase{teamRepo: teamRepo, teamMemberRepo: teamMemberRepo, roleRepo: roleRepo, logger: logger}
}

func (uc *CreateTeamUseCase) Execute(ctx context.Context, input dto.CreateTeamInput, userID string) (*dto.TeamOutput, error) {
	ownerRole, err := uc.roleRepo.FindSystemBySlug(ctx, string(entity.SystemRoleOwner))
	if err != nil {
		return nil, errors.ErrInternalError
	}

	baseSlug := utils.GenerateSlug(input.Name)
	slug := utils.GenerateUniqueSlug(baseSlug, func(s string) bool {
		_, e := uc.teamRepo.FindBySlug(ctx, s)
		return e == nil
	})

	now := time.Now()
	team := &entity.Team{
		ID: uuid.New().String(), Name: input.Name, Slug: slug,
		CreatedBy: userID, CreatedAt: now, UpdatedAt: now,
	}
	if err := uc.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}

	joined := now
	member := &entity.TeamMember{
		ID: uuid.New().String(), TeamID: team.ID, UserID: userID,
		RoleID: ownerRole.ID, Status: entity.MemberStatusActive,
		InvitedAt: now, JoinedAt: &joined, CreatedAt: now, UpdatedAt: now,
	}
	if err := uc.teamMemberRepo.Create(ctx, member); err != nil {
		uc.logger.Error("failed to create team owner membership: %v", err)
		return nil, err
	}

	return toTeamOutput(team), nil
}

type ListMyTeamsUseCase struct {
	teamRepo repository.TeamRepository
}

func NewListMyTeamsUseCase(teamRepo repository.TeamRepository) *ListMyTeamsUseCase {
	return &ListMyTeamsUseCase{teamRepo: teamRepo}
}

func (uc *ListMyTeamsUseCase) Execute(ctx context.Context, userID string) ([]*dto.TeamOutput, error) {
	teams, err := uc.teamRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]*dto.TeamOutput, len(teams))
	for i, t := range teams {
		out[i] = toTeamOutput(t)
	}
	return out, nil
}

type GetTeamUseCase struct {
	teamRepo   repository.TeamRepository
	authService *service.AuthorizationService
}

func NewGetTeamUseCase(teamRepo repository.TeamRepository, authService *service.AuthorizationService) *GetTeamUseCase {
	return &GetTeamUseCase{teamRepo: teamRepo, authService: authService}
}

func (uc *GetTeamUseCase) Execute(ctx context.Context, teamID, userID string) (*dto.TeamOutput, error) {
	// Security: use NOT_FOUND when the caller cannot see the team (no existence leak).
	ok, err := uc.authService.HasScope(ctx, userID, teamID, entity.ScopeTeamRead)
	if err != nil || !ok {
		return nil, errors.NewNotFoundError("team")
	}
	team, err := uc.teamRepo.FindByID(ctx, teamID)
	if err != nil {
		return nil, errors.NewNotFoundError("team")
	}
	return toTeamOutput(team), nil
}

type UpdateTeamUseCase struct {
	teamRepo    repository.TeamRepository
	authService *service.AuthorizationService
}

func NewUpdateTeamUseCase(teamRepo repository.TeamRepository, authService *service.AuthorizationService) *UpdateTeamUseCase {
	return &UpdateTeamUseCase{teamRepo: teamRepo, authService: authService}
}

func (uc *UpdateTeamUseCase) Execute(ctx context.Context, teamID, userID string, input dto.UpdateTeamInput) (*dto.TeamOutput, error) {
	ok, err := uc.authService.HasScope(ctx, userID, teamID, entity.ScopeTeamManage)
	if err != nil || !ok {
		return nil, errors.NewNotFoundError("team")
	}
	team, err := uc.teamRepo.FindByID(ctx, teamID)
	if err != nil {
		return nil, errors.NewNotFoundError("team")
	}
	team.Name = input.Name
	team.UpdatedAt = time.Now()
	if err := uc.teamRepo.Update(ctx, team); err != nil {
		return nil, err
	}
	return toTeamOutput(team), nil
}

func toTeamOutput(t *entity.Team) *dto.TeamOutput {
	return &dto.TeamOutput{
		ID: t.ID, Name: t.Name, Slug: t.Slug,
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
	}
}
