package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"family_finance_back/config"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
)

// JWTAuthMiddleware проверяет наличие и корректность JWT токена в заголовке Authorization.
// Также проверяет, что токен не находится в blacklist (т.е. не был отозван).
func JWTAuthMiddleware(cfg *config.Config, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "требуется заголовок Authorization"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "формат заголовка: Bearer {token}"})
			return
		}
		tokenString := parts[1]

		// Проверяем, не находится ли токен в blacklist
		ctx := context.Background()
		exists, err := redisClient.Exists(ctx, "blacklist:"+tokenString).Result()
		if err == nil && exists == 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "токен недействителен"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "неверный токен"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "неверные данные токена"})
			return
		}
		// Проверяем время истечения токена
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "срок действия токена истёк"})
				return
			}
		}
		// Извлекаем email пользователя и сохраняем в контексте
		email, ok := claims["email"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "неверная нагрузка токена"})
			return
		}
		c.Set("user_email", email)
		c.Next()
	}
}
