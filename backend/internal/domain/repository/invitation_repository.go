package repository

import (
	"context"
	"github.com/startup-job-board/backend/internal/domain/entity"
)

type InvitationRepository interface {
	Create(ctx context.Context, invitation *entity.Invitation) error
	Update(ctx context.Context, invitation *entity.Invitation) error
	FindByID(ctx context.Context, id string) (*entity.Invitation, error)
	FindByToken(ctx context.Context, token string) (*entity.Invitation, error)
	FindByEmail(ctx context.Context, email string) ([]*entity.Invitation, error)
	FindPendingByStartup(ctx context.Context, startupID string) ([]*entity.Invitation, error)
	ExpireOldInvitations(ctx context.Context) error
}



