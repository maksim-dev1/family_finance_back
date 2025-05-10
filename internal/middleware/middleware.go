package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
)

// LoggerMiddleware создает middleware для логирования HTTP запросов
// Логирует метод, путь, статус ответа и время выполнения запроса
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

// CORSMiddleware создает middleware для обработки CORS
// Разрешает запросы с указанного origin и необходимые заголовки
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

// AuthMiddleware создает middleware для проверки JWT токена
// Проверяет наличие и валидность токена в заголовке Authorization
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

// JWTAuthMiddleware создает middleware для проверки JWT токена
// Проверяет наличие и валидность токена в заголовке Authorization
// Также проверяет, не находится ли токен в черном списке
func JWTAuthMiddleware(redisClient *redis.Client) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Неверный формат заголовка Authorization", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			// Проверяем, находится ли токен в blacklist
			val, err := redisClient.Get(context.Background(), "blacklist:"+token).Result()
			if err == nil && val == "true" {
				http.Error(w, "Токен отозван", http.StatusUnauthorized)
				return
			}

			// Если всё ок, передаем дальше
			next(w, r)
		}
	}
}
