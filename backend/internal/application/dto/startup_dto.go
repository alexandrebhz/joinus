package dto

type CreateStartupInput struct {
	Name            string `json:"name" validate:"required,min=2,max=100"`
	Description     string `json:"description" validate:"required,min=10"`
	Website         string `json:"website" validate:"required,url"`
	FoundedYear     int    `json:"founded_year" validate:"required,min=1900,max=2100"`
	Industry        string `json:"industry" validate:"required"`
	CompanySize     string `json:"company_size" validate:"required"`
	Location        string `json:"location" validate:"required"`
	AllowPublicJoin bool   `json:"allow_public_join"`
}

type UpdateStartupInput struct {
	ID              string  `json:"id"`
	Name            *string `json:"name" validate:"omitempty,min=2,max=100"`
	Description     *string `json:"description" validate:"omitempty,min=10"`
	LogoURL         *string `json:"logo_url" validate:"omitempty,url"`
	Website         *string `json:"website" validate:"omitempty,url"`
	AllowPublicJoin *bool   `json:"allow_public_join"`
	JoinCode        *string `json:"join_code"`
	Industry        *string `json:"industry"`
	CompanySize     *string `json:"company_size"`
	Location        *string `json:"location"`
}

type StartupOutput struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Slug            string  `json:"slug"`
	Description     string  `json:"description"`
	LogoURL         *string `json:"logo_url"`
	Website         string  `json:"website"`
	FoundedYear     int     `json:"founded_year"`
	Industry        string  `json:"industry"`
	CompanySize     string  `json:"company_size"`
	Location        string  `json:"location"`
	AllowPublicJoin bool    `json:"allow_public_join"`
	MemberCount     int     `json:"member_count,omitempty"`
	JobCount        int     `json:"job_count,omitempty"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}



