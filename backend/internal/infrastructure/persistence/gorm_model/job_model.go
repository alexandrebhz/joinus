package gorm_model

import (
	"time"
)

type Job struct {
	ID              string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	StartupID       string     `gorm:"type:uuid;not null;index"`
	Title           string     `gorm:"type:varchar(255);not null"`
	Description     string     `gorm:"type:text;not null"`
	Requirements    string     `gorm:"type:text;not null"`
	JobType         string     `gorm:"type:varchar(50);not null"`
	LocationType    string     `gorm:"type:varchar(50);not null"`
	City            string     `gorm:"type:varchar(255)"`
	Country         string     `gorm:"type:varchar(255);not null"`
	SalaryMin       *int       `gorm:"type:integer"`
	SalaryMax       *int       `gorm:"type:integer"`
	Currency        string     `gorm:"type:varchar(3);not null"`
	ApplicationURL  *string    `gorm:"type:varchar(500)"`
	ApplicationEmail *string   `gorm:"type:varchar(255)"`
	Status          string     `gorm:"type:varchar(50);not null;default:'active'"`
	ExpiresAt       *time.Time `gorm:"type:timestamp"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (Job) TableName() string {
	return "jobs"
}



