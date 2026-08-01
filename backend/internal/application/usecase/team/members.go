package team

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/domain/service"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
)

type ListMembersUseCase struct {
	teamMemberRepo repository.TeamMemberRepository
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	authService    *service.AuthorizationService
}

func NewListMembersUseCase(
	teamMemberRepo repository.TeamMemberRepository,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	authService *service.AuthorizationService,
) *ListMembersUseCase {
	return &ListMembersUseCase{
		teamMemberRepo: teamMemberRepo, userRepo: userRepo, roleRepo: roleRepo, authService: authService,
	}
}

func (uc *ListMembersUseCase) Execute(ctx context.Context, teamID, userID string) ([]*dto.TeamMemberOutput, error) {
	ok, err := uc.authService.HasScope(ctx, userID, teamID, entity.ScopeMembersRead)
	if err != nil || !ok {
		return nil, errors.NewNotFoundError("team")
	}
	members, err := uc.teamMemberRepo.FindByTeamID(ctx, teamID)
	if err != nil {
		return nil, err
	}
	out := make([]*dto.TeamMemberOutput, 0, len(members))
	for _, m := range members {
		item := &dto.TeamMemberOutput{
			ID: m.ID, TeamID: m.TeamID, UserID: m.UserID,
			RoleID: m.RoleID, Status: string(m.Status),
		}
		if user, err := uc.userRepo.FindByID(ctx, m.UserID); err == nil && user != nil {
			item.UserName = user.Name
			item.UserEmail = user.Email
		}
		if role, err := uc.roleRepo.FindByID(ctx, m.RoleID); err == nil && role != nil {
			item.RoleSlug = role.Slug
			item.Scopes = scopesToStrings(role.Scopes)
		}
		out = append(out, item)
	}
	return out, nil
}

type InviteMemberUseCase struct {
	invitationRepo repository.TeamInvitationRepository
	teamRepo       repository.TeamRepository
	teamMemberRepo repository.TeamMemberRepository
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	emailService   port.EmailService
	tokenGen       port.TokenService
	authService    *service.AuthorizationService
	appURL         string
	logger         logger.Logger
}

func NewInviteMemberUseCase(
	invitationRepo repository.TeamInvitationRepository,
	teamRepo repository.TeamRepository,
	teamMemberRepo repository.TeamMemberRepository,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	emailService port.EmailService,
	tokenGen port.TokenService,
	authService *service.AuthorizationService,
	appURL string,
	logger logger.Logger,
) *InviteMemberUseCase {
	return &InviteMemberUseCase{
		invitationRepo: invitationRepo, teamRepo: teamRepo, teamMemberRepo: teamMemberRepo,
		userRepo: userRepo, roleRepo: roleRepo, emailService: emailService,
		tokenGen: tokenGen, authService: authService, appURL: appURL, logger: logger,
	}
}

func (uc *InviteMemberUseCase) Execute(ctx context.Context, teamID, inviterID string, input dto.InviteTeamMemberInput) error {
	ok, err := uc.authService.HasScope(ctx, inviterID, teamID, entity.ScopeMembersInvite)
	if err != nil || !ok {
		return errors.NewNotFoundError("team")
	}

	role, err := uc.roleRepo.FindByID(ctx, input.RoleID)
	if err != nil {
		return errors.NewBadRequestError("invalid role")
	}
	// Do not allow assigning roles from another team.
	if role.TeamID != nil && *role.TeamID != teamID {
		return errors.NewBadRequestError("invalid role")
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))
	if existing, _ := uc.userRepo.FindByEmail(ctx, email); existing != nil {
		if member, _ := uc.teamMemberRepo.FindByUserAndTeam(ctx, existing.ID, teamID); member != nil && member.IsActive() {
			return errors.NewBadRequestError("user is already a team member")
		}
	}

	token, err := uc.tokenGen.GenerateInvitationToken()
	if err != nil {
		return err
	}
	now := time.Now()
	inv := &entity.TeamInvitation{
		ID: uuid.New().String(), TeamID: teamID, Email: email, RoleID: input.RoleID,
		Token: token, InvitedBy: inviterID, Status: entity.InvitationStatusPending,
		ExpiresAt: now.Add(24 * time.Hour), CreatedAt: now, UpdatedAt: now,
	}
	if err := uc.invitationRepo.Create(ctx, inv); err != nil {
		return err
	}

	team, err := uc.teamRepo.FindByID(ctx, teamID)
	teamName := "a JoinUs team"
	if err == nil && team != nil {
		teamName = team.Name
	}
	link := uc.appURL + "/invitations/accept?token=" + token
	if err := uc.emailService.SendTeamInvitationEmail(ctx, email, teamName, link); err != nil {
		uc.logger.Warn("failed to send team invitation email: %v", err)
	}
	return nil
}

type AcceptInvitationUseCase struct {
	invitationRepo repository.TeamInvitationRepository
	teamMemberRepo repository.TeamMemberRepository
	userRepo       repository.UserRepository
}

