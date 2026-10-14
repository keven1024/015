package utils

import (
	"errors"
	"fmt"
	"pkg/utils"

	"golang.org/x/crypto/argon2"
)

var (
	ErrPasswordSaltNotSet = errors.New("PasswordSaltNotSet")
)

func GeneratePasswordHash(password string) (string, error) {
	salt := utils.GetConfig().Share.PasswordSalt
	if salt == "" {
		return "", ErrPasswordSaltNotSet
	}
	hash := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
	return fmt.Sprintf("%x", hash), nil
}
