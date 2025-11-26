package entity

import "time"

type StartupMember struct {
	ID        string
	StartupID string
	UserID    string
	Role      MemberRole
	Status    MemberStatus
	InvitedBy *string
	InvitedAt time.Time
	JoinedAt  *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MemberRole string

const (
	MemberRoleOwner     MemberRole = "owner"
	MemberRoleAdmin     MemberRole = "admin"
	MemberRoleMember    MemberRole = "member"
	MemberRoleRecruiter MemberRole = "recruiter" // Can only manage jobs
)

type MemberStatus string

const (
	MemberStatusActive  MemberStatus = "active"
	MemberStatusPending MemberStatus = "pending"
	MemberStatusRemoved MemberStatus = "removed"
)

// Domain methods
func (m *StartupMember) CanManageJobs() bool {
	return m.Status == MemberStatusActive &&
		(m.Role == MemberRoleOwner || m.Role == MemberRoleAdmin || m.Role == MemberRoleRecruiter)
}

func (m *StartupMember) CanManageMembers() bool {
	return m.Status == MemberStatusActive &&
		(m.Role == MemberRoleOwner || m.Role == MemberRoleAdmin)
}

func (m *StartupMember) CanManageStartup() bool {
	return m.Status == MemberStatusActive &&
		(m.Role == MemberRoleOwner || m.Role == MemberRoleAdmin)
}



