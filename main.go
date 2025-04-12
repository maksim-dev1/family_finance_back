package main

import (
	"log"
	"net/http"

	"family_finance_back/config"
	"family_finance_back/internal/db"
	"family_finance_back/internal/handlers"
	"family_finance_back/internal/repository"
	"family_finance_back/internal/service"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Загружаем конфигурацию из .env
	cfg := config.LoadConfig()

	// Инициализируем PostgreSQL
	postgresDB := db.InitPostgres(cfg)

	// Инициализируем Redis клиент
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	// Инициализируем репозитории
	userRepo := repository.NewUserRepository(postgresDB)

	// Инициализируем сервисы
	emailSvc := service.NewEmailService(cfg)
	authSvc := service.NewAuthService(userRepo, emailSvc, redisClient, cfg.JWTSecret)

	// Инициализируем обработчики
	authHandler := handlers.NewAuthHandler(authSvc)

	// Настраиваем маршруты
	http.HandleFunc("/login/request", authHandler.RequestLoginCodeHandler)
	http.HandleFunc("/login/verify", authHandler.VerifyLoginCodeHandler)
	http.HandleFunc("/register/request", authHandler.RequestRegistrationCodeHandler)
	http.HandleFunc("/register/verify", authHandler.VerifyRegistrationCodeHandler)

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
