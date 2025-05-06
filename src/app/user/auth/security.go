package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CreatePasswordHash генерирует хеш пароля с использованием соли
//
// Параметры:
//   - pwd: пароль для хеширования
//   - salt: соль для хеширования (если пустая строка, генерируется автоматически)
//
// Возвращает:
//   - хеш парля
//   - использованную соль
//   - ошибку (если возникла)
func CreatePasswordHash(pwd string, salt string) (string, string, error) {
	if pwd == "" {
		return "", "", errors.New("пароль не может быть пустым")
	}

	if salt == "" {
		saltBytes := make([]byte, 16)
		if _, err := rand.Read(saltBytes); err != nil {
			return "", "", fmt.Errorf("ошибка генерации соли: %w", err)
		}
		salt = hex.EncodeToString(saltBytes)
	}

	hash := sha256.New()
	if _, err := hash.Write([]byte(pwd + salt)); err != nil {
		return "", "", fmt.Errorf("ошибка хеширования: %w", err)
	}
	hashSum := hash.Sum(nil)
	hashHex := hex.EncodeToString(hashSum)

	return hashHex, salt, nil
}

func CreateConfirmCode() (string, error) {
	size := 6
	scheme := "0123456789"
	res := make([]byte, size)

	for i := 0; i < size; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(scheme))))

		if err != nil {
			return "", err

		}
		res[i] = scheme[num.Int64()]
	}
	return string(res), nil
}

func CreateJWT(payload map[string]interface{}, private_key string) (string, error) {
	delta := 24 * 60 * 60
	now := time.Now().UTC()
	exp := now.Add(time.Duration(delta) * time.Second)

	claims := jwt.MapClaims{
		"iat":    now.Unix(),
		"iss":    "MedodsTechTask",
		"exp":    exp.Unix(),
		"exp_at": exp.Format(time.RFC3339),
		"exp_in": delta,
	}

	for k, v := range payload {
		claims[k] = v
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(private_key))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func DecodeJWT(token_string string, public_key string) (map[string]interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(strings.TrimSpace(public_key)))
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("JWT token has expired")
		}
		return nil, errors.New("incorrect jwt")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid JWT payload")
}
