package redis

import (
	"family_finance_back/config"

	"github.com/go-redis/redis/v8"
)

// InitRedis инициализирует и возвращает Redis клиент.
func InitRedis(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	return rdb
}
