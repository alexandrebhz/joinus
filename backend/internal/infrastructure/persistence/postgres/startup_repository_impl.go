package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"github.com/startup-job-board/backend/internal/infrastructure/persistence/gorm_model"
	"gorm.io/gorm"
)

type StartupRepositoryImpl struct {
	db *gorm.DB
}

func NewStartupRepository(db *gorm.DB) repository.StartupRepository {
	return &StartupRepositoryImpl{db: db}
}

func (r *StartupRepositoryImpl) Create(ctx context.Context, startup *entity.Startup) error {
	model := r.toModel(startup)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *StartupRepositoryImpl) Update(ctx context.Context, startup *entity.Startup) error {
	model := r.toModel(startup)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *StartupRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&gorm_model.Startup{}, id).Error
}

func (r *StartupRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Startup, error) {
	var model gorm_model.Startup
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *StartupRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*entity.Startup, error) {
	var model gorm_model.Startup
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *StartupRepositoryImpl) FindByAPIToken(ctx context.Context, token string) (*entity.Startup, error) {
	var model gorm_model.Startup
	if err := r.db.WithContext(ctx).Where("api_token = ?", token).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *StartupRepositoryImpl) List(ctx context.Context, filter repository.StartupFilter) ([]*entity.Startup, int64, error) {
	query := r.db.WithContext(ctx).Model(&gorm_model.Startup{})

	if filter.Industry != "" {
		query = query.Where("industry = ?", filter.Industry)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", string(filter.Status))
	}
	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", search, search)
	}
	if filter.Location != "" {
		location := "%" + strings.ToLower(filter.Location) + "%"
		query = query.Where("LOWER(location) LIKE ?", location)
	}
	if filter.CompanySize != "" {
		companySize := "%" + strings.ToLower(filter.CompanySize) + "%"
		query = query.Where("LOWER(company_size) LIKE ?", companySize)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// Apply ordering
	orderBy := "created_at"
	if filter.OrderBy != "" {
		orderBy = filter.OrderBy
	}
	orderDir := "DESC"
	if filter.OrderDir != "" {
		orderDir = strings.ToUpper(filter.OrderDir)
	}
	query = query.Order(fmt.Sprintf("%s %s", orderBy, orderDir))

	var models []gorm_model.Startup
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	startups := make([]*entity.Startup, len(models))
	for i, m := range models {
		startups[i] = r.toDomain(&m)
	}

	return startups, total, nil
}

func (r *StartupRepositoryImpl) toModel(startup *entity.Startup) *gorm_model.Startup {
	return &gorm_model.Startup{
		ID:              startup.ID,
		Name:            startup.Name,
		Slug:            startup.Slug,
		Description:     startup.Description,
		LogoURL:         startup.LogoURL,
		Website:         startup.Website,
		FoundedYear:     startup.FoundedYear,
		Industry:        startup.Industry,
		CompanySize:     startup.CompanySize,
		Location:        startup.Location,
		LinkedIn:        startup.SocialLinks.LinkedIn,
		Twitter:         startup.SocialLinks.Twitter,
		GitHub:          startup.SocialLinks.GitHub,
		APIToken:        startup.APIToken,
		AllowPublicJoin: startup.AllowPublicJoin,
		JoinCode:        startup.JoinCode,
		Status:          string(startup.Status),
		CreatedAt:       startup.CreatedAt,
		UpdatedAt:       startup.UpdatedAt,
	}
}

func (r *StartupRepositoryImpl) toDomain(model *gorm_model.Startup) *entity.Startup {
	return &entity.Startup{
		ID:              model.ID,
		Name:            model.Name,
		Slug:            model.Slug,
		Description:     model.Description,
		LogoURL:         model.LogoURL,
		Website:         model.Website,
		FoundedYear:     model.FoundedYear,
		Industry:        model.Industry,
		CompanySize:     model.CompanySize,
		Location:        model.Location,
		SocialLinks: entity.SocialLinks{
			LinkedIn: model.LinkedIn,
			Twitter:  model.Twitter,
			GitHub:   model.GitHub,
		},
		APIToken:        model.APIToken,
		AllowPublicJoin: model.AllowPublicJoin,
		JoinCode:        model.JoinCode,
		Status:          entity.StartupStatus(model.Status),
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}



