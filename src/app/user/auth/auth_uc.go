package auth

import (
	"context"
	"strings"

	"github.com/MedodsTechTask/app/core"
	"github.com/MedodsTechTask/app/user/auth/configs"
	"github.com/MedodsTechTask/app/user/auth/repo"
	"github.com/MedodsTechTask/app/user/auth/share"
)

type AuthUseCase struct {
	cfg  *configs.Config
	repo repo.IAuthRepo
}

// NewAuthUseCase создает новый экземпляр AuthUseCase с заданной конфигурацией и репозиторием.
//
// Параметры:
//   - cfg: конфигурация приложения, содержащая параметры для аутентификации
//   - repo: интерфейс репозитория для работы с данными аутентификации
//
// Возвращает:
//   - указатель на новый экземпляр AuthUseCase
func NewAuthUseCase(cfg *configs.Config, repo repo.IAuthRepo) *AuthUseCase {
	return &AuthUseCase{cfg, repo}
}

// SignupEmail обрабатывает процесс регистрации пользователя через email. Проверяет совпадение паролей, валидирует email и пароль,
// генерирует код подтверждения и хеширует пароль, после чего сохраняет данные в базе данных.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем жизни операции
//   - req: структура с данными для регистрации пользователя (email, пароль, подтвержденный пароль)
//
// Возвращает:
//   - указатель на структуру ZEmailSignup с данными пользователя, если регистрация прошла успешно
//   - указатель на структуру ZError с описанием ошибки, если произошла ошибка на любом этапе
func (s *AuthUseCase) SignupEmail(ctx context.Context, req *share.QEmailSignup) (*share.ZEmailSignup, *core.ZError) {
	if !equal_passwords(req.Password, req.ConfirmedPwd) {
		return nil, &core.ZError{
			Code:      400,
			Where:     "UseCase",
			Message:   "Пароли не совпадают",
			Exception: nil,
		}
	}

	res, err := ValidateCredentials(req.Email, req.Password)
	if !res && err != nil {
		switch e := err.(type) {
		case *core.ErrInvalidLenPassword:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase",
				Message:   "Неверная длина пароля (мин. 6 символов)",
				Exception: e.ErrMessage,
			}
		case *core.ErrEmailValidate:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase",
				Message:   "Неверная почта",
				Exception: e.ErrMessage,
			}
		}
	}

	code, err := CreateConfirmCode()
	if err != nil {
		return nil, &core.ZError{
			Code:      400,
			Where:     "UseCase",
			Message:   "Ошибка генерации кода подтверждения",
			Exception: nil,
		}
	}
	passwd_hash, salt, err := CreatePasswordHash(req.Password, "")
	if err != nil {
		switch e := err.(type) {
		case *core.ErrPasswordEmpty:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase",
				Message:   "Пароль пустой",
				Exception: nil,
			}
		case *core.ErrGenerationSalt:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase",
				Message:   "Ошибка генерации соли пароля",
				Exception: e.ErrMessage,
			}
		case *core.ErrGenerationHash:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase",
				Message:   "Ошибка генерации хеша пароля",
				Exception: e.ErrMessage,
			}
		}
	}

	xres, err := s.repo.CreateEmailSignup(ctx, req.Email, passwd_hash, code, salt)
	if err != nil {
		switch e := err.(type) {
		case *core.ErrCreateSignup:
			return nil, &core.ZError{
				Code:      400,
				Where:     "Repo",
				Message:   "Ошибка регистрации пользователя",
				Exception: e.ErrMessage,
			}
		case *core.ErrPGRepo:
			return nil, &core.ZError{
				Code:      500,
				Where:     "Repo",
				Message:   "Неизвестная ошибка базы данных",
				Exception: e.ErrMessage,
			}
		}
	}

	return &share.ZEmailSignup{
		ID:           xres.ID,
		Email:        xres.Email,
		Code:         xres.Code,
		PasswordHash: xres.PasswordHash,
		Salt:         xres.Salt,
		CreatedAt:    xres.CreatedAt,
		UpdatedAt:    xres.UpdatedAt,
	}, nil
}

