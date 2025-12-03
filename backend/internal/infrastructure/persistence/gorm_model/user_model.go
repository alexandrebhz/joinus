package gorm_model

import (
	"time"
)

type User struct {
	ID        string    `gorm:"type:uuid;primary_key"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Role      string    `gorm:"type:varchar(50);not null;default:'candidate'"`
	StartupID *string   `gorm:"type:uuid"`
	Status    string    `gorm:"type:varchar(50);not null;default:'active'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "users"
}



