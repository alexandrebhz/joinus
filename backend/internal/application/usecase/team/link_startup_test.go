package team_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/startup-job-board/backend/internal/application/usecase/team"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/domain/service"
	apperrors "github.com/startup-job-board/backend/pkg/errors"
)

var errNF = errors.New("not found")

type linkFixture struct {
	authz    *service.AuthorizationService
	startups map[string]*entity.Startup
	users    map[string]*entity.User
	members  map[string]*entity.TeamMember
	roles    map[string]*entity.Role
}

func newLinkFixture() *linkFixture {
	f := &linkFixture{
		startups: map[string]*entity.Startup{},
		users:    map[string]*entity.User{},
		members:  map[string]*entity.TeamMember{},
		roles:    map[string]*entity.Role{},
	}
	f.authz = service.NewAuthorizationService(
		&lfUser{f}, &lfMember{f}, &lfRole{f}, &lfStartup{f}, &lfLegacy{},
	)
	return f
}

func TestLinkStartupCannotStealOtherTeamStartup(t *testing.T) {
	f := newLinkFixture()
	f.users["u1"] = &entity.User{ID: "u1", Role: entity.UserRoleCandidate}
	f.roles["owner"] = &entity.Role{ID: "owner", Scopes: entity.AllScopes()}
	f.members["u1:team-a"] = &entity.TeamMember{
		UserID: "u1", TeamID: "team-a", RoleID: "owner", Status: entity.MemberStatusActive,
	}
	other := "team-b"
	f.startups["s1"] = &entity.Startup{ID: "s1", TeamID: &other}

	uc := team.NewLinkStartupUseCase(&lfStartup{f}, f.authz)
	err := uc.Execute(context.Background(), "team-a", "s1", "u1")
	if err == nil {
		t.Fatal("expected error when linking another team's startup")
	}
	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) || appErr.Code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND to avoid leaking ownership, got %v", err)
	}
	if f.startups["s1"].TeamID == nil || *f.startups["s1"].TeamID != "team-b" {
		t.Fatal("startup team_id must remain unchanged")
	}
}

func TestLinkStartupAllowsOrphan(t *testing.T) {
	f := newLinkFixture()
	f.users["u1"] = &entity.User{ID: "u1", Role: entity.UserRoleCandidate}
	f.roles["owner"] = &entity.Role{ID: "owner", Scopes: entity.AllScopes()}
	f.members["u1:team-a"] = &entity.TeamMember{
		UserID: "u1", TeamID: "team-a", RoleID: "owner", Status: entity.MemberStatusActive,
		InvitedAt: time.Now(),
	}
	f.startups["s1"] = &entity.Startup{ID: "s1", TeamID: nil}

	uc := team.NewLinkStartupUseCase(&lfStartup{f}, f.authz)
	if err := uc.Execute(context.Background(), "team-a", "s1", "u1"); err != nil {
		t.Fatalf("link orphan should succeed: %v", err)
	}
	if f.startups["s1"].TeamID == nil || *f.startups["s1"].TeamID != "team-a" {
		t.Fatal("orphan should be linked to team-a")
	}
}

func TestGetTeamHidesForeignTeam(t *testing.T) {
	f := newLinkFixture()
	f.users["u1"] = &entity.User{ID: "u1", Role: entity.UserRoleCandidate}
	teams := &lfTeam{byID: map[string]*entity.Team{
		"team-b": {ID: "team-b", Name: "Secret", Slug: "secret"},
	}}
	uc := team.NewGetTeamUseCase(teams, f.authz)
	_, err := uc.Execute(context.Background(), "team-b", "u1")
	if err == nil {
		t.Fatal("expected not found for foreign team")
	}
	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) || appErr.Code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND, got %v", err)
	}
}
type lfUser struct{ f *linkFixture }

func (r *lfUser) Create(ctx context.Context, user *entity.User) error { return nil }
func (r *lfUser) Update(ctx context.Context, user *entity.User) error { return nil }
func (r *lfUser) FindByID(ctx context.Context, id string) (*entity.User, error) {
	u, ok := r.f.users[id]
	if !ok {
		return nil, errNF
	}
	return u, nil
}
func (r *lfUser) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, errNF
}
func (r *lfUser) List(ctx context.Context, page, pageSize int, search string) ([]*entity.User, int64, error) {
	return nil, 0, nil
}

type lfMember struct{ f *linkFixture }

