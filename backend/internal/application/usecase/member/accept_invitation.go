package member

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/application/dto"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/pkg/errors"
	"github.com/startup-job-board/backend/pkg/logger"
)

type AcceptInvitationUseCase struct {
	invitationRepo repository.InvitationRepository
	memberRepo     repository.StartupMemberRepository
	userRepo       repository.UserRepository
	logger         logger.Logger
}

func NewAcceptInvitationUseCase(
	invitationRepo repository.InvitationRepository,
	memberRepo repository.StartupMemberRepository,
	userRepo repository.UserRepository,
	logger logger.Logger,
) *AcceptInvitationUseCase {
	return &AcceptInvitationUseCase{
		invitationRepo: invitationRepo,
		memberRepo:     memberRepo,
		userRepo:       userRepo,
		logger:         logger,
	}
}

func (uc *AcceptInvitationUseCase) Execute(ctx context.Context, token string, userID string) (*dto.MemberOutput, error) {
	// Find invitation
	invitation, err := uc.invitationRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, errors.NewNotFoundError("invitation")
	}

	// Validate invitation
	if !invitation.IsValid() {
		return nil, errors.NewBadRequestError("invitation is expired or already used")
	}

	// Verify user email matches invitation email
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewNotFoundError("user")
	}

	if user.Email != invitation.Email {
		return nil, errors.NewForbiddenError("this invitation was sent to a different email address")
	}

	// Check if already a member
	existing, _ := uc.memberRepo.FindByUserAndStartup(ctx, userID, invitation.StartupID)
	if existing != nil && existing.Status == entity.MemberStatusActive {
		return nil, errors.NewBadRequestError("you are already a member of this startup")
	}

	// Create or update member record
	now := time.Now()
	member := &entity.StartupMember{
		ID:        uuid.New().String(),
		StartupID: invitation.StartupID,
		UserID:    userID,
		Role:      invitation.Role,
		Status:    entity.MemberStatusActive,
		InvitedBy: &invitation.InvitedBy,
		InvitedAt: invitation.CreatedAt,
		JoinedAt:  &now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if existing != nil {
		// Update existing member
		member.ID = existing.ID
		member.CreatedAt = existing.CreatedAt
		if err := uc.memberRepo.Update(ctx, member); err != nil {
			return nil, err
		}
	} else {
		// Create new member
		if err := uc.memberRepo.Create(ctx, member); err != nil {
			return nil, err
		}
	}

	// Update invitation status
	invitation.Status = entity.InvitationStatusAccepted
	invitation.AcceptedAt = &now
	if err := uc.invitationRepo.Update(ctx, invitation); err != nil {
		uc.logger.Warn("Failed to update invitation status: %v", err)
	}

	return &dto.MemberOutput{
		ID:        member.ID,
		UserID:    member.UserID,
		StartupID: member.StartupID,
		Role:      string(member.Role),
		Status:    string(member.Status),
		User: dto.UserSummary{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		JoinedAt: member.JoinedAt,
	}, nil
}



