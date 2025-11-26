package gorm_model

import (
	"time"
)

type StartupMember struct {
	ID        string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	StartupID string     `gorm:"type:uuid;not null;index"`
	UserID    string     `gorm:"type:uuid;not null;index"`
	Role      string     `gorm:"type:varchar(50);not null;default:'member'"`
	Status    string     `gorm:"type:varchar(50);not null;default:'active'"`
	InvitedBy *string    `gorm:"type:uuid"`
	InvitedAt time.Time  `gorm:"not null"`
	JoinedAt  *time.Time `gorm:"type:timestamp"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (StartupMember) TableName() string {
	return "startup_members"
}



