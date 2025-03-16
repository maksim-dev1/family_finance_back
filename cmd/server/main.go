package main

import (
	"log"

	"family_finance_back/config"
	"family_finance_back/internal/db"
	redisclient "family_finance_back/internal/redis"
	authpkg "family_finance_back/internal/auth"
	userpkg "family_finance_back/internal/user"
	"family_finance_back/internal/ping"
	"family_finance_back/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфигурацию из .env
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализируем подключение к базе данных PostgreSQL
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer database.Close()

	// Инициализируем клиент Redis
	redisClient := redisclient.InitRedis(cfg)

	// Инициализируем сервисы
	authService := authpkg.NewAuthService(database, redisClient, cfg)
	userService := userpkg.NewUserService(database)

	// Инициализируем обработчики
	authHandler := authpkg.NewAuthHandler(authService)
	userHandler := userpkg.NewUserHandler(userService)

	// Создаем роутер Gin
	router := gin.Default()

	// Public routes
	api := router.Group("/api/v1")
	{
		api.GET("/ping", ping.PingHandler)

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/verify", authHandler.Verify)
			authGroup.POST("/refresh", authHandler.Refresh)
			authGroup.POST("/logout", authHandler.Logout)
		}
	}

	// Protected routes (требуют валидного JWT токена)
	protected := api.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(cfg, redisClient))
	{
		protected.GET("/users", userHandler.GetAllUsers)
		protected.GET("/user", userHandler.GetCurrentUser)
		protected.DELETE("/user", userHandler.DeleteUser)
	}

	// Запускаем сервер на порту 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
