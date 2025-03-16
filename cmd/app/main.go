package main

import (
	"log"
	"net/http"
	"time"

	"family_finance_back/config"
	"family_finance_back/internal/auth"
	"family_finance_back/internal/db"
	redisstore "family_finance_back/internal/redisstore"
	"family_finance_back/internal/smtp"
	"family_finance_back/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Инициализация подключения к PostgreSQL
	dbConn, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer dbConn.Close()

	// Инициализация клиента Redis
	redisClient := redisstore.NewRedisClient(cfg)

	// Инициализация SMTP mailer
	mailer := smtp.NewMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword)

	// Инициализация репозитория и сервиса аутентификации
	authRepo := auth.NewPostgresRepository(dbConn)
	authService := auth.NewAuthService(authRepo, redisClient, mailer, cfg.JWTSecret)
	authHandler := auth.NewAuthHandler(authService)

	// Переводим Gin в release режим для продакшена
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Эндпоинт для проверки работоспособности сервера
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Роуты для авторизации
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/send-code", authHandler.SendCode)
		authRoutes.POST("/verify-code", authHandler.VerifyCode)
	}

	// Пример защищённого маршрута
	protected := router.Group("/api")
	protected.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	{
		protected.GET("/profile", func(c *gin.Context) {
			user, exists := c.Get("user")
			if !exists {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found in context"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"user": user})
		})
	}

	// Запуск сервера
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Server is running on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
