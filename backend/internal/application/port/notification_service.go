package port

import (
	"context"
	"github.com/startup-job-board/backend/internal/domain/entity"
)

type NotificationService interface {
	NotifyJoinRequest(ctx context.Context, member *entity.StartupMember, startup *entity.Startup) error
}



