package repository

import (
	"context"

	"github.com/startup-job-board/backend/internal/domain/entity"
)

type ContactRepository interface {
	Create(ctx context.Context, contact *entity.Contact) error
}
