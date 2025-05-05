package repo

import (
	"context"
)

type AuthRepo interface {
	CreateEmailSignup(ctx context.Context, user *XEmailSignup) error
	GetEmailSignup(ctx context.Context, email string) (*XEmailSignup, error)
	DeleteEmailSignup(ctx context.Context, id string) error
	GetRefreshTokenForAccount(ctx context.Context, account_id string) (XRefreshToken, error)
	SaveRefreshToken(ctx context.Context, account_id string, token string)
}
