package auth

import (
	"github.com/MedodsTechTask/app/user/auth/repo"
)

type AuthUseCase struct {
	repo repo.AuthRepo
}

func NewAuthUseCase(repo repo.AuthRepo) *AuthUseCase {
	return &AuthUseCase{repo}
}

func (s *AuthUseCase) HealthCheck() error {
	if err := s.repo.HealthCheck(); err != nil {

	}
}
