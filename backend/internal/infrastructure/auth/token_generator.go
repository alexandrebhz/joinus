package auth

import (
	"github.com/startup-job-board/backend/internal/application/port"
	"github.com/startup-job-board/backend/internal/domain/valueobject"
)

type TokenGenerator struct{}

func NewTokenGenerator() port.TokenService {
	return &TokenGenerator{}
}

func (g *TokenGenerator) GenerateToken() (string, error) {
	token, err := valueobject.NewAPIToken("")
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

func (g *TokenGenerator) GenerateInvitationToken() (string, error) {
	token, err := valueobject.NewAPIToken("")
	if err != nil {
		return "", err
	}
	return token.String(), nil
}



