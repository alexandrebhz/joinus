package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/gorm_model"
	"gorm.io/gorm"
)

type TeamRepositoryImpl struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) repository.TeamRepository {
	return &TeamRepositoryImpl{db: db}
}

func (r *TeamRepositoryImpl) Create(ctx context.Context, team *entity.Team) error {
	return r.db.WithContext(ctx).Create(&gorm_model.Team{
		ID: team.ID, Name: team.Name, Slug: team.Slug, CreatedBy: team.CreatedBy,
		CreatedAt: team.CreatedAt, UpdatedAt: team.UpdatedAt,
	}).Error
}

func (r *TeamRepositoryImpl) Update(ctx context.Context, team *entity.Team) error {
	return r.db.WithContext(ctx).Save(&gorm_model.Team{
		ID: team.ID, Name: team.Name, Slug: team.Slug, CreatedBy: team.CreatedBy,
		CreatedAt: team.CreatedAt, UpdatedAt: team.UpdatedAt,
	}).Error
}

func (r *TeamRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Team, error) {
	var m gorm_model.Team
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *TeamRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*entity.Team, error) {
	var m gorm_model.Team
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *TeamRepositoryImpl) ListByUserID(ctx context.Context, userID string) ([]*entity.Team, error) {
	var models []gorm_model.Team
	err := r.db.WithContext(ctx).
		Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ? AND team_members.status = ?", userID, string(entity.MemberStatusActive)).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	out := make([]*entity.Team, len(models))
	for i := range models {
		out[i] = r.toDomain(&models[i])
	}
	return out, nil
}

func (r *TeamRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*entity.Team, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&gorm_model.Team{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	var models []gorm_model.Team
	if err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]*entity.Team, len(models))
	for i := range models {
		out[i] = r.toDomain(&models[i])
	}
	return out, total, nil
}

func (r *TeamRepositoryImpl) toDomain(m *gorm_model.Team) *entity.Team {
	return &entity.Team{
		ID: m.ID, Name: m.Name, Slug: m.Slug, CreatedBy: m.CreatedBy,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

type TeamMemberRepositoryImpl struct {
	db *gorm.DB
}

func NewTeamMemberRepository(db *gorm.DB) repository.TeamMemberRepository {
	return &TeamMemberRepositoryImpl{db: db}
}

func (r *TeamMemberRepositoryImpl) Create(ctx context.Context, member *entity.TeamMember) error {
	return r.db.WithContext(ctx).Create(r.toModel(member)).Error
}

func (r *TeamMemberRepositoryImpl) Update(ctx context.Context, member *entity.TeamMember) error {
	return r.db.WithContext(ctx).Save(r.toModel(member)).Error
}

func (r *TeamMemberRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&gorm_model.TeamMember{}, "id = ?", id).Error
}

func (r *TeamMemberRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.TeamMember, error) {
	var m gorm_model.TeamMember
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *TeamMemberRepositoryImpl) FindByUserAndTeam(ctx context.Context, userID, teamID string) (*entity.TeamMember, error) {
	var m gorm_model.TeamMember
	if err := r.db.WithContext(ctx).Where("user_id = ? AND team_id = ?", userID, teamID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *TeamMemberRepositoryImpl) FindByTeamID(ctx context.Context, teamID string) ([]*entity.TeamMember, error) {
	var models []gorm_model.TeamMember
	if err := r.db.WithContext(ctx).Where("team_id = ?", teamID).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*entity.TeamMember, len(models))
	for i := range models {
		out[i] = r.toDomain(&models[i])
	}
	return out, nil
}

func (r *TeamMemberRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*entity.TeamMember, error) {
	var models []gorm_model.TeamMember
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*entity.TeamMember, len(models))
	for i := range models {
		out[i] = r.toDomain(&models[i])
	}
	return out, nil
}

