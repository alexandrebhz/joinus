package dto

import "time"

type InviteMemberInput struct {
	StartupID string `json:"startup_id" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Role      string `json:"role" validate:"required,oneof=admin member recruiter"`
	Message   string `json:"message"`
}

type InvitationOutput struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	InviteURL string    `json:"invite_url"`
}

type AcceptInvitationInput struct {
	Token string `json:"token" validate:"required"`
}

type JoinRequestInput struct {
	StartupSlug string `json:"startup_slug" validate:"required"`
	JoinCode    string `json:"join_code"`
	Message     string `json:"message"`
}

type MemberOutput struct {
	ID        string       `json:"id"`
	UserID    string       `json:"user_id"`
	StartupID string       `json:"startup_id"`
	Role      string       `json:"role"`
	Status    string       `json:"status"`
	User      UserSummary  `json:"user"`
	JoinedAt  *time.Time   `json:"joined_at"`
}

type UserSummary struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateMemberRoleInput struct {
	MemberID string `json:"member_id" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=owner admin member recruiter"`
}



