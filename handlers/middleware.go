package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"family_finance_back/utils"
)

type contextKey string

const (
	// ContextKeyUserEmail используется для сохранения email пользователя в контексте запроса.
	ContextKeyUserEmail contextKey = "userEmail"
)

// TokenAuthMiddleware проверяет наличие и корректность JWT-токена.
func TokenAuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "Отсутствует токен")
				return
			}
			// Ожидаем формат "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondWithError(w, http.StatusUnauthorized, "Неверный формат токена")
				return
			}
			tokenStr := parts[1]
			email, err := utils.ValidateJWTToken(tokenStr, jwtSecret)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Неверный токен: %v", err))
				return
			}
			// Записываем email пользователя в контекст запроса
			ctx := context.WithValue(r.Context(), ContextKeyUserEmail, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
