package dto

type CreateJobInput struct {
	StartupID        string  `json:"startup_id" validate:"required"`
	Title            string  `json:"title" validate:"required,min=5,max=100"`
	Description      string  `json:"description" validate:"required,min=20"`
	Requirements     string  `json:"requirements" validate:"required,min=10"`
	JobType          string  `json:"job_type" validate:"required,oneof=full_time part_time contract internship"`
	LocationType     string  `json:"location_type" validate:"required,oneof=remote hybrid onsite"`
	City             string  `json:"city"`
	Country          string  `json:"country" validate:"required"`
	SalaryMin        *int    `json:"salary_min" validate:"omitempty,min=0"`
	SalaryMax        *int    `json:"salary_max" validate:"omitempty,min=0"`
	Currency         string  `json:"currency" validate:"required,len=3"`
	ApplicationURL   *string `json:"application_url" validate:"omitempty,url"`
	ApplicationEmail *string `json:"application_email" validate:"omitempty,email"`
	ExpiresAt        *string `json:"expires_at" validate:"omitempty"`
}

type UpdateJobInput struct {
	ID               string  `json:"id"`
	Title            *string `json:"title" validate:"omitempty,min=5,max=100"`
	Description      *string `json:"description" validate:"omitempty,min=20"`
	Requirements     *string `json:"requirements" validate:"omitempty,min=10"`
	JobType          *string `json:"job_type" validate:"omitempty,oneof=full_time part_time contract internship"`
	LocationType     *string `json:"location_type" validate:"omitempty,oneof=remote hybrid onsite"`
	City             *string `json:"city"`
	Country          *string `json:"country"`
	SalaryMin        *int    `json:"salary_min" validate:"omitempty,min=0"`
	SalaryMax        *int    `json:"salary_max" validate:"omitempty,min=0"`
	Currency         *string `json:"currency" validate:"omitempty,len=3"`
	ApplicationURL   *string `json:"application_url" validate:"omitempty,url"`
	ApplicationEmail *string `json:"application_email" validate:"omitempty,email"`
	Status           *string `json:"status" validate:"omitempty,oneof=active closed"`
	ExpiresAt        *string `json:"expires_at" validate:"omitempty"`
}

type JobOutput struct {
	ID               string  `json:"id"`
	StartupID        string  `json:"startup_id"`
	StartupName      string  `json:"startup_name,omitempty"`
	StartupSlug      string  `json:"startup_slug,omitempty"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Requirements     string  `json:"requirements"`
	JobType          string  `json:"job_type"`
	LocationType     string  `json:"location_type"`
	City             string  `json:"city"`
	Country          string  `json:"country"`
	SalaryMin        *int    `json:"salary_min"`
	SalaryMax        *int    `json:"salary_max"`
	Currency         string  `json:"currency"`
	ApplicationURL   *string `json:"application_url"`
	ApplicationEmail *string `json:"application_email"`
	Status           string  `json:"status"`
	ExpiresAt        *string `json:"expires_at"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}
