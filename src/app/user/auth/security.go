package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
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
		return "", "", errors.New("Пароль не может быть пустым")
	}

	if salt == "" {
		saltBytes := make([]byte, 16)
		if _, err := rand.Read(saltBytes); err != nil {
			return "", "", fmt.Errorf("Ошибка генерации соли: %w", err)
		}
		salt = hex.EncodeToString(saltBytes)
	}

	hash := sha256.New()
	if _, err := hash.Write([]byte(pwd + salt)); err != nil {
		return "", "", fmt.Errorf("Ошибка хеширования: %w", err)
	}
	hashSum := hash.Sum(nil)
	hashHex := hex.EncodeToString(hashSum)

	return hashHex, salt, nil
}
