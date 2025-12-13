package gorm_model

import (
	"time"
)

type Contact struct {
	ID        string `gorm:"type:uuid;primary_key"`
	Name      string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"type:varchar(255);not null"`
	Subject   string `gorm:"type:varchar(500)"`
	Message   string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Contact) TableName() string {
	return "contacts"
}
