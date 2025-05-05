package core

const (
	// BasePath - общий роут api
	BasePath = "/api/v1"
)

const (
	// UserAuthPath - роут auth сервиса
	UserAuthPath = "/user/auth"

	// UserAuthHealthCheck - Проверка состояния сервиса auth
	UserAuthHealthCheck = "/health"

	// UserAuthSignUpEmail - Регистрация пользователя через email + пароль
	UserAuthSignUpEmail = "/signup/email"

	// UserAuthConfirmEmail - Подтверждение емейла вводом кода или по ссылке
	UserAuthConfirmEmail = "/confirm/email"

	// UserAuthLoginEmail - Вход через емейл+пароль
	UserAuthLoginEmail = "/login/email"

	// UserAuthRefreshToken - Рефреш токена
	UserAuthRefreshToken = "/refresh/token"
)
