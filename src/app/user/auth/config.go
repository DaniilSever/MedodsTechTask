package auth

import (
	"os"
	"sync"
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
	cfg  *Config
	once sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
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
	})
	return cfg
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("Environment variable " + key + " is not set")
	}
	return val
}
