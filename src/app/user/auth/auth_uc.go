package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/MedodsTechTask/app/user/auth/configs"
	"github.com/MedodsTechTask/app/user/auth/repo"
	"github.com/MedodsTechTask/app/user/auth/share"
)

type AuthUseCase struct {
	cfg  *configs.Config
	repo repo.IAuthRepo
}

func NewAuthUseCase(cfg *configs.Config, repo repo.IAuthRepo) *AuthUseCase {
	return &AuthUseCase{cfg, repo}
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

	del, err := s.repo.DeleteEmailSignup(ctx, req.SignupID)
	if !del && err != nil {
		return nil, err
	}

	xres, err := s.repo.CreateAccount(ctx, signup_acc)
	if err != nil {
		return nil, err
	}

	return &share.ZAccount{
		ID:         xres.ID,
		Email:      xres.Email,
		PasswdHash: xres.PasswordHash,
		Salt:       xres.Salt,
		CreatedAt:  xres.CreatedAt,
		UpdatedAt:  xres.UpdatedAt,
	}, nil
}

func (s *AuthUseCase) LoginEmail(ctx context.Context, login *share.QLoginEmail) (*share.ZToken, error) {
	acc, err := s.repo.GetAccountForEmail(ctx, login.Email)
	if err != nil {
		return nil, err
	}

	pwd_hash, _, err := CreatePasswordHash(login.Password, acc.Salt)
	if err != nil {
		return nil, err
	}
	if pwd_hash != acc.PasswordHash {
		return nil, fmt.Errorf("не верный пароль")
	}

	access_payload := map[string]interface{}{
		"sub":  acc.ID,
		"type": "access",
	}

	token, err := CreateJWT(access_payload, s.cfg.JWTPrivateKey)
	if err != nil {
		return nil, err
	}

	refresh_payload := map[string]interface{}{
		"sub":  acc.ID,
		"type": "refresh",
	}

	refresh_token, err := CreateJWT(refresh_payload, s.cfg.JWTPrivateKey)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.SaveRefreshToken(ctx, acc.ID, refresh_token)
	if err != nil {
		return nil, err
	}

	return &share.ZToken{
		AccessToken:  token,
		RefreshToken: refresh_token,
		TokenType:    "bearer",
	}, nil
}

func (s *AuthUseCase) RefreshToken(ctx context.Context, req *share.QRefreshToken) (*share.ZToken, error) {
	payload, err := DecodeJWT(req.RefreshToken, s.cfg.JWTPublicKey)
	if err != nil {
		return nil, err
	}
	if payload["type"] != "refresh" {
		return nil, fmt.Errorf("не верный токен")
	}

	acc_id := payload["sub"].(string)

	_, err = s.repo.GetRefreshTokenForAccount(ctx, acc_id, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	access_payload := map[string]interface{}{
		"sub":  acc_id,
		"type": "access",
	}

	token, err := CreateJWT(access_payload, s.cfg.JWTPrivateKey)
	if err != nil {
		return nil, err
	}

	return &share.ZToken{
		AccessToken:  token,
		RefreshToken: req.RefreshToken,
		TokenType:    "bearer",
	}, nil
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
