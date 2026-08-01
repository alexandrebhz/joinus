package seed

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
	"gorm.io/gorm"
)

// SystemRoles ensures global role templates + scopes exist (idempotent).
func SystemRoles(ctx context.Context, roleRepo repository.RoleRepository) error {
	slugs := []entity.SystemRoleSlug{
		entity.SystemRoleOwner,
		entity.SystemRoleAdmin,
		entity.SystemRoleMember,
		entity.SystemRoleRecruiter,
	}
	for _, slug := range slugs {
		existing, err := roleRepo.FindSystemBySlug(ctx, string(slug))
		if err == nil && existing != nil {
			continue
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		now := time.Now()
		role := &entity.Role{
			ID:        uuid.New().String(),
			TeamID:    nil,
			Name:      string(slug),
			Slug:      string(slug),
			IsSystem:  true,
			Scopes:    entity.DefaultScopesForRole(slug),
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := roleRepo.Create(ctx, role); err != nil {
			return err
		}
	}
	return nil
}
