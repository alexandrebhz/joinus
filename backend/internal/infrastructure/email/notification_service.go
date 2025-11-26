package email

import (
	"context"

	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/domain/entity"
)

type EmailNotificationService struct {
	emailService port.EmailService
}

func NewEmailNotificationService(emailService port.EmailService) port.NotificationService {
	return &EmailNotificationService{
		emailService: emailService,
	}
}

func (s *EmailNotificationService) NotifyJoinRequest(ctx context.Context, member *entity.StartupMember, startup *entity.Startup) error {
	return s.emailService.SendJoinRequestNotification(ctx, member, startup)
}



