package startup

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
	"github.com/startup-job-board/backend/pkg/utils"
)

type CreateStartupUseCase struct {
	startupRepo    repository.StartupRepository
	teamRepo       repository.TeamRepository
	teamMemberRepo repository.TeamMemberRepository
	roleRepo       repository.RoleRepository
	memberRepo     repository.StartupMemberRepository
	userRepo       repository.UserRepository
	tokenGen       port.TokenService
	logger         logger.Logger
}

func NewCreateStartupUseCase(
	startupRepo repository.StartupRepository,
	teamRepo repository.TeamRepository,
	teamMemberRepo repository.TeamMemberRepository,
	roleRepo repository.RoleRepository,
	memberRepo repository.StartupMemberRepository,
	userRepo repository.UserRepository,
	tokenGen port.TokenService,
	logger logger.Logger,
) *CreateStartupUseCase {
	return &CreateStartupUseCase{
		startupRepo: startupRepo, teamRepo: teamRepo, teamMemberRepo: teamMemberRepo,
		roleRepo: roleRepo, memberRepo: memberRepo, userRepo: userRepo,
		tokenGen: tokenGen, logger: logger,
	}
}

func (uc *CreateStartupUseCase) Execute(ctx context.Context, input dto.CreateStartupInput, userID string, userRole entity.UserRole) (*dto.StartupOutput, error) {
	// Platform admins create orphan startups (no team) via the admin API.
	// Regular users create a team + become owner + link the startup.
	if userRole.IsPlatformAdmin() {
		return nil, errors.NewForbiddenError("platform admins must create startups via admin API (without a team)")
	}

	tokenStr, err := uc.tokenGen.GenerateToken()
	if err != nil {
		return nil, err
	}

	ownerRole, err := uc.roleRepo.FindSystemBySlug(ctx, string(entity.SystemRoleOwner))
	if err != nil {
		return nil, errors.ErrInternalError
	}

	teamBase := utils.GenerateSlug(input.Name + "-team")
	teamSlug := utils.GenerateUniqueSlug(teamBase, func(s string) bool {
		_, e := uc.teamRepo.FindBySlug(ctx, s)
		return e == nil
	})
	now := time.Now()
	team := &entity.Team{
		ID: uuid.New().String(), Name: input.Name + " Team", Slug: teamSlug,
		CreatedBy: userID, CreatedAt: now, UpdatedAt: now,
	}
	if err := uc.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}
	joined := now
	if err := uc.teamMemberRepo.Create(ctx, &entity.TeamMember{
		ID: uuid.New().String(), TeamID: team.ID, UserID: userID,
		RoleID: ownerRole.ID, Status: entity.MemberStatusActive,
		InvitedAt: now, JoinedAt: &joined, CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		return nil, err
	}

	baseSlug := utils.GenerateSlug(input.Name)
	slug := utils.GenerateUniqueSlug(baseSlug, func(s string) bool {
		_, e := uc.startupRepo.FindBySlug(ctx, s)
		return e == nil
	})

	teamID := team.ID
	startup := &entity.Startup{
		ID: uuid.New().String(), Name: input.Name, Slug: slug, Description: input.Description,
		Website: input.Website, FoundedYear: input.FoundedYear, Industry: input.Industry,
		CompanySize: input.CompanySize, Location: input.Location, APIToken: utils.HashToken(tokenStr),
		AllowPublicJoin: input.AllowPublicJoin, Status: entity.StartupStatusActive,
		TeamID: &teamID, Plan: string(entity.StartupPlanFree), CreatedAt: now, UpdatedAt: now,
	}
	if err := uc.startupRepo.Create(ctx, startup); err != nil {
		return nil, err
	}

	// Keep legacy startup_members in sync for gradual migration / crawler paths.
	_ = uc.memberRepo.Create(ctx, &entity.StartupMember{
		ID: uuid.New().String(), StartupID: startup.ID, UserID: userID,
		Role: entity.MemberRoleOwner, Status: entity.MemberStatusActive,
		InvitedAt: now, JoinedAt: &joined, CreatedAt: now, UpdatedAt: now,
	})

	if userRole == entity.UserRoleCandidate {
		if user, err := uc.userRepo.FindByID(ctx, userID); err == nil && user != nil {
			user.Role = entity.UserRoleStartupOwner
			user.UpdatedAt = time.Now()
			_ = uc.userRepo.Update(ctx, user)
		}
	}

	out := uc.toOutput(startup)
	out.APIToken = tokenStr // plaintext shown once
	return out, nil
}

func (uc *CreateStartupUseCase) toOutput(startup *entity.Startup) *dto.StartupOutput {
	return &dto.StartupOutput{
		ID: startup.ID, Name: startup.Name, Slug: startup.Slug, Description: startup.Description,
		LogoURL: startup.LogoURL, Website: startup.Website, FoundedYear: startup.FoundedYear,
		Industry: startup.Industry, CompanySize: startup.CompanySize, Location: startup.Location,
		AllowPublicJoin: startup.AllowPublicJoin, Status: string(startup.Status), Plan: startup.Plan,
		CreatedAt: startup.CreatedAt.Format(time.RFC3339), UpdatedAt: startup.UpdatedAt.Format(time.RFC3339),
	}
}
