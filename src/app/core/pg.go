package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAuditInsertFailed = errors.New("audit insert failed")
)

type PgRepo struct {
	pool *pgxpool.Pool
	dsn  string
}

// NewPgRepo создает новый репозиторий с подключением к PostgreSQL
func NewPgRepo(dsn string) *PgRepo {
	return &PgRepo{dsn: dsn}
}

// InitPool инициализирует пул подключений
func (r *PgRepo) InitPool(ctx context.Context) error {
	config, err := pgxpool.ParseConfig(r.dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	config.MinConns = 1
	config.MaxConns = 2
	config.MaxConnLifetime = 500 * time.Second

	// Настройка кодеков для JSON и UUID
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterDefaultPgType(json.RawMessage{}, "jsonb")
		return nil
	}

	r.pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	return nil
}

// Acquire возвращает соединение из пула
func (r *PgRepo) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	if r.pool == nil {
		if err := r.InitPool(ctx); err != nil {
			return nil, err
		}
	}
	return r.pool.Acquire(ctx)
}

// HealthCheck проверяет доступность БД
func (r *PgRepo) HealthCheck(ctx context.Context) (bool, error) {
	conn, err := r.Acquire(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, "SELECT 1"); err != nil {
		return false, nil
	}

	return true, nil
}
