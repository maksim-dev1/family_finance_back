package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWTToken генерирует JWT-токен с email в качестве claims.
func GenerateJWTToken(email string, secret string, expiryMinutes int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Minute * time.Duration(expiryMinutes)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateJWTToken проверяет JWT-токен и возвращает email из claims.
func ValidateJWTToken(tokenStr string, secret string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return "", errors.New("не удалось получить email из токена")
		}
		return email, nil
	}
	return "", errors.New("недействительный токен")
}
