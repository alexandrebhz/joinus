package port

import (
	"context"
	"github.com/startup-job-board/backend/internal/domain/entity"
)

type EmailService interface {
	SendInvitationEmail(ctx context.Context, invitation *entity.Invitation, startup *entity.Startup, inviteURL string) error
	SendJoinRequestNotification(ctx context.Context, member *entity.StartupMember, startup *entity.Startup) error
	SendMemberApprovedEmail(ctx context.Context, member *entity.StartupMember, startup *entity.Startup) error
}



