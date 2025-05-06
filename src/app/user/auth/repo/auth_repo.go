package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/MedodsTechTask/app/core"
	"github.com/MedodsTechTask/app/user/auth/configs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type IAuthRepo interface {
	CreateEmailSignup(ctx context.Context, email string, passwd_hash string, code string, salt string) (*XEmailSignup, error)
	CreateAccount(ctx context.Context, req *XEmailSignup) (*XAccount, error)
	GetEmailSignup(ctx context.Context, id string) (*XEmailSignup, error)
	GetAccountForEmail(ctx context.Context, email string) (*XAccount, error)
	DeleteEmailSignup(ctx context.Context, id string) (bool, error)
	SaveRefreshToken(ctx context.Context, account_id string, user_agent string, ip_address string, token string) (*XRefreshToken, error)
	GetRefreshTokenForAccount(ctx context.Context, account_id string, token string) (*XRefreshToken, error)
	RevokeToken(ctx context.Context, account_id string) (bool, error)
}

type AuthRepo struct {
	pgRepo *core.PgRepo
	cfg    *configs.Config
}

// NewAuthRepo создает новый репозиторий аутентификации с инициализацией пула соединений с базой данных.
//
// Параметры:
//   - cfg: конфигурация приложения, содержащая параметры для подключения к базе данных
//
// Возвращает:
//   - указатель на новый экземпляр AuthRepo
//   - ошибку, если не удалось инициализировать репозиторий или пул соединений
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

// CreateEmailSignup создает новую запись о регистрации с email в базе данных.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем жизни операции
//   - email: email пользователя
//   - passwd_hash: хеш пароля пользователя
//   - code: код подтверждения
//   - salt: соль для хеширования пароля
//
// Возвращает:
//   - указатель на структуру XEmailSignup, содержащую информацию о регистрации
//   - ошибку, если операция не удалась (например, ошибка базы данных или нарушение уникальности)
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
		return nil, &core.ErrPGRepo{ErrMessage: err}
	}
	defer conn.Release()

	var res XEmailSignup
	err = conn.QueryRow(ctx, q, email, code, passwd_hash, salt).Scan(&res.ID, &res.Email, &res.Code, &res.PasswordHash, &res.Salt, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, &core.ErrCreateSignup{ErrMessage: err}
		}

		return nil, &core.ErrPGRepo{ErrMessage: err}
	}

	return &res, nil
}

// CreateAccount создает новый аккаунт в базе данных.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем жизни операции
//   - req: структура XEmailSignup, содержащая данные для создания аккаунта (email, хеш пароля, соль)
//
// Возвращает:
//   - указатель на структуру XAccount, содержащую информацию о созданном аккаунте
//   - ошибку, если операция не удалась (например, ошибка базы данных или нарушение уникальности)
func (r *AuthRepo) CreateAccount(ctx context.Context, req *XEmailSignup) (*XAccount, error) {
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
		return nil, &core.ErrPGRepo{ErrMessage: err}
	}
	defer conn.Release()

	var res XAccount
	err = conn.QueryRow(ctx, q, req.Email, req.PasswordHash, req.Salt).Scan(&res.ID, &res.Email, &res.PasswordHash, &res.Salt, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, &core.ErrCreateAccount{ErrMessage: err}
		}

		return nil, &core.ErrPGRepo{ErrMessage: err}
	}
	return &res, nil
}

// GetEmailSignup извлекает запись о регистрации с email из базы данных по заданному идентификатору.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем жизни операции
//   - id: идентификатор записи о регистрации
//
// Возвращает:
//   - указатель на структуру XEmailSignup с данными о регистрации
//   - ошибку, если запись не найдена или произошла ошибка при запросе к базе данных
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
		return nil, &core.ErrPGRepo{ErrMessage: err}
	}
	defer conn.Release()

	var res XEmailSignup
	err = conn.QueryRow(ctx, q, id).Scan(&res.ID, &res.Email, &res.Code, &res.PasswordHash, &res.Salt, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &core.ErrEmailSignupNotFound{ErrMessage: err}
		}

		return nil, &core.ErrPGRepo{ErrMessage: err}
	}

	return &res, nil
}

