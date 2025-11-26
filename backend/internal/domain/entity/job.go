package entity

import "time"

type Job struct {
	ID              string
	StartupID       string
	Title           string
	Description     string
	Requirements    string
	JobType         JobType
	LocationType    LocationType
	City            string
	Country         string
	SalaryMin       *int
	SalaryMax       *int
	Currency        string
	ApplicationURL  *string
	ApplicationEmail *string
	Status          JobStatus
	ExpiresAt       *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type JobType string

const (
	JobTypeFullTime  JobType = "full_time"
	JobTypePartTime  JobType = "part_time"
	JobTypeContract  JobType = "contract"
	JobTypeInternship JobType = "internship"
)

type LocationType string

const (
	LocationRemote LocationType = "remote"
	LocationHybrid LocationType = "hybrid"
	LocationOnsite LocationType = "onsite"
)

type JobStatus string

const (
	JobStatusActive JobStatus = "active"
	JobStatusClosed JobStatus = "closed"
)



