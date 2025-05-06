package share

import (
	"time"
)

type QEmailSignup struct {
	Email        string `json:"email" example:"user@example.com"`
	Password     string `json:"password" example:"123123"`
	ConfirmedPwd string `json:"confim_pwd" example:"123123"`
}

type QConfirmEmail struct {
	SignupID string `json:"signup_id" example:"592af5b5-4f60-4ddd-b080-be674c86eda8"`
	Code     string `json:"code" example:"123456"`
}

type QCreateEmailSignup struct {
	Email      string     `json:"email"`
	PasswdHash string     `json:"passwd_hash"`
	Salt       string     `json:"salt"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type QCreateAccount struct {
	Email      string     `json:"email"`
	PasswdHash string     `json:"passwd_hash"`
	Salt       string     `json:"salt"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type ZEmailSignup struct {
	ID           string     `json:"id" example:"592af5b5-4f60-4ddd-b080-be674c86eda8"`
	Email        string     `json:"email" example:"user@example.com"`
	Code         string     `json:"code" example:"123456"` // Добавлено для local debug
	PasswordHash string     `json:"passwd_hash" example:"592af5b54f604dddb080be674c86eda8"`
	Salt         string     `json:"salt"  example:"592af5b54f604dddb080be674c86eda8"`
	CreatedAt    time.Time  `json:"created_at" example:"2024-02-13 05:37:40.483836"`
	UpdatedAt    *time.Time `json:"updated_at" example:"2024-02-13 05:37:40.483836"`
}

type ZToken struct {
	AccessToken  string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NSIsImV4cCI6MTk5OTk5OTk5OSwiaWF0IjoxNzAwMDAwMDAwLCJpc3MiOiJleGFtcGxlLWFwcCJ9.rnH9fqOBlB4tfbgIhJX_yta9Z9yVtOmMFLhy5aC_cC8"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NSIsImV4cCI6MTk5OTk5OTk5OSwiaWF0IjoxNzAwMDAwMDAwLCJpc3MiOiJleGFtcGxlLWFwcCJ9.rnH9fqOBlB4tfbgIhJX_yta9Z9yVtOmMFLhy5aC_cC8"`
	TokenType    string `json:"bearer" example:"Authefication"`
}

type ZAccountID struct {
	ID string `json:"id" example:"592af5b5-4f60-4ddd-b080-be674c86eda8"`
}

type ZAccount struct {
	ID        string     `json:"id" example:"592af5b5-4f60-4ddd-b080-be674c86eda8"`
	Email     string     `json:"email" example:"user@example.com"`
	CreatedAt time.Time  `json:"created_at" example:"2024-02-13 05:37:40.483836"`
	UpdatedAt *time.Time `json:"updated_at" example:"2024-02-13 05:37:40.483836"`
}

type QRefreshToken struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NSIsImV4cCI6MTk5OTk5OTk5OSwiaWF0IjoxNzAwMDAwMDAwLCJpc3MiOiJleGFtcGxlLWFwcCJ9.rnH9fqOBlB4tfbgIhJX_yta9Z9yVtOmMFLhy5aC_cC8"`
}