// GetAccountForEmail извлекает аккаунт из базы данных по заданному email.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем жизни операции
//   - email: email для поиска соответствующего аккаунта
//
// Возвращает:
//   - указатель на структуру XAccount с данными аккаунта
//   - ошибку, если аккаунт не найден или произошла ошибка при запросе к базе данных
func (r *AuthRepo) GetAccountForEmail(ctx context.Context, email string) (*XAccount, error) {
	const q = `
		SELECT
			id
			, email
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
		return nil, &core.ErrPGRepo{ErrMessage: err}
	}
	defer conn.Release()

	var res XAccount
	err = conn.QueryRow(ctx, q, email).Scan(&res.ID, &res.Email, &res.PasswordHash, &res.Salt, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &core.ErrAccountNotFound{ErrMessage: err}
		}

		return nil, &core.ErrPGRepo{ErrMessage: err}
	}
	return &res, nil
}

// DeleteEmailSignup удаляет запись о регистрации с email из базы данных по заданному идентификатору.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем жизни операции
//   - id: идентификатор записи о регистрации для удаления
//
// Возвращает:
//   - булевое значение, указывающее на успешность операции (true, если запись удалена)
//   - ошибку, если операция не удалась или запись не найдена
func (r *AuthRepo) DeleteEmailSignup(ctx context.Context, id string) (bool, error) {
	const q = `
		DELETE FROM "SignupEmail"
		WHERE id = $1;
	`

	conn, err := r.pgRepo.Acquire(ctx)
	if err != nil {
		return false, &core.ErrPGRepo{ErrMessage: err}
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, q, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, &core.ErrEmailSignupNotFound{ErrMessage: err}
		}

		return false, &core.ErrPGRepo{ErrMessage: err}
	}

	return true, nil
}

// SaveRefreshToken сохраняет новый refresh-токен в базе данных для заданного аккаунта.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем жизни операции
//   - account_id: идентификатор аккаунта, к которому привязан токен
//   - user_agent: строка, представляющая user-agent устройства пользователя
//   - ip_address: IP-адрес пользователя
//   - token: сам refresh-токен для сохранения
//
// Возвращает:
//   - указатель на структуру XRefreshToken с данными сохраненного токена
//   - ошибку, если операция не удалась (например, ошибка базы данных или нарушение уникальности)
func (r *AuthRepo) SaveRefreshToken(ctx context.Context, account_id string, user_agent string, ip_address string, token string) (*XRefreshToken, error) {
	const q = `
		INSERT INTO "RefreshToken"
		(
			account_id
			, token
			, user_agent
			, ip_address
			, expires_at
		)
		VALUES ($1, $2, $3, $4, NOW() + INTERVAL '5 days')
		RETURNING *;
	`

	conn, err := r.pgRepo.Acquire(ctx)
	if err != nil {
		return nil, &core.ErrPGRepo{ErrMessage: err}
	}
	defer conn.Release()

	var res XRefreshToken
	err = conn.QueryRow(ctx, q, account_id, token, user_agent, ip_address).Scan(&res.ID, &res.AccountID, &res.Token, &res.UserAgent, &res.IpAddress, &res.ExpiresAt, &res.IsRevoked, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, &core.ErrSaveToken{ErrMessage: err}
		}

		return nil, &core.ErrPGRepo{ErrMessage: err}
	}

	return &res, nil
}

// GetRefreshTokenForAccount извлекает refresh-токен для заданного аккаунта из базы данных.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем жизни операции
//   - account_id: идентификатор аккаунта, для которого нужно получить токен
//   - token: сам refresh-токен для поиска в базе данных
//
// Возвращает:
//   - указатель на структуру XRefreshToken с данными найденного токена
//   - ошибку, если токен не найден или произошла ошибка при запросе к базе данных
func (r *AuthRepo) GetRefreshTokenForAccount(ctx context.Context, account_id string, token string) (*XRefreshToken, error) {
	const q = `
		SELECT 
			id
			, account_id
			, token
			, user_agent
			, ip_address
			, expires_at
			, is_revoked
			, created_at
			, updated_at
		FROM "RefreshToken"
		WHERE True
			AND account_id = $1
			AND token = $2
			AND expires_at > NOW()
	`

	conn, err := r.pgRepo.Acquire(ctx)
	if err != nil {
		return nil, &core.ErrPGRepo{ErrMessage: err}
	}
	defer conn.Release()

	var res XRefreshToken
	err = conn.QueryRow(ctx, q, account_id, token).Scan(&res.ID, &res.AccountID, &res.Token, &res.UserAgent, &res.IpAddress, &res.ExpiresAt, &res.IsRevoked, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &core.ErrTokenNotFound{ErrMessage: err}
		}

		return nil, &core.ErrPGRepo{ErrMessage: err}
	}

	return &res, nil
}

// RevokeToken обновляет статус refresh-токенов для заданного аккаунта, помечая их как отозванные.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем жизни операции
//   - account_id: идентификатор аккаунта, для которого нужно отозвать токены
//
// Возвращает:
//   - булевое значение, указывающее на успешность операции (true, если токены были отозваны)
//   - ошибку, если операция не удалась
func (r *AuthRepo) RevokeToken(ctx context.Context, account_id string) (bool, error) {
	const q = `
		UPDATE "RefreshToken"
		SET is_revoked = TRUE,
		updated_at = NOW()
		WHERE account_id = $1
	`

	conn, err := r.pgRepo.Acquire(ctx)
	if err != nil {
		return false, &core.ErrPGRepo{ErrMessage: err}
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, q, account_id)
	if err != nil {
		return false, &core.ErrPGRepo{ErrMessage: err}
	}

	return true, nil
}
