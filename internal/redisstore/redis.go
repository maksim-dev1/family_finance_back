package redisstore

import (
    "context"
    "time"

    "github.com/go-redis/redis/v8"
    "family_finance_back/config"
)

// NewRedisClient инициализирует новый клиент Redis.
func NewRedisClient(cfg *config.Config) *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     cfg.RedisAddr,
        Password: cfg.RedisPassword,
        DB:       0, // используется база по умолчанию
    })
    return client
}

// SetCode сохраняет в Redis верификационный код для email с указанным временем жизни.
func SetCode(ctx context.Context, client *redis.Client, email, code string, expiration time.Duration) error {
    return client.Set(ctx, "auth_code:"+email, code, expiration).Err()
}

// GetCode возвращает верификационный код для email.
func GetCode(ctx context.Context, client *redis.Client, email string) (string, error) {
    return client.Get(ctx, "auth_code:"+email).Result()
}

// DeleteCode удаляет верификационный код для email.
func DeleteCode(ctx context.Context, client *redis.Client, email string) error {
    return client.Del(ctx, "auth_code:"+email).Err()
}
