package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"family_finance_back/config"
	"family_finance_back/internal/handler"
	"family_finance_back/internal/middleware"
	"family_finance_back/internal/repository"
	"family_finance_back/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить .env файл")
	}

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	// Подключаемся к PostgreSQL
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Инициализируем Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	// Инициализируем репозитории
	userRepo := repository.NewUserRepository(db)
	familyRepo := repository.NewFamilyRepository(db)
	familyMemberRepo := repository.NewFamilyMemberRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	savingsGoalRepo := repository.NewSavingsGoalRepository(db)

	// Инициализируем сервисы
	authService := service.NewAuthService(userRepo, rdb, cfg)
	emailService := service.NewEmailService(cfg)
	familyService := service.NewFamilyService(familyRepo, familyMemberRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	savingsService := service.NewSavingsService(savingsGoalRepo)
	gamificationService := service.NewGamificationService(transactionRepo)
	syncService := service.NewSyncService(transactionRepo)

	// Инициализируем обработчики
	authHandler := handler.NewAuthHandler(authService, emailService)
	familyHandler := handler.NewFamilyHandler(familyService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	savingsHandler := handler.NewSavingsHandler(savingsService)
	gamificationHandler := handler.NewGamificationHandler(gamificationService)
	syncHandler := handler.NewSyncHandler(syncService)

	// Настраиваем роутер Gin
	router := gin.Default()

	// Публичные маршруты (регистрация, логин, верификация кода)
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/verify", authHandler.VerifyCode)
	}

	// Протектед маршруты с JWT middleware
	apiRoutes := router.Group("/api")
	apiRoutes.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	{
		// Family routes
		apiRoutes.POST("/families", familyHandler.CreateFamily)
		apiRoutes.GET("/families", familyHandler.GetFamilies)

		// Transaction routes
		apiRoutes.POST("/transactions", transactionHandler.CreateTransaction)
		apiRoutes.GET("/transactions/personal", transactionHandler.GetPersonalTransactions)
		apiRoutes.GET("/transactions/group", transactionHandler.GetGroupTransactions)

		// Savings routes
		apiRoutes.POST("/savings", savingsHandler.CreateSavingsGoal)
		apiRoutes.GET("/savings", savingsHandler.GetSavingsGoals)
		apiRoutes.POST("/savings/calculate", savingsHandler.CalculateSavingPlan)

		// Gamification routes
		apiRoutes.GET("/gamification/score", gamificationHandler.GetUserScore)

		// Sync routes
		apiRoutes.POST("/sync/transactions", syncHandler.SyncTransactions)
	}

	// Запускаем сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Сервер запущен на порту %s", port)
	router.Run(fmt.Sprintf(":%s", port))
}
