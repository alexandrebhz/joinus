package entity_test

import (
	"testing"

	"github.com/startup-job-board/backend/internal/domain/entity"
)

func TestDefaultScopesForOwnerIncludesBilling(t *testing.T) {
	scopes := entity.DefaultScopesForRole(entity.SystemRoleOwner)
	found := false
	for _, s := range scopes {
		if s == entity.ScopeBillingManage {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("owner should include billing:manage")
	}
}

func TestDefaultScopesForMemberExcludesWrite(t *testing.T) {
	scopes := entity.DefaultScopesForRole(entity.SystemRoleMember)
	role := &entity.Role{Scopes: scopes}
	if role.HasScope(entity.ScopeJobsWrite) {
		t.Fatal("member must not have jobs:write")
	}
	if !role.HasScope(entity.ScopeJobsRead) {
		t.Fatal("member should have jobs:read")
	}
}

func TestUserRoleIsPlatformAdmin(t *testing.T) {
	if !entity.UserRolePlatformAdmin.IsPlatformAdmin() {
		t.Fatal("platform_admin should be platform admin")
	}
	if !entity.UserRoleAdmin.IsPlatformAdmin() {
		t.Fatal("legacy admin should be platform admin")
	}
	if entity.UserRoleCandidate.IsPlatformAdmin() {
		t.Fatal("candidate must not be platform admin")
	}
}