func NewAcceptInvitationUseCase(
	invitationRepo repository.TeamInvitationRepository,
	teamMemberRepo repository.TeamMemberRepository,
	userRepo repository.UserRepository,
) *AcceptInvitationUseCase {
	return &AcceptInvitationUseCase{invitationRepo: invitationRepo, teamMemberRepo: teamMemberRepo, userRepo: userRepo}
}

func (uc *AcceptInvitationUseCase) Execute(ctx context.Context, userID string, input dto.AcceptTeamInvitationInput) error {
	inv, err := uc.invitationRepo.FindByToken(ctx, input.Token)
	if err != nil || inv == nil || !inv.IsValid() {
		return errors.NewBadRequestError("invalid or expired invitation")
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.NewNotFoundError("user")
	}
	if !strings.EqualFold(user.Email, inv.Email) {
		// Do not reveal invitation email details.
		return errors.NewForbiddenError("invitation does not match this account")
	}

	if existing, _ := uc.teamMemberRepo.FindByUserAndTeam(ctx, userID, inv.TeamID); existing != nil {
		if existing.IsActive() {
			inv.Status = entity.InvitationStatusAccepted
			_ = uc.invitationRepo.Update(ctx, inv)
			return nil
		}
		existing.RoleID = inv.RoleID
		existing.Status = entity.MemberStatusActive
		now := time.Now()
		existing.JoinedAt = &now
		existing.UpdatedAt = now
		if err := uc.teamMemberRepo.Update(ctx, existing); err != nil {
			return err
		}
	} else {
		now := time.Now()
		member := &entity.TeamMember{
			ID: uuid.New().String(), TeamID: inv.TeamID, UserID: userID,
			RoleID: inv.RoleID, Status: entity.MemberStatusActive,
			InvitedBy: &inv.InvitedBy, InvitedAt: inv.CreatedAt, JoinedAt: &now,
			CreatedAt: now, UpdatedAt: now,
		}
		if err := uc.teamMemberRepo.Create(ctx, member); err != nil {
			return err
		}
	}

	inv.Status = entity.InvitationStatusAccepted
	inv.UpdatedAt = time.Now()
	return uc.invitationRepo.Update(ctx, inv)
}

type UpdateMemberUseCase struct {
	teamMemberRepo repository.TeamMemberRepository
	roleRepo       repository.RoleRepository
	authService    *service.AuthorizationService
}

func NewUpdateMemberUseCase(
	teamMemberRepo repository.TeamMemberRepository,
	roleRepo repository.RoleRepository,
	authService *service.AuthorizationService,
) *UpdateMemberUseCase {
	return &UpdateMemberUseCase{teamMemberRepo: teamMemberRepo, roleRepo: roleRepo, authService: authService}
}

func (uc *UpdateMemberUseCase) Execute(ctx context.Context, teamID, targetUserID, actorID string, input dto.UpdateTeamMemberInput) error {
	ok, err := uc.authService.HasScope(ctx, actorID, teamID, entity.ScopeMembersManage)
	if err != nil || !ok {
		return errors.NewNotFoundError("team")
	}
	member, err := uc.teamMemberRepo.FindByUserAndTeam(ctx, targetUserID, teamID)
	if err != nil || member == nil {
		return errors.NewNotFoundError("member")
	}
	role, err := uc.roleRepo.FindByID(ctx, input.RoleID)
	if err != nil {
		return errors.NewBadRequestError("invalid role")
	}
	if role.TeamID != nil && *role.TeamID != teamID {
		return errors.NewBadRequestError("invalid role")
	}
	member.RoleID = input.RoleID
	if input.Status != "" {
		member.Status = entity.MemberStatus(input.Status)
	}
	member.UpdatedAt = time.Now()
	return uc.teamMemberRepo.Update(ctx, member)
}

type RemoveMemberUseCase struct {
	teamMemberRepo repository.TeamMemberRepository
	authService    *service.AuthorizationService
}

func NewRemoveMemberUseCase(teamMemberRepo repository.TeamMemberRepository, authService *service.AuthorizationService) *RemoveMemberUseCase {
	return &RemoveMemberUseCase{teamMemberRepo: teamMemberRepo, authService: authService}
}

func (uc *RemoveMemberUseCase) Execute(ctx context.Context, teamID, targetUserID, actorID string) error {
	ok, err := uc.authService.HasScope(ctx, actorID, teamID, entity.ScopeMembersManage)
	if err != nil || !ok {
		return errors.NewNotFoundError("team")
	}
	if targetUserID == actorID {
		return errors.NewBadRequestError("cannot remove yourself")
	}
	member, err := uc.teamMemberRepo.FindByUserAndTeam(ctx, targetUserID, teamID)
	if err != nil || member == nil {
		return errors.NewNotFoundError("member")
	}
	member.Status = entity.MemberStatusRemoved
	member.UpdatedAt = time.Now()
	return uc.teamMemberRepo.Update(ctx, member)
}

type ListRolesUseCase struct {
	roleRepo    repository.RoleRepository
	authService *service.AuthorizationService
}

