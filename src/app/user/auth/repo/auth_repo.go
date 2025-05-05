package repo

import (
	"context"
	"fmt"
	"sync"

	"github.com/MedodsTechTask/app/core"
	"github.com/MedodsTechTask/app/user/auth"
)

type IAuthRepo interface {
	CreateEmailSignup(ctx context.Context, user *XEmailSignup) error
	GetEmailSignup(ctx context.Context, email string) (*XEmailSignup, error)
	DeleteEmailSignup(ctx context.Context, id string) error
	GetRefreshTokenForAccount(ctx context.Context, account_id string) (XRefreshToken, error)
	SaveRefreshToken(ctx context.Context, account_id string, token string)
}

type AuthRepo struct {
	pgRepo *core.PgRepo
	cfg    *auth.Config
}

var (
	repo *AuthRepo
	once sync.Once
)

func NewAuthRepo(cfg *auth.Config) (*AuthRepo, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.AuthDBUSR,
		cfg.AuthDBPWD,
		cfg.DBMasterHost,
		cfg.DBMasterPort,
		cfg.AuthDBDBN,
	)
	pgRepo := core.NewPgRepo(dsn)
	if err := pgRepo.InitPool(context.Background()); err != nil {
		return nil, fmt.Errorf("Failed to init DB pool %w", err)
	}

	return &AuthRepo{
		pgRepo: pgRepo,
		cfg:    cfg,
	}, nil
}
