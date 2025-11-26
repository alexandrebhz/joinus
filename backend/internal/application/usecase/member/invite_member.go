package member

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/domain/service"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
)

type InviteMemberUseCase struct {
	invitationRepo repository.InvitationRepository
	memberRepo     repository.StartupMemberRepository
	startupRepo    repository.StartupRepository
	userRepo       repository.UserRepository
	emailService   port.EmailService
	tokenGen       port.TokenService
	appURL         string
	logger         logger.Logger
}

func NewInviteMemberUseCase(
	invitationRepo repository.InvitationRepository,
	memberRepo repository.StartupMemberRepository,
	startupRepo repository.StartupRepository,
	userRepo repository.UserRepository,
	emailService port.EmailService,
	tokenGen port.TokenService,
	appURL string,
	logger logger.Logger,
) *InviteMemberUseCase {
	return &InviteMemberUseCase{
		invitationRepo: invitationRepo,
		memberRepo:     memberRepo,
		startupRepo:    startupRepo,
		userRepo:       userRepo,
		emailService:   emailService,
		tokenGen:       tokenGen,
		appURL:         appURL,
		logger:         logger,
	}
}

func (uc *InviteMemberUseCase) Execute(ctx context.Context, input dto.InviteMemberInput, inviterID string) (*dto.InvitationOutput, error) {
	// Verify startup exists
	startup, err := uc.startupRepo.FindByID(ctx, input.StartupID)
	if err != nil {
		return nil, errors.NewNotFoundError("startup")
	}

	// Check if inviter has permission
	authService := service.NewAuthorizationService(uc.memberRepo)
	canManage, err := authService.CanManageMembers(ctx, inviterID, input.StartupID)
	if err != nil || !canManage {
		return nil, errors.NewForbiddenError("you don't have permission to invite members")
	}

	// Check if user already exists
	user, _ := uc.userRepo.FindByEmail(ctx, input.Email)
	if user != nil {
		// Check if already a member
		member, _ := uc.memberRepo.FindByUserAndStartup(ctx, user.ID, input.StartupID)
		if member != nil && member.Status == entity.MemberStatusActive {
			return nil, errors.NewBadRequestError("user is already a member of this startup")
		}
	}

	// Generate invitation token
	token, err := uc.tokenGen.GenerateInvitationToken()
	if err != nil {
		return nil, err
	}

	// Create invitation (expires in 24 hours)
	invitation := &entity.Invitation{
		ID:        uuid.New().String(),
		StartupID: input.StartupID,
		Email:     input.Email,
		Token:     token,
		Role:      entity.MemberRole(input.Role),
		InvitedBy: inviterID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Status:    entity.InvitationStatusPending,
		CreatedAt: time.Now(),
	}

	if err := uc.invitationRepo.Create(ctx, invitation); err != nil {
		return nil, err
	}

	// Send invitation email
	inviteURL := uc.appURL + "/invitations/" + token + "/accept"
	if err := uc.emailService.SendInvitationEmail(ctx, invitation, startup, inviteURL); err != nil {
		uc.logger.Warn("Failed to send invitation email: %v", err)
		// Don't fail the request if email fails
	}

	return &dto.InvitationOutput{
		ID:        invitation.ID,
		Email:     invitation.Email,
		Role:      string(invitation.Role),
		Token:     invitation.Token,
		ExpiresAt: invitation.ExpiresAt,
		InviteURL: inviteURL,
	}, nil
}



