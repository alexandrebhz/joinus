package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/domain/service"
)

var errNotFound = errors.New("not found")

type fixture struct {
	authz   *service.AuthorizationService
	users   map[string]*entity.User
	members map[string]*entity.TeamMember
	roles   map[string]*entity.Role
	starts  map[string]*entity.Startup
}

func newFixture() *fixture {
	f := &fixture{
		users:   map[string]*entity.User{},
		members: map[string]*entity.TeamMember{},
		roles:   map[string]*entity.Role{},
		starts:  map[string]*entity.Startup{},
	}
	f.authz = service.NewAuthorizationService(
		&userRepo{f: f},
		&teamMemberRepo{f: f},
		&roleRepo{f: f},
		&startupRepo{f: f},
		&legacyMemberRepo{},
	)
	return f
}

func (f *fixture) addUser(id string, role entity.UserRole) {
	f.users[id] = &entity.User{ID: id, Role: role, Status: entity.UserStatusActive}
}

func (f *fixture) addRole(id string, scopes ...entity.Scope) {
	f.roles[id] = &entity.Role{ID: id, Slug: id, Scopes: scopes, IsSystem: true}
}

func (f *fixture) addMember(userID, teamID, roleID string) {
	f.members[userID+":"+teamID] = &entity.TeamMember{
		ID: userID + teamID, UserID: userID, TeamID: teamID, RoleID: roleID,
		Status: entity.MemberStatusActive, InvitedAt: time.Now(),
	}
}

func (f *fixture) addStartup(id string, teamID *string) {
	f.starts[id] = &entity.Startup{ID: id, TeamID: teamID, Status: entity.StartupStatusActive}
}

func TestPlatformAdminBypassesTeamMembership(t *testing.T) {
	f := newFixture()
	f.addUser("admin1", entity.UserRolePlatformAdmin)
	teamID := "team-a"
	f.addStartup("s1", &teamID)

	ok, err := f.authz.HasScope(context.Background(), "admin1", teamID, entity.ScopeJobsWrite)
	if err != nil || !ok {
		t.Fatalf("platform admin should have scope without membership, ok=%v err=%v", ok, err)
	}
	ok, err = f.authz.CanAccessStartup(context.Background(), "admin1", "s1", entity.ScopeStartupManage)
	if err != nil || !ok {
		t.Fatalf("platform admin should access any startup, ok=%v err=%v", ok, err)
	}
}

func TestLegacyAdminRoleIsPlatformAdmin(t *testing.T) {
	f := newFixture()
	f.addUser("legacy", entity.UserRoleAdmin)
	ok, err := f.authz.IsPlatformAdmin(context.Background(), "legacy")
	if err != nil || !ok {
		t.Fatalf("legacy admin should be platform admin, ok=%v err=%v", ok, err)
	}
}

func TestOrphanStartupDeniedToNonAdmin(t *testing.T) {
	f := newFixture()
	f.addUser("user1", entity.UserRoleCandidate)
	f.addStartup("orphan", nil)

	ok, err := f.authz.CanAccessStartup(context.Background(), "user1", "orphan", entity.ScopeStartupManage)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if ok {
		t.Fatal("non-admin must not access unlinked startup")
	}
}

func TestMemberScopeIsolationAcrossTeams(t *testing.T) {
	f := newFixture()
	f.addUser("alice", entity.UserRoleCandidate)
	f.addRole("recruiter", entity.ScopeJobsWrite, entity.ScopeJobsRead)
	f.addMember("alice", "team-a", "recruiter")
	teamA, teamB := "team-a", "team-b"
	f.addStartup("sa", &teamA)
	f.addStartup("sb", &teamB)

	ok, err := f.authz.CanAccessStartup(context.Background(), "alice", "sa", entity.ScopeJobsWrite)
	if err != nil || !ok {
		t.Fatalf("alice should manage jobs on her team startup, ok=%v err=%v", ok, err)
	}

	ok, err = f.authz.CanAccessStartup(context.Background(), "alice", "sb", entity.ScopeJobsWrite)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if ok {
		t.Fatal("alice must not access another team's startup (security isolation)")
	}

	ok, err = f.authz.HasScope(context.Background(), "alice", "team-b", entity.ScopeTeamRead)
	if err != nil || ok {
		t.Fatalf("alice must not read another team, ok=%v err=%v", ok, err)
	}
}

func TestInactiveMemberDenied(t *testing.T) {
	f := newFixture()
	f.addUser("bob", entity.UserRoleCandidate)
	f.addRole("owner", entity.AllScopes()...)
	f.addMember("bob", "team-a", "owner")
	f.members["bob:team-a"].Status = entity.MemberStatusRemoved

	ok, err := f.authz.HasScope(context.Background(), "bob", "team-a", entity.ScopeTeamManage)
	if err != nil || ok {
		t.Fatalf("removed member must be denied, ok=%v err=%v", ok, err)
	}
}

func TestRoleWithoutScopeDenied(t *testing.T) {
	f := newFixture()
	f.addUser("cara", entity.UserRoleCandidate)
	f.addRole("member", entity.ScopeTeamRead, entity.ScopeJobsRead)
	f.addMember("cara", "team-a", "member")

	ok, err := f.authz.HasScope(context.Background(), "cara", "team-a", entity.ScopeBillingManage)
	if err != nil || ok {
		t.Fatalf("member role must not have billing:manage, ok=%v err=%v", ok, err)
	}
}

