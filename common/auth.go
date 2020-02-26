package common

import (
	"github.com/PPIO/pi-cloud-monitor-backend/config"
	"github.com/dgrijalva/jwt-go"
)

type JWTCustomClaims struct {
	UID uint64 `json:"uid"`
	jwt.StandardClaims
}

func GetToken(uid uint64, expireAt int64) (*string, error) {
	claims := &JWTCustomClaims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: expireAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	tokenStr, err := token.SignedString(config.GetJwtPrivateKey())
	if err != nil {
		return nil, err
	}
	return &tokenStr, nil
}
