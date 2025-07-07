package utils

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func GeneratePasswordHash(password string) (string, error) {
	salt := GetEnv("PASSWORD_SALT")
	if salt == "" {
		return "", errors.New("请配置PASSWORD_SALT")
	}
	hash := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
	return fmt.Sprintf("%x", hash), nil
}