func NewListRolesUseCase(roleRepo repository.RoleRepository, authService *service.AuthorizationService) *ListRolesUseCase {
	return &ListRolesUseCase{roleRepo: roleRepo, authService: authService}
}

func (uc *ListRolesUseCase) Execute(ctx context.Context, teamID, userID string) ([]*dto.RoleOutput, error) {
	ok, err := uc.authService.HasScope(ctx, userID, teamID, entity.ScopeRolesRead)
	if err != nil || !ok {
		return nil, errors.NewNotFoundError("team")
	}
	roles, err := uc.roleRepo.FindByTeamID(ctx, teamID)
	if err != nil {
		return nil, err
	}
	out := make([]*dto.RoleOutput, len(roles))
	for i, r := range roles {
		out[i] = &dto.RoleOutput{
			ID: r.ID, Name: r.Name, Slug: r.Slug, IsSystem: r.IsSystem,
			Scopes: scopesToStrings(r.Scopes),
		}
	}
	return out, nil
}

type LinkStartupUseCase struct {
	startupRepo repository.StartupRepository
	authService *service.AuthorizationService
}

func NewLinkStartupUseCase(startupRepo repository.StartupRepository, authService *service.AuthorizationService) *LinkStartupUseCase {
	return &LinkStartupUseCase{startupRepo: startupRepo, authService: authService}
}

func (uc *LinkStartupUseCase) Execute(ctx context.Context, teamID, startupID, userID string) error {
	ok, err := uc.authService.HasScope(ctx, userID, teamID, entity.ScopeStartupManage)
	if err != nil || !ok {
		return errors.NewNotFoundError("team")
	}
	startup, err := uc.startupRepo.FindByID(ctx, startupID)
	if err != nil {
		return errors.NewNotFoundError("startup")
	}
	// Only link unassigned startups (or already linked to this team). Prevents stealing another team's startup.
	if startup.TeamID != nil && *startup.TeamID != "" && *startup.TeamID != teamID {
		return errors.NewNotFoundError("startup")
	}
	startup.TeamID = &teamID
	startup.UpdatedAt = time.Now()
	return uc.startupRepo.Update(ctx, startup)
}

type UnlinkStartupUseCase struct {
	startupRepo repository.StartupRepository
	authService *service.AuthorizationService
}

func NewUnlinkStartupUseCase(startupRepo repository.StartupRepository, authService *service.AuthorizationService) *UnlinkStartupUseCase {
	return &UnlinkStartupUseCase{startupRepo: startupRepo, authService: authService}
}

func (uc *UnlinkStartupUseCase) Execute(ctx context.Context, teamID, startupID, userID string) error {
	ok, err := uc.authService.HasScope(ctx, userID, teamID, entity.ScopeStartupManage)
	if err != nil || !ok {
		return errors.NewNotFoundError("team")
	}
	startup, err := uc.startupRepo.FindByID(ctx, startupID)
	if err != nil {
		return errors.NewNotFoundError("startup")
	}
	if startup.TeamID == nil || *startup.TeamID != teamID {
		return errors.NewNotFoundError("startup")
	}
	startup.TeamID = nil
	startup.UpdatedAt = time.Now()
	return uc.startupRepo.Update(ctx, startup)
}

type ListTeamStartupsUseCase struct {
	startupRepo repository.StartupRepository
	authService *service.AuthorizationService
}

func NewListTeamStartupsUseCase(startupRepo repository.StartupRepository, authService *service.AuthorizationService) *ListTeamStartupsUseCase {
	return &ListTeamStartupsUseCase{startupRepo: startupRepo, authService: authService}
}

func (uc *ListTeamStartupsUseCase) Execute(ctx context.Context, teamID, userID string) ([]*dto.StartupOutput, error) {
	ok, err := uc.authService.HasScope(ctx, userID, teamID, entity.ScopeStartupRead)
	if err != nil || !ok {
		return nil, errors.NewNotFoundError("team")
	}
	startups, err := uc.startupRepo.FindByTeamID(ctx, teamID)
	if err != nil {
		return nil, err
	}
	out := make([]*dto.StartupOutput, len(startups))
	for i, s := range startups {
		out[i] = &dto.StartupOutput{
			ID: s.ID, Name: s.Name, Slug: s.Slug, Description: s.Description,
			LogoURL: s.LogoURL, Website: s.Website, FoundedYear: s.FoundedYear,
			Industry: s.Industry, CompanySize: s.CompanySize, Location: s.Location,
			AllowPublicJoin: s.AllowPublicJoin, Status: string(s.Status), Plan: s.Plan,
			CreatedAt: s.CreatedAt.Format(time.RFC3339), UpdatedAt: s.UpdatedAt.Format(time.RFC3339),
		}
	}
	return out, nil
}

func scopesToStrings(scopes []entity.Scope) []string {
	out := make([]string, len(scopes))
	for i, s := range scopes {
		out[i] = string(s)
	}
	return out
}