// ConfirmEmail подтверждает регистрацию пользователя по коду, отправленному на email.
// В случае успешного подтверждения, создает аккаунт пользователя и удаляет запись о регистрации.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем жизни операции
//   - req: структура с данными для подтверждения регистрации (ID регистрации и код подтверждения)
//
// Возвращает:
//   - указатель на структуру ZAccount с данными созданного аккаунта, если подтверждение прошло успешно
//   - указатель на структуру ZError с описанием ошибки, если произошла ошибка на любом этапе
func (s *AuthUseCase) ConfirmEmail(ctx context.Context, req *share.QConfirmEmail) (*share.ZAccount, *core.ZError) {
	signup_acc, err := s.repo.GetEmailSignup(ctx, req.SignupID)
	if err != nil {
		switch e := err.(type) {
		case *core.ErrEmailSignupNotFound:
			return nil, &core.ZError{
				Code:      404,
				Where:     "Repo",
				Message:   "Аккаунт не найден",
				Exception: e.ErrMessage,
			}
		case *core.ErrPGRepo:
			return nil, &core.ZError{
				Code:      500,
				Where:     "Repo",
				Message:   "Неизвестная ошибка базы данных",
				Exception: e.ErrMessage,
			}
		}
	}

	if req.Code != signup_acc.Code {
		return nil, &core.ZError{
			Code:      400,
			Where:     "UseCase",
			Message:   "Неверный код подтвержден",
			Exception: nil,
		}
	}

	del, err := s.repo.DeleteEmailSignup(ctx, req.SignupID)
	if !del && err != nil {
		switch e := err.(type) {
		case *core.ErrEmailSignupNotFound:
			return nil, &core.ZError{
				Code:      404,
				Where:     "Repo",
				Message:   "Аккаунт не найден",
				Exception: e.ErrMessage,
			}
		case *core.ErrPGRepo:
			return nil, &core.ZError{
				Code:      500,
				Where:     "Repo",
				Message:   "Неизвестная ошибка базы данных",
				Exception: e.ErrMessage,
			}
		}
	}

	xres, err := s.repo.CreateAccount(ctx, signup_acc)
	if err != nil {
		switch e := err.(type) {
		case *core.ErrCreateAccount:
			return nil, &core.ZError{
				Code:      400,
				Where:     "Repo",
				Message:   "Не удалось создать аккаунт",
				Exception: e.ErrMessage,
			}
		case *core.ErrPGRepo:
			return nil, &core.ZError{
				Code:      500,
				Where:     "Repo",
				Message:   "Неизвестная ошибка базы данных",
				Exception: e.ErrMessage,
			}
		}
	}

	return &share.ZAccount{
		ID:         xres.ID,
		Email:      xres.Email,
		PasswdHash: xres.PasswordHash,
		Salt:       xres.Salt,
		CreatedAt:  xres.CreatedAt,
		UpdatedAt:  xres.UpdatedAt,
	}, nil
}

