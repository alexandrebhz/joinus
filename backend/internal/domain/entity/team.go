package entity

import "time"

type Team struct {
	ID        string
	Name      string
	Slug      string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TeamMember struct {
	ID        string
	TeamID    string
	UserID    string
	RoleID    string
	Status    MemberStatus
	InvitedBy *string
	InvitedAt time.Time
	JoinedAt  *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *TeamMember) IsActive() bool {
	return m.Status == MemberStatusActive
}

type Role struct {
	ID       string
	TeamID   *string // nil = global system template
	Name     string
	Slug     string
	IsSystem bool
	Scopes   []Scope
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Role) HasScope(scope Scope) bool {
	for _, s := range r.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

type TeamInvitation struct {
	ID        string
	TeamID    string
	Email     string
	RoleID    string
	Token     string
	InvitedBy string
	Status    InvitationStatus
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (i *TeamInvitation) IsValid() bool {
	return i.Status == InvitationStatusPending && time.Now().Before(i.ExpiresAt)
}