func (r *TeamMemberRepositoryImpl) toModel(m *entity.TeamMember) *gorm_model.TeamMember {
	return &gorm_model.TeamMember{
		ID: m.ID, TeamID: m.TeamID, UserID: m.UserID, RoleID: m.RoleID,
		Status: string(m.Status), InvitedBy: m.InvitedBy, InvitedAt: m.InvitedAt,
		JoinedAt: m.JoinedAt, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func (r *TeamMemberRepositoryImpl) toDomain(m *gorm_model.TeamMember) *entity.TeamMember {
	return &entity.TeamMember{
		ID: m.ID, TeamID: m.TeamID, UserID: m.UserID, RoleID: m.RoleID,
		Status: entity.MemberStatus(m.Status), InvitedBy: m.InvitedBy, InvitedAt: m.InvitedAt,
		JoinedAt: m.JoinedAt, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

func (r *RoleRepositoryImpl) Create(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&gorm_model.Role{
			ID: role.ID, TeamID: role.TeamID, Name: role.Name, Slug: role.Slug,
			IsSystem: role.IsSystem, CreatedAt: role.CreatedAt, UpdatedAt: role.UpdatedAt,
		}).Error; err != nil {
			return err
		}
		return r.replaceScopesTx(tx, role.ID, role.Scopes)
	})
}

func (r *RoleRepositoryImpl) Update(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&gorm_model.Role{
			ID: role.ID, TeamID: role.TeamID, Name: role.Name, Slug: role.Slug,
			IsSystem: role.IsSystem, CreatedAt: role.CreatedAt, UpdatedAt: role.UpdatedAt,
		}).Error; err != nil {
			return err
		}
		return r.replaceScopesTx(tx, role.ID, role.Scopes)
	})
}

func (r *RoleRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Role, error) {
	var m gorm_model.Role
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toDomain(ctx, &m)
}

func (r *RoleRepositoryImpl) FindSystemBySlug(ctx context.Context, slug string) (*entity.Role, error) {
	var m gorm_model.Role
	if err := r.db.WithContext(ctx).Where("slug = ? AND is_system = ? AND team_id IS NULL", slug, true).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toDomain(ctx, &m)
}

func (r *RoleRepositoryImpl) FindByTeamID(ctx context.Context, teamID string) ([]*entity.Role, error) {
	var models []gorm_model.Role
	// Team-specific roles + global system templates
	if err := r.db.WithContext(ctx).
		Where("team_id = ? OR (is_system = ? AND team_id IS NULL)", teamID, true).
		Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*entity.Role, 0, len(models))
	for i := range models {
		role, err := r.toDomain(ctx, &models[i])
		if err != nil {
			return nil, err
		}
		out = append(out, role)
	}
	return out, nil
}

func (r *RoleRepositoryImpl) ListSystem(ctx context.Context) ([]*entity.Role, error) {
	var models []gorm_model.Role
	if err := r.db.WithContext(ctx).Where("is_system = ? AND team_id IS NULL", true).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*entity.Role, 0, len(models))
	for i := range models {
		role, err := r.toDomain(ctx, &models[i])
		if err != nil {
			return nil, err
		}
		out = append(out, role)
	}
	return out, nil
}

func (r *RoleRepositoryImpl) ReplaceScopes(ctx context.Context, roleID string, scopes []entity.Scope) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return r.replaceScopesTx(tx, roleID, scopes)
	})
}

