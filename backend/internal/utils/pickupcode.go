package utils

import (
	"crypto/rand"
	"encoding/base32"
)

func GeneratePickupCode() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	encoding := base32.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567")
	return encoding.EncodeToString(bytes)[:4]
}
