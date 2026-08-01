package dto

type CreateTeamInput struct {
	Name string `json:"name" validate:"required,min=2"`
}

type UpdateTeamInput struct {
	Name string `json:"name" validate:"required,min=2"`
}

type TeamOutput struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TeamMemberOutput struct {
	ID        string   `json:"id"`
	TeamID    string   `json:"team_id"`
	UserID    string   `json:"user_id"`
	UserName  string   `json:"user_name,omitempty"`
	UserEmail string   `json:"user_email,omitempty"`
	RoleID    string   `json:"role_id"`
	RoleSlug  string   `json:"role_slug,omitempty"`
	Status    string   `json:"status"`
	Scopes    []string `json:"scopes,omitempty"`
}

type InviteTeamMemberInput struct {
	Email  string `json:"email" validate:"required,email"`
	RoleID string `json:"role_id" validate:"required"`
}

type UpdateTeamMemberInput struct {
	RoleID string `json:"role_id" validate:"required"`
	Status string `json:"status" validate:"omitempty,oneof=active pending removed"`
}

type AcceptTeamInvitationInput struct {
	Token string `json:"token" validate:"required"`
}

type LinkStartupInput struct {
	StartupID string `json:"startup_id" validate:"required"`
}

type RoleOutput struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Slug     string   `json:"slug"`
	IsSystem bool     `json:"is_system"`
	Scopes   []string `json:"scopes"`
}

type AdminUpdateUserInput struct {
	Role   string `json:"role" validate:"omitempty,oneof=platform_admin admin user candidate startup_owner"`
	Status string `json:"status" validate:"omitempty,oneof=active pending inactive"`
}

type AdminLinkStartupTeamInput struct {
	TeamID *string `json:"team_id"`
}

type MeOutput struct {
	ID     string            `json:"id"`
	Email  string            `json:"email"`
	Name   string            `json:"name"`
	Role   string            `json:"role"`
	Status string            `json:"status"`
	Teams  []MeTeamMembership `json:"teams"`
}

type MeTeamMembership struct {
	TeamID   string   `json:"team_id"`
	TeamName string   `json:"team_name"`
	RoleID   string   `json:"role_id"`
	RoleSlug string   `json:"role_slug"`
	Scopes   []string `json:"scopes"`
}