// LoginEmail обрабатывает процесс входа пользователя через email и пароль.
// Он проверяет существование аккаунта, валидирует пароль и генерирует токены доступа и обновления.
// Также сохраняет refresh токен в базе данных.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем жизни операции
//   - login: структура с email и паролем для авторизации
//   - user_agent: строка с информацией о пользовательском агенте
//   - ip: строка с IP-адресом пользователя
//
// Возвращает:
//   - указатель на структуру ZToken с access и refresh токенами, если авторизация прошла успешно
//   - указатель на структуру ZError с описанием ошибки, если произошла ошибка на любом этапе
func (s *AuthUseCase) LoginEmail(ctx context.Context, login *share.QLoginEmail, user_agent string, ip string) (*share.ZToken, *core.ZError) {
	acc, err := s.repo.GetAccountForEmail(ctx, login.Email)
	if err != nil {
		switch e := err.(type) {
		case *core.ErrAccountNotFound:
			return nil, &core.ZError{
				Code:      404,
				Where:     "Repo",
				Message:   "Аккаунт не найден",
				Exception: e.ErrMessage,
			}
		case *core.ErrPGRepo:
			return nil, &core.ZError{
				Code:      500,
				Where:     "Repo",
				Message:   "Неизвестная ошибка базы данных",
				Exception: e.ErrMessage,
			}
		}
	}

	pwd_hash, _, err := CreatePasswordHash(login.Password, acc.Salt)
	if err != nil {
		switch e := err.(type) {
		case *core.ErrPasswordEmpty:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase",
				Message:   "Пароль пустой",
				Exception: e.ErrMessage,
			}
		case *core.ErrGenerationHash:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase",
				Message:   "Ошибка генерации хеша пароля",
				Exception: e.ErrMessage,
			}
		}
	}
	if pwd_hash != acc.PasswordHash {
		return nil, &core.ZError{
			Code:      400,
			Where:     "UseCase",
			Message:   "Пароли не совпадают",
			Exception: nil,
		}
	}

	access_payload := map[string]interface{}{
		"sub":  acc.ID,
		"type": "access",
	}

	token, err := CreateJWT(access_payload, s.cfg.JWTPrivateKey)
	if err != nil {
		switch e := err.(type) {
		case *core.ErrParsePrivateKey:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase/Security",
				Message:   "Ошибка парсинга приватного ключа",
				Exception: e.ErrMessage,
			}
		case *core.ErrSignedJwt:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase/Security",
				Message:   "Ошибка подписания jwt ключа",
				Exception: e.ErrMessage,
			}
		}
	}

	refresh_payload := map[string]interface{}{
		"sub":  acc.ID,
		"type": "refresh",
	}

	refresh_token, err := CreateJWT(refresh_payload, s.cfg.JWTPrivateKey)
	if err != nil {
		switch e := err.(type) {
		case *core.ErrParsePrivateKey:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCasUseCase/Securitye",
				Message:   "Ошибка парсинга приватного ключа",
				Exception: e.ErrMessage,
			}
		case *core.ErrSignedJwt:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase/Security",
				Message:   "Ошибка подписания jwt ключа",
				Exception: e.ErrMessage,
			}
		}
	}

	_, err = s.repo.SaveRefreshToken(ctx, acc.ID, user_agent, ip, refresh_token)
	if err != nil {
		switch e := err.(type) {
		case *core.ErrSaveToken:
			return nil, &core.ZError{
				Code:      400,
				Where:     "Repo",
				Message:   "Не удалось сохранить токен",
				Exception: e.ErrMessage,
			}
		case *core.ErrPGRepo:
			return nil, &core.ZError{
				Code:      500,
				Where:     "Repo",
				Message:   "Неизвестная ошибка базы данных",
				Exception: e.ErrMessage,
			}
		}
	}

	return &share.ZToken{
		AccessToken:  token,
		RefreshToken: refresh_token,
		TokenType:    "bearer",
	}, nil
}

