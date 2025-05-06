package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
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
