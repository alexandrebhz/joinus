package gorm_model

import "time"

type Team struct {
	ID        string    `gorm:"type:uuid;primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Slug      string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	CreatedBy string    `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Team) TableName() string { return "teams" }

type TeamMember struct {
	ID        string     `gorm:"type:uuid;primary_key"`
	TeamID    string     `gorm:"type:uuid;not null;uniqueIndex:idx_team_user"`
	UserID    string     `gorm:"type:uuid;not null;uniqueIndex:idx_team_user;index"`
	RoleID    string     `gorm:"type:uuid;not null;index"`
	Status    string     `gorm:"type:varchar(50);not null;default:'active'"`
	InvitedBy *string    `gorm:"type:uuid"`
	InvitedAt time.Time
	JoinedAt  *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (TeamMember) TableName() string { return "team_members" }

type Role struct {
	ID        string    `gorm:"type:uuid;primary_key"`
	TeamID    *string   `gorm:"type:uuid;index"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Slug      string    `gorm:"type:varchar(100);not null;index"`
	IsSystem  bool      `gorm:"not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Role) TableName() string { return "roles" }

type RoleScope struct {
	ID     string `gorm:"type:uuid;primary_key"`
	RoleID string `gorm:"type:uuid;not null;uniqueIndex:idx_role_scope"`
	Scope  string `gorm:"type:varchar(100);not null;uniqueIndex:idx_role_scope"`
}

func (RoleScope) TableName() string { return "role_scopes" }

type TeamInvitation struct {
	ID        string    `gorm:"type:uuid;primary_key"`
	TeamID    string    `gorm:"type:uuid;not null;index"`
	Email     string    `gorm:"type:varchar(255);not null;index"`
	RoleID    string    `gorm:"type:uuid;not null"`
	Token     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	InvitedBy string    `gorm:"type:uuid;not null"`
	Status    string    `gorm:"type:varchar(50);not null;default:'pending'"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (TeamInvitation) TableName() string { return "team_invitations" }

type OAuthAccount struct {
	ID             string    `gorm:"type:uuid;primary_key"`
	UserID         string    `gorm:"type:uuid;not null;index"`
	Provider       string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_oauth_provider_user"`
	ProviderUserID string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_oauth_provider_user"`
	Email          string    `gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (OAuthAccount) TableName() string { return "oauth_accounts" }