// RefreshToken обрабатывает запрос на обновление токена доступа с использованием refresh токена.
// Он проверяет действительность refresh токена, его тип и соответствие с данными пользователя,
// а также проверяет, не был ли токен отозван. В случае успеха возвращает новый токен доступа.
//
// Параметры:
//   - ctx: контекст выполнения для управления временем жизни операции
//   - req: структура с refresh токеном для обновления
//   - user_agent: строка с информацией о пользовательском агенте
//   - ip: строка с IP-адресом пользователя
//
// Возвращает:
//   - указатель на структуру ZToken с новым access токеном и старым refresh токеном, если обновление прошло успешно
//   - указатель на структуру ZError с описанием ошибки, если произошла ошибка на любом этапе
func (s *AuthUseCase) RefreshToken(ctx context.Context, req *share.QRefreshToken, user_agent string, ip string) (*share.ZToken, *core.ZError) {
	payload, err := DecodeJWT(req.RefreshToken, s.cfg.JWTPublicKey)
	if err != nil {
		switch e := err.(type) {
		case *core.ErrParsePublicKey:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase/Security",
				Message:   "Ошибка парсинга публичного ключа",
				Exception: e.ErrMessage,
			}
		case *core.ErrUnExpectedSign:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase/Security",
				Message:   "Ошибка не верный метод подписания JWT ключа",
				Exception: e.ErrMessage,
			}
		case *core.ErrJwtExpired:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase/Security",
				Message:   "JWT ключ истек",
				Exception: e.ErrMessage,
			}
		case *core.ErrIncorrectJwt:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase/Security",
				Message:   "Неверный JWT ключ",
				Exception: e.ErrMessage,
			}
		case *core.ErrInvalidJwtPayload:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase/Security",
				Message:   "Недействительная полезная нагрузка JWT",
				Exception: e.ErrMessage,
			}
		}
	}
	if payload["type"] != "refresh" {
		return nil, &core.ZError{
			Code:      400,
			Where:     "UseCase",
			Message:   "Неверный тип JWT ключа",
			Exception: nil,
		}
	}

	acc_id := payload["sub"].(string)

	res, err := s.repo.GetRefreshTokenForAccount(ctx, acc_id, req.RefreshToken)
	if err != nil {
		switch e := err.(type) {
		case *core.ErrTokenNotFound:
			return nil, &core.ZError{
				Code:      404,
				Where:     "Repo",
				Message:   "Токен не найден",
				Exception: e.ErrMessage,
			}
		case *core.ErrPGRepo:
			return nil, &core.ZError{
				Code:      500,
				Where:     "Repo",
				Message:   "Неизвестная ошибка базы данных",
				Exception: e.ErrMessage,
			}
		}
	}
	if res.IsRevoked {
		return nil, &core.ZError{
			Code:      400,
			Where:     "UseCase",
			Message:   "Токен был отозван",
			Exception: nil,
		}
	} else {
		if user_agent != res.UserAgent || ip != res.IpAddress {
			_, err = s.repo.RevokeToken(ctx, acc_id)
			if err != nil {
				return nil, &core.ZError{
					Code:      500,
					Where:     "Repo",
					Message:   "Неизвестная ошибка базы данных",
					Exception: err,
				}
			}
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase",
				Message:   "Токен отозван",
				Exception: nil,
			}
		}
	}

	access_payload := map[string]interface{}{
		"sub":  acc_id,
		"type": "access",
	}

	token, err := CreateJWT(access_payload, s.cfg.JWTPrivateKey)
	if err != nil {
		switch e := err.(type) {
		case *core.ErrParsePrivateKey:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase/Security",
				Message:   "Ошибка парсинга приватного ключа",
				Exception: e.ErrMessage,
			}
		case *core.ErrSignedJwt:
			return nil, &core.ZError{
				Code:      400,
				Where:     "UseCase/Security",
				Message:   "Ошибка подписания jwt ключа",
				Exception: e.ErrMessage,
			}
		}
	}

	return &share.ZToken{
		AccessToken:  token,
		RefreshToken: req.RefreshToken,
		TokenType:    "bearer",
	}, nil
}

// ----------- Tools -----------

// equal_passwords сравнивает пароль и подтверждение пароля на совпадение.
//
// Параметры:
//   - password: строка с паролем
//   - confirm_pwd: строка с подтверждением пароля
//
// Возвращает:
//   - true, если пароли совпадают
//   - false, если пароли не совпадают
func equal_passwords(password string, confirm_pwd string) bool {
	return password == confirm_pwd
}

// ValidateCredentials проверяет валидность учетных данных пользователя.
//
// Параметры:
//   - email: строка с адресом электронной почты
//   - password: строка с паролем
//
// Возвращает:
//   - true, если учетные данные валидны
//   - ошибку, если длина пароля меньше 6 символов или email не содержит символ "@"
func ValidateCredentials(email, password string) (bool, error) {
	if len(password) < 6 {
		return false, &core.ErrInvalidLenPassword{ErrMessage: nil}
	}
	if !strings.Contains(email, "@") {
		return false, &core.ErrEmailValidate{ErrMessage: nil}
	}
	return true, nil

}
