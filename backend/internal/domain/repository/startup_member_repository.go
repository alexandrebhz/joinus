package repository

import (
	"context"
	"github.com/startup-job-board/backend/internal/domain/entity"
)

type StartupMemberRepository interface {
	Create(ctx context.Context, member *entity.StartupMember) error
	Update(ctx context.Context, member *entity.StartupMember) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entity.StartupMember, error)
	FindByUserAndStartup(ctx context.Context, userID, startupID string) (*entity.StartupMember, error)
	FindByStartupID(ctx context.Context, startupID string) ([]*entity.StartupMember, error)
	FindByUserID(ctx context.Context, userID string) ([]*entity.StartupMember, error)
}



