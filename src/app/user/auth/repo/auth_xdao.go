package repo

import (
	"time"
)

type XEmailSignup struct {
	ID           string     `db:"id"`
	Email        string     `db:"email"`
	Code         string     `db:"code"` // Добавлено для local debug
	PasswordHash string     `db:"passwd_hash"`
	Salt         string     `db:"salt"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

type XAccount struct {
	ID           string     `db:"id"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"passwd_hash"`
	Salt         string     `db:"salt"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

type XConfirmEmail struct {
	ID       string `db:"id"`
	Password string `db:"password"`
}

type XRefreshToken struct {
	ID        string     `db:"id"`
	AccountID string     `db:"accouint_id"`
	Token     string     `db:"token"`
	UserAgent string     `db:"user_agent"`
	IpAddress string     `db:"ip_address"`
	ExpiresAt time.Time  `db:"expires_at"`
	IsRevoked bool       `db:"is_revoked"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