func (r *lfMember) Create(ctx context.Context, m *entity.TeamMember) error { return nil }
func (r *lfMember) Update(ctx context.Context, m *entity.TeamMember) error { return nil }
func (r *lfMember) Delete(ctx context.Context, id string) error             { return nil }
func (r *lfMember) FindByID(ctx context.Context, id string) (*entity.TeamMember, error) {
	return nil, errNF
}
func (r *lfMember) FindByUserAndTeam(ctx context.Context, userID, teamID string) (*entity.TeamMember, error) {
	m, ok := r.f.members[userID+":"+teamID]
	if !ok {
		return nil, nil
	}
	return m, nil
}
func (r *lfMember) FindByTeamID(ctx context.Context, teamID string) ([]*entity.TeamMember, error) {
	return nil, nil
}
func (r *lfMember) FindByUserID(ctx context.Context, userID string) ([]*entity.TeamMember, error) {
	return nil, nil
}

type lfRole struct{ f *linkFixture }

func (r *lfRole) Create(ctx context.Context, role *entity.Role) error { return nil }
func (r *lfRole) Update(ctx context.Context, role *entity.Role) error { return nil }
func (r *lfRole) FindByID(ctx context.Context, id string) (*entity.Role, error) {
	role, ok := r.f.roles[id]
	if !ok {
		return nil, errNF
	}
	return role, nil
}
func (r *lfRole) FindSystemBySlug(ctx context.Context, slug string) (*entity.Role, error) {
	return nil, errNF
}
func (r *lfRole) FindByTeamID(ctx context.Context, teamID string) ([]*entity.Role, error) {
	return nil, nil
}
func (r *lfRole) ListSystem(ctx context.Context) ([]*entity.Role, error) { return nil, nil }
func (r *lfRole) ReplaceScopes(ctx context.Context, roleID string, scopes []entity.Scope) error {
	return nil
}

type lfStartup struct{ f *linkFixture }

func (r *lfStartup) Create(ctx context.Context, s *entity.Startup) error { return nil }
func (r *lfStartup) Update(ctx context.Context, s *entity.Startup) error {
	r.f.startups[s.ID] = s
	return nil
}
func (r *lfStartup) Delete(ctx context.Context, id string) error { return nil }
func (r *lfStartup) FindByID(ctx context.Context, id string) (*entity.Startup, error) {
	s, ok := r.f.startups[id]
	if !ok {
		return nil, errNF
	}
	return s, nil
}
func (r *lfStartup) FindBySlug(ctx context.Context, slug string) (*entity.Startup, error) {
	return nil, errNF
}
func (r *lfStartup) FindByAPIToken(ctx context.Context, token string) (*entity.Startup, error) {
	return nil, errNF
}
func (r *lfStartup) FindByStripeSubscriptionID(ctx context.Context, id string) (*entity.Startup, error) {
	return nil, errNF
}
func (r *lfStartup) FindByTeamID(ctx context.Context, teamID string) ([]*entity.Startup, error) {
	return nil, nil
}
func (r *lfStartup) List(ctx context.Context, filter repository.StartupFilter) ([]*entity.Startup, int64, error) {
	return nil, 0, nil
}

type lfLegacy struct{}

func (r *lfLegacy) Create(ctx context.Context, member *entity.StartupMember) error { return nil }
func (r *lfLegacy) Update(ctx context.Context, member *entity.StartupMember) error { return nil }
func (r *lfLegacy) Delete(ctx context.Context, id string) error                     { return nil }
func (r *lfLegacy) FindByID(ctx context.Context, id string) (*entity.StartupMember, error) {
	return nil, errNF
}
func (r *lfLegacy) FindByUserAndStartup(ctx context.Context, userID, startupID string) (*entity.StartupMember, error) {
	return nil, errNF
}
func (r *lfLegacy) FindByStartupID(ctx context.Context, startupID string) ([]*entity.StartupMember, error) {
	return nil, nil
}
func (r *lfLegacy) FindByUserID(ctx context.Context, userID string) ([]*entity.StartupMember, error) {
	return nil, nil
}

type lfTeam struct {
	byID map[string]*entity.Team
}

func (r *lfTeam) Create(ctx context.Context, team *entity.Team) error { return nil }
func (r *lfTeam) Update(ctx context.Context, team *entity.Team) error { return nil }
func (r *lfTeam) FindByID(ctx context.Context, id string) (*entity.Team, error) {
	t, ok := r.byID[id]
	if !ok {
		return nil, errNF
	}
	return t, nil
}
func (r *lfTeam) FindBySlug(ctx context.Context, slug string) (*entity.Team, error) {
	return nil, errNF
}
func (r *lfTeam) ListByUserID(ctx context.Context, userID string) ([]*entity.Team, error) {
	return nil, nil
}
func (r *lfTeam) List(ctx context.Context, page, pageSize int) ([]*entity.Team, int64, error) {
	return nil, 0, nil
}
