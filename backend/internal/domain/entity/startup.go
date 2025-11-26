package entity

import "time"

type Startup struct {
	ID              string
	Name            string
	Slug            string
	Description     string
	LogoURL         *string
	Website         string
	FoundedYear     int
	Industry        string
	CompanySize     string
	Location        string
	SocialLinks     SocialLinks
	APIToken        string
	AllowPublicJoin bool
	JoinCode        *string
	Status          StartupStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

type SocialLinks struct {
	LinkedIn string
	Twitter  string
	GitHub   string
}

type StartupStatus string

const (
	StartupStatusActive   StartupStatus = "active"
	StartupStatusInactive StartupStatus = "inactive"
)



