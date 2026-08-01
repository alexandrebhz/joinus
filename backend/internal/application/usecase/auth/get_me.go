package auth

import (
	"context"

	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
)

type GetMeUseCase struct {
	userRepo       repository.UserRepository
	teamRepo       repository.TeamRepository
	teamMemberRepo repository.TeamMemberRepository
	roleRepo       repository.RoleRepository
	logger         logger.Logger
}

func NewGetMeUseCase(
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
	teamMemberRepo repository.TeamMemberRepository,
	roleRepo repository.RoleRepository,
	logger logger.Logger,
) *GetMeUseCase {
	return &GetMeUseCase{
		userRepo: userRepo, teamRepo: teamRepo, teamMemberRepo: teamMemberRepo,
		roleRepo: roleRepo, logger: logger,
	}
}

func (uc *GetMeUseCase) Execute(ctx context.Context, userID string) (*dto.MeOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewNotFoundError("user")
	}

	members, err := uc.teamMemberRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	teams := make([]dto.MeTeamMembership, 0, len(members))
	for _, m := range members {
		if !m.IsActive() {
			continue
		}
		team, err := uc.teamRepo.FindByID(ctx, m.TeamID)
		if err != nil {
			continue
		}
		role, err := uc.roleRepo.FindByID(ctx, m.RoleID)
		if err != nil {
			continue
		}
		scopes := make([]string, len(role.Scopes))
		for i, s := range role.Scopes {
			scopes[i] = string(s)
		}
		teams = append(teams, dto.MeTeamMembership{
			TeamID: team.ID, TeamName: team.Name,
			RoleID: role.ID, RoleSlug: role.Slug, Scopes: scopes,
		})
	}

	return &dto.MeOutput{
		ID: user.ID, Email: user.Email, Name: user.Name,
		Role: string(user.Role), Status: string(user.Status), Teams: teams,
	}, nil
}
