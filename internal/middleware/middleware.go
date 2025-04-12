package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
)

func JWTAuthMiddleware(redisClient *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			next.ServeHTTP(w, r)
		})
	}
}
