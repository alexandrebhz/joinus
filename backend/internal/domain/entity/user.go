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
	UserRoleAdmin        UserRole = "admin"
	UserRoleStartupOwner UserRole = "startup_owner"
	UserRoleCandidate    UserRole = "candidate"
	UserRoleMember       UserRole = "member" // Deprecated: use candidate instead
)

type UserStatus string

const (
	UserStatusActive  UserStatus = "active"
	UserStatusPending UserStatus = "pending"
	UserStatusInactive UserStatus = "inactive"
)



