package auth

import (
	"time"
)

type QEmailSignup struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"123123"`
}

type QConfirmEmail struct {
	SignupID string `json:"signup_id"`
	Password string `json:"password"`
}

type ZEmailSignup struct {
	ID        string    `json:"id" example:"592af5b5-4f60-4ddd-b080-be674c86eda8"`
	CreatedAt time.Time `json:"created_at" example:"2024-02-13 05:37:40.483836"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-02-13 05:37:40.483836"`
}

type ZToken struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"bearer"`
}

type ZAccountID struct {
	ID string `json:"id" example:"592af5b5-4f60-4ddd-b080-be674c86eda8"`
}

type ZAccount struct {
	ID        string    `json:"id" example:"592af5b5-4f60-4ddd-b080-be674c86eda8"`
	Email     string    `json:"email" example:"user@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2024-02-13 05:37:40.483836"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-02-13 05:37:40.483836"`
}

type QRefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}
