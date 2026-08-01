package repository

import (
	"context"
	"time"

	"github.com/startup-job-board/backend/internal/domain/entity"
)

type OAuthLoginCodeRepository interface {
	Create(ctx context.Context, code *entity.OAuthLoginCode) error
	FindByCodeHash(ctx context.Context, codeHash string) (*entity.OAuthLoginCode, error)
	MarkUsed(ctx context.Context, id string, usedAt time.Time) error
	DeleteExpired(ctx context.Context) error
}