func (r *RoleRepositoryImpl) replaceScopesTx(tx *gorm.DB, roleID string, scopes []entity.Scope) error {
	if err := tx.Where("role_id = ?", roleID).Delete(&gorm_model.RoleScope{}).Error; err != nil {
		return err
	}
	for _, s := range scopes {
		if err := tx.Create(&gorm_model.RoleScope{
			ID: uuid.New().String(), RoleID: roleID, Scope: string(s),
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *RoleRepositoryImpl) toDomain(ctx context.Context, m *gorm_model.Role) (*entity.Role, error) {
	var scopes []gorm_model.RoleScope
	if err := r.db.WithContext(ctx).Where("role_id = ?", m.ID).Find(&scopes).Error; err != nil {
		return nil, err
	}
	out := make([]entity.Scope, len(scopes))
	for i, s := range scopes {
		out[i] = entity.Scope(s.Scope)
	}
	return &entity.Role{
		ID: m.ID, TeamID: m.TeamID, Name: m.Name, Slug: m.Slug, IsSystem: m.IsSystem,
		Scopes: out, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}, nil
}

type TeamInvitationRepositoryImpl struct {
	db *gorm.DB
}

func NewTeamInvitationRepository(db *gorm.DB) repository.TeamInvitationRepository {
	return &TeamInvitationRepositoryImpl{db: db}
}

func (r *TeamInvitationRepositoryImpl) Create(ctx context.Context, inv *entity.TeamInvitation) error {
	return r.db.WithContext(ctx).Create(r.toModel(inv)).Error
}

func (r *TeamInvitationRepositoryImpl) Update(ctx context.Context, inv *entity.TeamInvitation) error {
	return r.db.WithContext(ctx).Save(r.toModel(inv)).Error
}

func (r *TeamInvitationRepositoryImpl) FindByToken(ctx context.Context, token string) (*entity.TeamInvitation, error) {
	var m gorm_model.TeamInvitation
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *TeamInvitationRepositoryImpl) FindByTeamID(ctx context.Context, teamID string) ([]*entity.TeamInvitation, error) {
	var models []gorm_model.TeamInvitation
	if err := r.db.WithContext(ctx).Where("team_id = ?", teamID).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*entity.TeamInvitation, len(models))
	for i := range models {
		out[i] = r.toDomain(&models[i])
	}
	return out, nil
}

func (r *TeamInvitationRepositoryImpl) toModel(inv *entity.TeamInvitation) *gorm_model.TeamInvitation {
	return &gorm_model.TeamInvitation{
		ID: inv.ID, TeamID: inv.TeamID, Email: inv.Email, RoleID: inv.RoleID,
		Token: inv.Token, InvitedBy: inv.InvitedBy, Status: string(inv.Status),
		ExpiresAt: inv.ExpiresAt, CreatedAt: inv.CreatedAt, UpdatedAt: inv.UpdatedAt,
	}
}

func (r *TeamInvitationRepositoryImpl) toDomain(m *gorm_model.TeamInvitation) *entity.TeamInvitation {
	return &entity.TeamInvitation{
		ID: m.ID, TeamID: m.TeamID, Email: m.Email, RoleID: m.RoleID,
		Token: m.Token, InvitedBy: m.InvitedBy, Status: entity.InvitationStatus(m.Status),
		ExpiresAt: m.ExpiresAt, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

type OAuthAccountRepositoryImpl struct {
	db *gorm.DB
}

func NewOAuthAccountRepository(db *gorm.DB) repository.OAuthAccountRepository {
	return &OAuthAccountRepositoryImpl{db: db}
}

func (r *OAuthAccountRepositoryImpl) Create(ctx context.Context, account *entity.OAuthAccount) error {
	return r.db.WithContext(ctx).Create(&gorm_model.OAuthAccount{
		ID: account.ID, UserID: account.UserID, Provider: string(account.Provider),
		ProviderUserID: account.ProviderUserID, Email: account.Email,
		CreatedAt: account.CreatedAt, UpdatedAt: account.UpdatedAt,
	}).Error
}

func (r *OAuthAccountRepositoryImpl) FindByProviderAndUserID(ctx context.Context, provider entity.OAuthProvider, providerUserID string) (*entity.OAuthAccount, error) {
	var m gorm_model.OAuthAccount
	if err := r.db.WithContext(ctx).Where("provider = ? AND provider_user_id = ?", string(provider), providerUserID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity.OAuthAccount{
		ID: m.ID, UserID: m.UserID, Provider: entity.OAuthProvider(m.Provider),
		ProviderUserID: m.ProviderUserID, Email: m.Email,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}, nil
}

func (r *OAuthAccountRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*entity.OAuthAccount, error) {
	var models []gorm_model.OAuthAccount
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*entity.OAuthAccount, len(models))
	for i, m := range models {
		out[i] = &entity.OAuthAccount{
			ID: m.ID, UserID: m.UserID, Provider: entity.OAuthProvider(m.Provider),
			ProviderUserID: m.ProviderUserID, Email: m.Email,
			CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
		}
	}
	return out, nil
}
