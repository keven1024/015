package services

import (
	u "pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

type DownloadShareClaims struct {
	ShareId string `json:"share_id"`
	jwt.RegisteredClaims
}

func ParseDownloadToken(token string) (DownloadShareClaims, error) {
	claims := DownloadShareClaims{}
	t, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.GetConfig().Share.DownloadSecret), nil
	})
	if err != nil {
		return DownloadShareClaims{}, err
	}
	if t == nil || !t.Valid || claims.ShareId == "" {
		return DownloadShareClaims{}, ErrInvalidDownloadToken
	}
	return claims, nil
}
