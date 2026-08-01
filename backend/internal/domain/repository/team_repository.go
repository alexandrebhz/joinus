package repository

import (
	"context"

	"github.com/startup-job-board/backend/internal/domain/entity"
)

type TeamRepository interface {
	Create(ctx context.Context, team *entity.Team) error
	Update(ctx context.Context, team *entity.Team) error
	FindByID(ctx context.Context, id string) (*entity.Team, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Team, error)
	ListByUserID(ctx context.Context, userID string) ([]*entity.Team, error)
	List(ctx context.Context, page, pageSize int) ([]*entity.Team, int64, error)
}

type TeamMemberRepository interface {
	Create(ctx context.Context, member *entity.TeamMember) error
	Update(ctx context.Context, member *entity.TeamMember) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entity.TeamMember, error)
	FindByUserAndTeam(ctx context.Context, userID, teamID string) (*entity.TeamMember, error)
	FindByTeamID(ctx context.Context, teamID string) ([]*entity.TeamMember, error)
	FindByUserID(ctx context.Context, userID string) ([]*entity.TeamMember, error)
}

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, role *entity.Role) error
	FindByID(ctx context.Context, id string) (*entity.Role, error)
	FindSystemBySlug(ctx context.Context, slug string) (*entity.Role, error)
	FindByTeamID(ctx context.Context, teamID string) ([]*entity.Role, error)
	ListSystem(ctx context.Context) ([]*entity.Role, error)
	ReplaceScopes(ctx context.Context, roleID string, scopes []entity.Scope) error
}

type TeamInvitationRepository interface {
	Create(ctx context.Context, invitation *entity.TeamInvitation) error
	Update(ctx context.Context, invitation *entity.TeamInvitation) error
	FindByToken(ctx context.Context, token string) (*entity.TeamInvitation, error)
	FindByTeamID(ctx context.Context, teamID string) ([]*entity.TeamInvitation, error)
}

type OAuthAccountRepository interface {
	Create(ctx context.Context, account *entity.OAuthAccount) error
	FindByProviderAndUserID(ctx context.Context, provider entity.OAuthProvider, providerUserID string) (*entity.OAuthAccount, error)
	FindByUserID(ctx context.Context, userID string) ([]*entity.OAuthAccount, error)
}
