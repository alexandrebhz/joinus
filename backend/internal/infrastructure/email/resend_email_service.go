package email

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v3"
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/infrastructure/config"
)

type ResendEmailService struct {
	client *resend.Client
	from   string
}

func NewResendEmailService(cfg config.EmailConfig) port.EmailService {
	client := resend.NewClient(cfg.ResendAPIKey)
	return &ResendEmailService{
		client: client,
		from:   "noreply@startupboard.com", // Should be configurable
	}
}

func (s *ResendEmailService) SendInvitationEmail(ctx context.Context, invitation *entity.Invitation, startup *entity.Startup, inviteURL string) error {
	subject := fmt.Sprintf("You've been invited to join %s", startup.Name)
	htmlBody := fmt.Sprintf(`
		<h1>You've been invited!</h1>
		<p>You've been invited to join <strong>%s</strong> as a %s.</p>
		<p>Click the link below to accept the invitation:</p>
		<p><a href="%s">Accept Invitation</a></p>
		<p>This invitation expires in 24 hours.</p>
	`, startup.Name, invitation.Role, inviteURL)

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{invitation.Email},
		Subject: subject,
		Html:    htmlBody,
	}

	_, err := s.client.Emails.Send(params)
	return err
}

func (s *ResendEmailService) SendJoinRequestNotification(ctx context.Context, member *entity.StartupMember, startup *entity.Startup) error {
	// This would notify startup owners/admins about a join request
	// For now, we'll implement a simple version
	subject := fmt.Sprintf("New join request for %s", startup.Name)
	htmlBody := fmt.Sprintf(`
		<h1>New Join Request</h1>
		<p>A user has requested to join <strong>%s</strong>.</p>
		<p>Please review the request in your dashboard.</p>
	`, startup.Name)

	// Get startup owners/admins emails (this would need to be implemented)
	// For now, this is a placeholder
	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{"admin@startupboard.com"}, // Should be actual admin emails
		Subject: subject,
		Html:    htmlBody,
	}

	_, err := s.client.Emails.Send(params)
	return err
}

func (s *ResendEmailService) SendMemberApprovedEmail(ctx context.Context, member *entity.StartupMember, startup *entity.Startup) error {
	// Get user email (would need user repository)
	subject := fmt.Sprintf("Welcome to %s!", startup.Name)
	htmlBody := fmt.Sprintf(`
		<h1>Welcome!</h1>
		<p>Your request to join <strong>%s</strong> has been approved.</p>
		<p>You can now access the startup dashboard.</p>
	`, startup.Name)

	// This would need the user's email
	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{"user@example.com"}, // Should be actual user email
		Subject: subject,
		Html:    htmlBody,
	}

	_, err := s.client.Emails.Send(params)
	return err
}



