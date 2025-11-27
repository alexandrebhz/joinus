package gorm_model

import (
	"time"
)

type Invitation struct {
	ID        string     `gorm:"type:uuid;primary_key"`
	StartupID string     `gorm:"type:uuid;not null;index"`
	Email     string     `gorm:"type:varchar(255);not null;index"`
	Token     string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	Role      string     `gorm:"type:varchar(50);not null;default:'member'"`
	InvitedBy string     `gorm:"type:uuid;not null"`
	ExpiresAt time.Time  `gorm:"type:timestamp;not null"`
	AcceptedAt *time.Time `gorm:"type:timestamp"`
	Status    string     `gorm:"type:varchar(50);not null;default:'pending'"`
	CreatedAt time.Time
}

func (Invitation) TableName() string {
	return "invitations"
}



