package entity

import "time"

type Invitation struct {
	ID         string
	StartupID  string
	Email      string
	Token      string // Unique invitation token
	Role       MemberRole
	InvitedBy  string // UserID
	ExpiresAt  time.Time
	AcceptedAt *time.Time
	Status     InvitationStatus
	CreatedAt  time.Time
}

type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusExpired  InvitationStatus = "expired"
	InvitationStatusRevoked  InvitationStatus = "revoked"
)

func (i *Invitation) IsValid() bool {
	return i.Status == InvitationStatusPending &&
		i.ExpiresAt.After(time.Now())
}



