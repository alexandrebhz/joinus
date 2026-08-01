package entity

// Scope is a fine-grained permission string bound to roles.
type Scope string

const (
	ScopeTeamRead       Scope = "team:read"
	ScopeTeamManage     Scope = "team:manage"
	ScopeMembersRead    Scope = "members:read"
	ScopeMembersInvite  Scope = "members:invite"
	ScopeMembersManage  Scope = "members:manage"
	ScopeRolesRead      Scope = "roles:read"
	ScopeRolesManage    Scope = "roles:manage"
	ScopeStartupRead    Scope = "startup:read"
	ScopeStartupManage  Scope = "startup:manage"
	ScopeJobsRead       Scope = "jobs:read"
	ScopeJobsWrite      Scope = "jobs:write"
	ScopeJobsDelete     Scope = "jobs:delete"
	ScopeBillingRead    Scope = "billing:read"
	ScopeBillingManage  Scope = "billing:manage"
)

// AllScopes returns the full catalog of known scopes.
func AllScopes() []Scope {
	return []Scope{
		ScopeTeamRead, ScopeTeamManage,
		ScopeMembersRead, ScopeMembersInvite, ScopeMembersManage,
		ScopeRolesRead, ScopeRolesManage,
		ScopeStartupRead, ScopeStartupManage,
		ScopeJobsRead, ScopeJobsWrite, ScopeJobsDelete,
		ScopeBillingRead, ScopeBillingManage,
	}
}

// SystemRoleSlug identifies seeded system roles.
type SystemRoleSlug string

const (
	SystemRoleOwner     SystemRoleSlug = "owner"
	SystemRoleAdmin     SystemRoleSlug = "admin"
	SystemRoleMember    SystemRoleSlug = "member"
	SystemRoleRecruiter SystemRoleSlug = "recruiter"
)

// DefaultScopesForRole returns the scope set for a system role template.
func DefaultScopesForRole(slug SystemRoleSlug) []Scope {
	switch slug {
	case SystemRoleOwner:
		return AllScopes()
	case SystemRoleAdmin:
		return []Scope{
			ScopeTeamRead, ScopeTeamManage,
			ScopeMembersRead, ScopeMembersInvite, ScopeMembersManage,
			ScopeRolesRead, ScopeRolesManage,
			ScopeStartupRead, ScopeStartupManage,
			ScopeJobsRead, ScopeJobsWrite, ScopeJobsDelete,
			ScopeBillingRead, ScopeBillingManage,
		}
	case SystemRoleRecruiter:
		return []Scope{
			ScopeTeamRead,
			ScopeMembersRead,
			ScopeStartupRead,
			ScopeJobsRead, ScopeJobsWrite, ScopeJobsDelete,
		}
	case SystemRoleMember:
		return []Scope{
			ScopeTeamRead,
			ScopeMembersRead,
			ScopeStartupRead,
			ScopeJobsRead,
		}
	default:
		return nil
	}
}
