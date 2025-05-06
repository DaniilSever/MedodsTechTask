package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/MedodsTechTask/app/user/auth/repo"
	"github.com/MedodsTechTask/app/user/auth/share"
)

type AuthUseCase struct {
	repo repo.IAuthRepo
}

func NewAuthUseCase(repo repo.IAuthRepo) *AuthUseCase {
	return &AuthUseCase{repo}
}

func (s *AuthUseCase) SignupEmail(ctx context.Context, req *share.QEmailSignup) (*share.ZEmailSignup, error) {
	if !equal_passwords(req.Password, req.ConfirmedPwd) {
		return nil, fmt.Errorf("пароли не совпадают")
	}
	if !ValidateCredentials(req.Email, req.Password) {
		return nil, fmt.Errorf("не верные данные")
	}

	code, _ := CreateConfirmCode()
	passwd_hash, salt, _ := CreatePasswordHash(req.Password, "")

	xres, err := s.repo.CreateEmailSignup(ctx, req.Email, passwd_hash, code, salt)

	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	fmt.Println(xres, err)
	return &share.ZEmailSignup{
		ID:           xres.ID,
		Email:        xres.Email,
		Code:         xres.Code,
		PasswordHash: xres.PasswordHash,
		Salt:         xres.Salt,
		CreatedAt:    xres.CreatedAt,
		UpdatedAt:    xres.UpdatedAt,
	}, nil
}

func (s *AuthUseCase) ConfirmEmail(ctx context.Context, req *share.QConfirmEmail) (*share.ZAccount, error) {
	signup_acc, err := s.repo.GetEmailSignup(ctx, req.SignupID)
	if err != nil {
		return nil, fmt.Errorf("аккаунт не найден")
	}

	if req.Code != signup_acc.Code {
		return nil, fmt.Errorf("не верный код")
	}

	return nil, nil
}

// ----------- Tools -----------

func equal_passwords(password string, confirm_pwd string) bool {
	return password == confirm_pwd
}

func ValidateCredentials(email, password string) bool {
	if len(password) < 6 {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	return true

}
