package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const ACCESS_TOKEN_NAME = "access-token"

type Claims struct {
	UID    int64  `json:"uid"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func CreateAccessToken(uid int64, userID, role string) (string, error) {
	at := jwt.New(jwt.SigningMethodHS256)

	claims := at.Claims.(jwt.MapClaims)
	claims["uid"] = uid
	claims["user_id"] = userID
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	t, err := at.SignedString([]byte(JWTSecret()))
	if err != nil {
		return "", err
	}

	return t, nil
}

func GetJWTToken(at string) (*Claims, *jwt.Token, error) {
	claims := Claims{}

	t, err := jwt.ParseWithClaims(at, &claims, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret(), nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &claims, t, nil
}
