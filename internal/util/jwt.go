package util

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(email string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().AddDate(100, 0, 0).Unix(), // очень большой срок жизни
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
