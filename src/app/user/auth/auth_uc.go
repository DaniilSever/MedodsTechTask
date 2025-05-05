package auth

import "github.com/MedodsTechTask/app/user/auth/repo"

type AuthUseCase struct {
	repo repo.AuthRepo
}
