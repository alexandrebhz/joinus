package service

import (
	"context"
	"github.com/startup-job-board/backend/internal/domain/entity"
	"github.com/startup-job-board/backend/internal/domain/repository"
)

type AuthorizationService struct {
	memberRepo repository.StartupMemberRepository
}

func NewAuthorizationService(memberRepo repository.StartupMemberRepository) *AuthorizationService {
	return &AuthorizationService{
		memberRepo: memberRepo,
	}
}

func (s *AuthorizationService) CanManageStartup(ctx context.Context, userID, startupID string) (bool, error) {
	member, err := s.memberRepo.FindByUserAndStartup(ctx, userID, startupID)
	if err != nil {
		return false, err
	}
	if member == nil {
		return false, nil
	}
	return member.CanManageStartup(), nil
}

func (s *AuthorizationService) CanManageJobs(ctx context.Context, userID, startupID string) (bool, error) {
	member, err := s.memberRepo.FindByUserAndStartup(ctx, userID, startupID)
	if err != nil {
		return false, err
	}
	if member == nil {
		return false, nil
	}
	return member.CanManageJobs(), nil
}

func (s *AuthorizationService) CanManageMembers(ctx context.Context, userID, startupID string) (bool, error) {
	member, err := s.memberRepo.FindByUserAndStartup(ctx, userID, startupID)
	if err != nil {
		return false, err
	}
	if member == nil {
		return false, nil
	}
	return member.CanManageMembers(), nil
}

func (s *AuthorizationService) GetUserStartups(ctx context.Context, userID string) ([]*entity.StartupMember, error) {
	return s.memberRepo.FindByUserID(ctx, userID)
}

