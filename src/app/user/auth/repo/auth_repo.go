package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/MedodsTechTask/app/core"
	"github.com/MedodsTechTask/app/user/auth/configs"
	"github.com/MedodsTechTask/app/user/auth/share"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type IAuthRepo interface {
	CreateEmailSignup(ctx context.Context, email string, passwd_hash string, code string, salt string) (*XEmailSignup, error)
	CreateAccount(ctx context.Context, req *share.QCreateAccount) (*XAccount, error)
	GetEmailSignup(ctx context.Context, id string) (*XEmailSignup, error)
	GetAccountForEmail(ctx context.Context, email string) (*XAccount, error)
	DeleteEmailSignup(ctx context.Context, id string) (bool, error)
	SaveRefreshToken(ctx context.Context, account_id string, token string) (*XRefreshToken, error)
	GetRefreshTokenForAccount(ctx context.Context, account_id string, token string) (*XRefreshToken, error)
}

type AuthRepo struct {
	pgRepo *core.PgRepo
	cfg    *configs.Config
}

func NewAuthRepo(cfg *configs.Config) (*AuthRepo, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.AuthDBUSR,
		cfg.AuthDBPWD,
		cfg.DBMasterHost,
		cfg.DBMasterPort,
		cfg.AuthDBDBN,
	)
	pgRepo := core.NewPgRepo(dsn)
	if err := pgRepo.InitPool(context.Background()); err != nil {
		return nil, fmt.Errorf("filed to init DB pool %w", err)
	}

	return &AuthRepo{
		pgRepo: pgRepo,
		cfg:    cfg,
	}, nil
}

func (r *AuthRepo) CreateEmailSignup(ctx context.Context, email string, passwd_hash string, code string, salt string) (*XEmailSignup, error) {
	const q = `
		INSERT INTO "SignupEmail"
		(
			email
			, code
			, passwd_hash
			, salt
		)
		VALUES ($1, $2, $3, $4)
		RETURNING *;
	`
	conn, err := r.pgRepo.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("error to connections on PG: %w", err)
	}
	defer conn.Release()

	var res XEmailSignup
	err = conn.QueryRow(ctx, q, email, code, passwd_hash, salt).Scan(&res.ID, &res.Email, &res.Code, &res.PasswordHash, &res.Salt, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, &core.ErrCreateSignup{Email: email, ErrMessage: err}
		}

		return nil, &core.ErrPGRepo{ErrMessage: err}
	}

	return &res, nil
}

func (r *AuthRepo) CreateAccount(ctx context.Context, req *share.QCreateAccount) (*XAccount, error) {
	const q = `
		INSERT INTO "Account"
		(
			email
			, passwd_hash
			, salt
		)
		VALUES ($1, $2, $3)
		RETURNING *;
	`

	conn, err := r.pgRepo.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("error to connections on PG: %w", err)
	}
	defer conn.Release()

	var res XAccount
	err = conn.QueryRow(ctx, q, req.Email, req.PasswdHash, req.Salt).Scan(&res.ID, &res.Email, &res.PasswordHash, &res.Salt, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, &core.ErrCreateAccount{Email: req.Email, ErrMessage: err}
		}

		return nil, &core.ErrPGRepo{ErrMessage: err}
	}
	return &res, nil
}

func (r *AuthRepo) GetEmailSignup(ctx context.Context, id string) (*XEmailSignup, error) {
	const q = `
		SELECT
			id
			, email
			, code
			, passwd_hash
			, salt
			, created_at
			, updated_at
		FROM "SignupEmail"
		WHERE id = $1
		LIMIT 1;
	`

	conn, err := r.pgRepo.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("error to connections on PG: %w", err)
	}
	defer conn.Release()

	var res XEmailSignup
	err = conn.QueryRow(ctx, q, id).Scan(&res.ID, &res.Email, &res.Code, &res.PasswordHash, &res.Salt, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &core.ErrEmailSignupNotFound{ID: id, ErrMessage: err}
		}

		return nil, &core.ErrPGRepo{ErrMessage: err}
	}

	return &res, nil
}

func (r *AuthRepo) GetAccountForEmail(ctx context.Context, email string) (*XAccount, error) {
	const q = `
		SELECT
			id
			, email
			, code
			, passwd_hash
			, salt
			, created_at
			, updated_at
		FROM "Account"
		WHERE email = $1
		LIMIT 1;
	`

	conn, err := r.pgRepo.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("error to connections on PG: %w", err)
	}
	defer conn.Release()

	var res XAccount
	err = conn.QueryRow(ctx, q, email).Scan(&res.ID, &res.Email, &res.PasswordHash, &res.Salt, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &core.ErrAccountNotFound{Email: email, ErrMessage: err}
		}

		return nil, &core.ErrPGRepo{ErrMessage: err}
	}
	return &res, nil
}

func (r *AuthRepo) DeleteEmailSignup(ctx context.Context, id string) (bool, error) {
	const q = `
		DELETE FROM "SignupEmail"
		WHERE id = $1;
	`

	conn, err := r.pgRepo.Acquire(ctx)
	if err != nil {
		return false, fmt.Errorf("error to connections on PG: %w", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, q, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, &core.ErrEmailSignupNotFound{ID: id, ErrMessage: err}
		}

		return false, &core.ErrPGRepo{ErrMessage: err}
	}

	return true, nil
}

// ------- tokens -------

func (r *AuthRepo) SaveRefreshToken(ctx context.Context, account_id string, token string) (*XRefreshToken, error) {
	const q = `
		INSERT INTO "RefreshToken"
		(
			account_id
			, token
			, expires_at
		)
		VALUES ($1, $2, NOW() + INTERVAL '5 days')
		RETURNING *
	`
	conn, err := r.pgRepo.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("error to connections on PG: %w", err)
	}
	defer conn.Release()

	var res XRefreshToken
	err = conn.QueryRow(ctx, q, account_id, token).Scan(&res.ID, &res.AccountID, &res.Token, &res.ExpiresAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, &core.ErrSaveToken{Token: token, ErrMessage: err}
		}

		return nil, &core.ErrPGRepo{ErrMessage: err}
	}

	return &res, nil
}

func (r *AuthRepo) GetRefreshTokenForAccount(ctx context.Context, account_id string, token string) (*XRefreshToken, error) {
	const q = `
		SELECT 
			id
			, account_id
			, token
			, expires_at
		FROM "RefreshToken"
		WHERE True
			AND account_id = $1
			AND token = $2
			AND expires_at > NOW()
	`

	conn, err := r.pgRepo.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("error to connections on PG: %w", err)
	}
	defer conn.Release()

	var res XRefreshToken
	err = conn.QueryRow(ctx, q, account_id, token).Scan(&res.ID, &res.AccountID, &res.Token, &res.ExpiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &core.ErrTokenNotFound{Token: token, ErrMessage: err}
		}

		return nil, &core.ErrPGRepo{ErrMessage: err}
	}

	return &res, nil
}
