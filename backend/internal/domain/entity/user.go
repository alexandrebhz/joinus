package entity

import "time"

type User struct {
	ID        string
	Email     string
	Password  string // Hashed
	Name      string
	Role      UserRole
	StartupID *string
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRole string

const (
	// PlatformAdmin has cross-tenant privilege and is never a team member for that power.
	UserRolePlatformAdmin UserRole = "platform_admin"
	// Legacy alias kept for existing JWTs / DB rows; treated as platform admin.
	UserRoleAdmin UserRole = "admin"
	// Regular authenticated user (job seeker and/or team member).
	UserRoleUser UserRole = "user"
	// Legacy roles — still accepted in JWT; map via helpers below.
	UserRoleStartupOwner UserRole = "startup_owner"
	UserRoleCandidate    UserRole = "candidate"
	UserRoleMember       UserRole = "member" // Deprecated
)

// IsPlatformAdmin reports whether the role has platform-wide privileges
// without requiring team membership.
func (r UserRole) IsPlatformAdmin() bool {
	return r == UserRolePlatformAdmin || r == UserRoleAdmin
}

type UserStatus string

const (
	UserStatusActive  UserStatus = "active"
	UserStatusPending UserStatus = "pending"
	UserStatusInactive UserStatus = "inactive"
)



