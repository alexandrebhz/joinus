package gorm_model

import "time"

type OAuthLoginCode struct {
	ID           string     `gorm:"type:uuid;primary_key"`
	CodeHash     string     `gorm:"type:varchar(64);uniqueIndex;not null"`
	AccessToken  string     `gorm:"type:text;not null"`
	RefreshToken string     `gorm:"type:text;not null"`
	UserID       string     `gorm:"type:uuid;not null;index"`
	ExpiresAt    time.Time  `gorm:"not null;index"`
	UsedAt       *time.Time `gorm:"index"`
	CreatedAt    time.Time
}

func (OAuthLoginCode) TableName() string { return "oauth_login_codes" }
