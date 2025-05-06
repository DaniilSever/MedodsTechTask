package core

import "fmt"

type ErrPGRepo struct {
	ErrMessage any
}

type ErrCreateSignup struct {
	Email      string
	ErrMessage any
}

type ErrCreateAccount struct {
	Email      string
	ErrMessage any
}

type ErrAccountNotFound struct {
	Email      string
	ErrMessage any
}

type ErrEmailSignupNotFound struct {
	ID         string
	ErrMessage any
}

type ErrSaveToken struct {
	Token      string
	ErrMessage any
}

type ErrTokenNotFound struct {
	Token      string
	ErrMessage any
}

// ----------------------------------

func (e *ErrPGRepo) Error() string {
	return fmt.Sprintf("произошла ошибка базы данных: %s", e.ErrMessage)
}

func (e *ErrCreateSignup) Error() string {
	return fmt.Sprintf("ошибка создания записи регистрации для %s, сообщение: %s", e.Email, e.ErrMessage)
}

func (e *ErrCreateAccount) Error() string {
	return fmt.Sprintf("ошибка создания аккаунта для %s, сообщение: %s", e.Email, e.ErrMessage)
}

func (e *ErrAccountNotFound) Error() string {
	return fmt.Sprintf("ошибка аккаунт %s не найден, сообщение: %s", e.Email, e.ErrMessage)
}

func (e *ErrEmailSignupNotFound) Error() string {
	return fmt.Sprintf("ошибка аккаунт %s не найден, сообщение: %s", e.ID, e.ErrMessage)
}

func (e *ErrSaveToken) Error() string {
	return fmt.Sprintf("не удалось сохранить токен %s, сообщение: %s", e.Token, e.ErrMessage)
}

func (e *ErrTokenNotFound) Error() string {
	return fmt.Sprintf("токен %s не найден, сообщение: %s", e.Token, e.ErrMessage)
}
