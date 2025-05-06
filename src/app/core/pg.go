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

// NewPgRepo создает новый экземпляр PgRepo с заданной строкой подключения.
//
// Параметры:
//   - dsn: строка подключения к базе данных
//
// Возвращает:
//   - указатель на новый экземпляр PgRepo
func NewPgRepo(dsn string) *PgRepo {
	return &PgRepo{dsn: dsn}
}

// InitPool инициализирует пул соединений к базе данных PostgreSQL.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем жизни операции
//
// Возвращает:
//   - ошибку, если инициализация пула завершилась неудачей
func (r *PgRepo) InitPool(ctx context.Context) error {
	config, err := pgxpool.ParseConfig(r.dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DSN: %w", err)
	}

	config.MinConns = 1
	config.MaxConns = 2
	config.MaxConnLifetime = 500 * time.Second

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

// Acquire получает подключение из пула соединений к базе данных.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем ожидания получения соединения
//
// Возвращает:
//   - указатель на соединение из пула
//   - ошибку, если не удалось получить соединение
func (r *PgRepo) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	if r.pool == nil {
		if err := r.InitPool(ctx); err != nil {
			return nil, err
		}
	}
	return r.pool.Acquire(ctx)
}
