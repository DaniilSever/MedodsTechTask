package core

import "fmt"

// ------------- for swagger -------------
type ZError struct {
	Code      int    `json:"code" example:"500"`
	Where     string `json:"where" example:"ExampleAPI"`
	Message   string `json:"message" example:"Пример ошибки"`
	Exception any    `json:"exception"`
}

// ------------- for UseCase -------------

type ErrInvalidLenPassword struct {
	ErrMessage any
}

type ErrEmailValidate struct {
	ErrMessage any
}

// ------------- Error Func to uc -------------

func (e *ErrInvalidLenPassword) Error() string {
	return fmt.Sprintf("неверная длина пароля \nerr: %s", e.ErrMessage)
}

func (e *ErrEmailValidate) Error() string {
	return fmt.Sprintf("не валидная почта \nerr: %s", e.ErrMessage)
}

// ------------- for repo -------------

type ErrPGRepo struct {
	ErrMessage any
}

type ErrCreateSignup struct {
	ErrMessage any
}

type ErrCreateAccount struct {
	ErrMessage any
}

type ErrAccountNotFound struct {
	ErrMessage any
}

type ErrEmailSignupNotFound struct {
	ErrMessage any
}

type ErrSaveToken struct {
	ErrMessage any
}

type ErrTokenNotFound struct {
	ErrMessage any
}

// ------------- Error Func to repo -------------

func (e *ErrPGRepo) Error() string {
	return fmt.Sprintf("произошла ошибка базы данных \nerr: %s", e.ErrMessage)
}

func (e *ErrCreateSignup) Error() string {
	return fmt.Sprintf("ошибка создания записи регистрации \nerr: %s", e.ErrMessage)
}

func (e *ErrCreateAccount) Error() string {
	return fmt.Sprintf("ошибка создания аккаунта для \nerr: %s", e.ErrMessage)
}

func (e *ErrAccountNotFound) Error() string {
	return fmt.Sprintf("ошибка подтвержденный аккаунт не найден \nerr: %s", e.ErrMessage)
}

func (e *ErrEmailSignupNotFound) Error() string {
	return fmt.Sprintf("ошибка зарегистрированный аккаунт не найден \nerr: %s", e.ErrMessage)
}

func (e *ErrSaveToken) Error() string {
	return fmt.Sprintf("не удалось сохранить токен \nerr: %s", e.ErrMessage)
}

func (e *ErrTokenNotFound) Error() string {
	return fmt.Sprintf("токен не найден \nerr: %s", e.ErrMessage)
}

// ------------- for security -------------

type ErrPasswordEmpty struct {
	ErrMessage any
}

type ErrGenerationSalt struct {
	ErrMessage any
}

type ErrGenerationHash struct {
	ErrMessage any
}

type ErrGenerationConfirmCode struct {
	ErrMessage any
}

type ErrParsePrivateKey struct {
	ErrMessage any
}

type ErrParsePublicKey struct {
	ErrMessage any
}

type ErrSignedJwt struct {
	ErrMessage any
}

type ErrUnExpectedSign struct {
	ErrMessage any
}

type ErrJwtExpired struct {
	ErrMessage any
}

type ErrIncorrectJwt struct {
	ErrMessage any
}

type ErrInvalidJwtPayload struct {
	ErrMessage any
}

// ------------- Error Func to security -------------

func (e *ErrPasswordEmpty) Error() string {
	return fmt.Sprintf("пароль пустой \nerr: %s", e.ErrMessage)
}

func (e *ErrGenerationSalt) Error() string {
	return fmt.Sprintf("ошибка генерации соли для пароля \nerr: %s", e.ErrMessage)
}

func (e *ErrGenerationHash) Error() string {
	return fmt.Sprintf("ошибка генерации хеша пароля \nerr: %s", e.ErrMessage)
}

func (e *ErrGenerationConfirmCode) Error() string {
	return fmt.Sprintf("ошибка генерации кода подтверждения аккаунта \nerr: %s", e.ErrMessage)
}

func (e *ErrParsePrivateKey) Error() string {
	return fmt.Sprintf("ошибка парсинга приватного ключа \nerr: %s", e.ErrMessage)
}

func (e *ErrParsePublicKey) Error() string {
	return fmt.Sprintf("ошибка парсинга публичного ключа \nerr: %s", e.ErrMessage)
}

func (e *ErrSignedJwt) Error() string {
	return fmt.Sprintf("ошибка подписи jwt ключа \nerr: %s", e.ErrMessage)
}

func (e *ErrUnExpectedSign) Error() string {
	return fmt.Sprintf("не правильный метод подписания jwt \nerr: %s", e.ErrMessage)
}

func (e *ErrJwtExpired) Error() string {
	return fmt.Sprintf("jwt ключ истек \nerr: %s", e.ErrMessage)
}

func (e *ErrIncorrectJwt) Error() string {
	return fmt.Sprintf("не верный jwt ключ \nerr: %s", e.ErrMessage)
}

func (e *ErrInvalidJwtPayload) Error() string {
	return fmt.Sprintf("недействительная полезная нагрузка JWT \nerr: %s", e.ErrMessage)
}
