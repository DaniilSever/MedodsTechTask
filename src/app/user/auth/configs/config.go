package configs

import (
	"os"
)

type Config struct {
	AppEnv string
	// DB
	AuthDBDBN    string
	AuthDBUSR    string
	AuthDBPWD    string
	DBMasterHost string
	DBMasterPort string
	// JWT
	JWTPublicKey          string
	JWTPrivateKey         string
	AuthJWTTokenExpireMin int
}

var (
	cfg *Config
)

// GetConfig возвращает конфигурацию приложения в зависимости от значения переменной окружения APP_ENV.
//
// Поддерживаемые значения APP_ENV:
//   - "local": конфигурация для локальной среды
//   - "test": конфигурация для тестовой среды
//
// Возвращает:
//   - указатель на структуру конфигурации Config
func GetConfig() *Config {
	env := os.Getenv("APP_ENV")
	switch env {
	case "local":
		cfg = &Config{
			AppEnv: "local",
			// DB
			AuthDBDBN:    "auth",
			AuthDBUSR:    "auth",
			AuthDBPWD:    getEnv("AUTH_DB_PWD"),
			DBMasterHost: "175.243.0.50",
			DBMasterPort: "5432",
			// JWT
			JWTPublicKey:          getEnv("JWT_PUBLIC_KEY"),
			JWTPrivateKey:         getEnv("JWT_PRIVATE_KEY"),
			AuthJWTTokenExpireMin: 60 * 24,
		}
	case "test":
		cfg = &Config{
			AppEnv: "local",
			// DB
			AuthDBDBN:    "auth",
			AuthDBUSR:    "auth",
			AuthDBPWD:    "auth_pwd",
			DBMasterHost: "175.243.0.50",
			DBMasterPort: "5433",
			// JWT
			JWTPublicKey:          "testjwt",
			JWTPrivateKey:         "testjwt",
			AuthJWTTokenExpireMin: 60 * 24,
		}
	}
	return cfg
}

// getEnv возвращает значение переменной окружения по ключу.
//
// Параметры:
//   - key: имя переменной окружения
//
// Возвращает:
//   - значение переменной окружения
//
// Паника:
//   - если переменная окружения не установлена
func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("Environment variable " + key + " is not set")
	}
	return val
}
