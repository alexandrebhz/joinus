package service

import (
	"context"
	"errors"

	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"gorm.io/gorm"
)

// AuthorizationService enforces team-scoped RBAC with platform-admin bypass.
// Platform admins are never required to be team members.
type AuthorizationService struct {
	userRepo       repository.UserRepository
	teamMemberRepo repository.TeamMemberRepository
	roleRepo       repository.RoleRepository
	startupRepo    repository.StartupRepository
	// Legacy membership used only as fallback while startups may still
	// have startup_members without a team link.
	memberRepo repository.StartupMemberRepository
}

func NewAuthorizationService(
	userRepo repository.UserRepository,
	teamMemberRepo repository.TeamMemberRepository,
	roleRepo repository.RoleRepository,
	startupRepo repository.StartupRepository,
	memberRepo repository.StartupMemberRepository,
) *AuthorizationService {
	return &AuthorizationService{
		userRepo:       userRepo,
		teamMemberRepo: teamMemberRepo,
		roleRepo:       roleRepo,
		startupRepo:    startupRepo,
		memberRepo:     memberRepo,
	}
}

func (s *AuthorizationService) IsPlatformAdmin(ctx context.Context, userID string) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return user.Role.IsPlatformAdmin(), nil
}

// HasScope checks whether the user has the given scope on a team.
// Platform admins always pass.
func (s *AuthorizationService) HasScope(ctx context.Context, userID, teamID string, scope entity.Scope) (bool, error) {
	ok, err := s.IsPlatformAdmin(ctx, userID)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	member, err := s.teamMemberRepo.FindByUserAndTeam(ctx, userID, teamID)
	if err != nil {
		return false, err
	}
	if member == nil || !member.IsActive() {
		return false, nil
	}

	role, err := s.roleRepo.FindByID(ctx, member.RoleID)
	if err != nil {
		return false, err
	}
	return role.HasScope(scope), nil
}

// CanAccessStartup evaluates scope against the startup's linked team.
// Unlinked startups (team_id null) are platform-admin only.
// Returns false (not an error) when denied — callers should use NotFound
// for private resources to avoid leaking existence across tenants.
func (s *AuthorizationService) CanAccessStartup(ctx context.Context, userID, startupID string, scope entity.Scope) (bool, error) {
	ok, err := s.IsPlatformAdmin(ctx, userID)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	startup, err := s.startupRepo.FindByID(ctx, startupID)
	if err != nil {
		return false, err
	}
	if startup.TeamID == nil || *startup.TeamID == "" {
		return false, nil
	}
	return s.HasScope(ctx, userID, *startup.TeamID, scope)
}

// CanManageJobs preserves the previous job API contract.
func (s *AuthorizationService) CanManageJobs(ctx context.Context, userID, startupID string) (bool, error) {
	ok, err := s.CanAccessStartup(ctx, userID, startupID, entity.ScopeJobsWrite)
	if err != nil || ok {
		return ok, err
	}
	return s.legacyCanManageJobs(ctx, userID, startupID)
}

func (s *AuthorizationService) CanManageMembers(ctx context.Context, userID, startupID string) (bool, error) {
	return s.CanAccessStartup(ctx, userID, startupID, entity.ScopeMembersManage)
}

func (s *AuthorizationService) CanManageStartup(ctx context.Context, userID, startupID string) (bool, error) {
	ok, err := s.CanAccessStartup(ctx, userID, startupID, entity.ScopeStartupManage)
	if err != nil || ok {
		return ok, err
	}
	return s.legacyCanManageStartup(ctx, userID, startupID)
}

func (s *AuthorizationService) GetUserStartups(ctx context.Context, userID string) ([]*entity.StartupMember, error) {
	return s.memberRepo.FindByUserID(ctx, userID)
}

func (s *AuthorizationService) legacyCanManageJobs(ctx context.Context, userID, startupID string) (bool, error) {
	member, err := s.memberRepo.FindByUserAndStartup(ctx, userID, startupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if member == nil {
		return false, nil
	}
	return member.CanManageJobs(), nil
}

func (s *AuthorizationService) legacyCanManageStartup(ctx context.Context, userID, startupID string) (bool, error) {
	member, err := s.memberRepo.FindByUserAndStartup(ctx, userID, startupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if member == nil {
		return false, nil
	}
	return member.CanManageStartup(), nil
}