// --- in-memory repos ---

type userRepo struct{ f *fixture }

func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
	r.f.users[user.ID] = user
	return nil
}
func (r *userRepo) Update(ctx context.Context, user *entity.User) error {
	r.f.users[user.ID] = user
	return nil
}
func (r *userRepo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	u, ok := r.f.users[id]
	if !ok {
		return nil, errNotFound
	}
	return u, nil
}
func (r *userRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, errNotFound
}
func (r *userRepo) List(ctx context.Context, page, pageSize int, search string) ([]*entity.User, int64, error) {
	return nil, 0, nil
}

type teamMemberRepo struct{ f *fixture }

func (r *teamMemberRepo) Create(ctx context.Context, m *entity.TeamMember) error {
	r.f.members[m.UserID+":"+m.TeamID] = m
	return nil
}
func (r *teamMemberRepo) Update(ctx context.Context, m *entity.TeamMember) error {
	r.f.members[m.UserID+":"+m.TeamID] = m
	return nil
}
func (r *teamMemberRepo) Delete(ctx context.Context, id string) error { return nil }
func (r *teamMemberRepo) FindByID(ctx context.Context, id string) (*entity.TeamMember, error) {
	return nil, errNotFound
}
func (r *teamMemberRepo) FindByUserAndTeam(ctx context.Context, userID, teamID string) (*entity.TeamMember, error) {
	m, ok := r.f.members[userID+":"+teamID]
	if !ok {
		return nil, nil
	}
	return m, nil
}
func (r *teamMemberRepo) FindByTeamID(ctx context.Context, teamID string) ([]*entity.TeamMember, error) {
	return nil, nil
}
func (r *teamMemberRepo) FindByUserID(ctx context.Context, userID string) ([]*entity.TeamMember, error) {
	return nil, nil
}

type roleRepo struct{ f *fixture }

func (r *roleRepo) Create(ctx context.Context, role *entity.Role) error {
	r.f.roles[role.ID] = role
	return nil
}
func (r *roleRepo) Update(ctx context.Context, role *entity.Role) error {
	r.f.roles[role.ID] = role
	return nil
}
func (r *roleRepo) FindByID(ctx context.Context, id string) (*entity.Role, error) {
	role, ok := r.f.roles[id]
	if !ok {
		return nil, errNotFound
	}
	return role, nil
}
func (r *roleRepo) FindSystemBySlug(ctx context.Context, slug string) (*entity.Role, error) {
	return nil, errNotFound
}
func (r *roleRepo) FindByTeamID(ctx context.Context, teamID string) ([]*entity.Role, error) {
	return nil, nil
}
func (r *roleRepo) ListSystem(ctx context.Context) ([]*entity.Role, error) { return nil, nil }
func (r *roleRepo) ReplaceScopes(ctx context.Context, roleID string, scopes []entity.Scope) error {
	return nil
}

type startupRepo struct{ f *fixture }

func (r *startupRepo) Create(ctx context.Context, s *entity.Startup) error {
	r.f.starts[s.ID] = s
	return nil
}
func (r *startupRepo) Update(ctx context.Context, s *entity.Startup) error {
	r.f.starts[s.ID] = s
	return nil
}
func (r *startupRepo) Delete(ctx context.Context, id string) error { return nil }
func (r *startupRepo) FindByID(ctx context.Context, id string) (*entity.Startup, error) {
	s, ok := r.f.starts[id]
	if !ok {
		return nil, errNotFound
	}
	return s, nil
}
func (r *startupRepo) FindBySlug(ctx context.Context, slug string) (*entity.Startup, error) {
	return nil, errNotFound
}
func (r *startupRepo) FindByAPIToken(ctx context.Context, token string) (*entity.Startup, error) {
	return nil, errNotFound
}
func (r *startupRepo) FindByStripeSubscriptionID(ctx context.Context, id string) (*entity.Startup, error) {
	return nil, errNotFound
}
func (r *startupRepo) FindByTeamID(ctx context.Context, teamID string) ([]*entity.Startup, error) {
	return nil, nil
}
func (r *startupRepo) List(ctx context.Context, filter repository.StartupFilter) ([]*entity.Startup, int64, error) {
	return nil, 0, nil
}

type legacyMemberRepo struct{}

func (r *legacyMemberRepo) Create(ctx context.Context, member *entity.StartupMember) error {
	return nil
}
func (r *legacyMemberRepo) Update(ctx context.Context, member *entity.StartupMember) error {
	return nil
}
func (r *legacyMemberRepo) Delete(ctx context.Context, id string) error { return nil }
func (r *legacyMemberRepo) FindByID(ctx context.Context, id string) (*entity.StartupMember, error) {
	return nil, errNotFound
}
func (r *legacyMemberRepo) FindByUserAndStartup(ctx context.Context, userID, startupID string) (*entity.StartupMember, error) {
	return nil, errNotFound
}
func (r *legacyMemberRepo) FindByStartupID(ctx context.Context, startupID string) ([]*entity.StartupMember, error) {
	return nil, nil
}
func (r *legacyMemberRepo) FindByUserID(ctx context.Context, userID string) ([]*entity.StartupMember, error) {
	return nil, nil
}
