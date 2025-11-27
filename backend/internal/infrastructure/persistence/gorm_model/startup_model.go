package gorm_model

import (
	"time"
	"gorm.io/gorm"
)

type Startup struct {
	ID              string     `gorm:"type:uuid;primary_key"`
	Name            string     `gorm:"type:varchar(255);not null"`
	Slug            string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	Description     string     `gorm:"type:text;not null"`
	LogoURL         *string    `gorm:"type:varchar(500)"`
	Website         string     `gorm:"type:varchar(500);not null"`
	FoundedYear     int        `gorm:"not null"`
	Industry        string     `gorm:"type:varchar(255);not null"`
	CompanySize     string     `gorm:"type:varchar(100);not null"`
	Location        string     `gorm:"type:varchar(255);not null"`
	LinkedIn        string     `gorm:"type:varchar(500)"`
	Twitter         string     `gorm:"type:varchar(500)"`
	GitHub          string     `gorm:"type:varchar(500)"`
	APIToken        string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	AllowPublicJoin bool       `gorm:"default:false"`
	JoinCode        *string    `gorm:"type:varchar(255)"`
	Status          string     `gorm:"type:varchar(50);not null;default:'active'"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (Startup) TableName() string {
	return "startups"
}



