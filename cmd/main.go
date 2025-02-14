package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"family_finance_back/config"
	"family_finance_back/internal/db"
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

	// Инициализируем PostgreSQL с созданием таблиц (миграцией)
	db := database.InitDB(*cfg)
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
	userService := service.NewUserService(userRepo)
	emailService := service.NewEmailService(cfg)
	familyService := service.NewFamilyService(familyRepo, familyMemberRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	savingsService := service.NewSavingsService(savingsGoalRepo)
	gamificationService := service.NewGamificationService(transactionRepo)
	syncService := service.NewSyncService(transactionRepo)

	// Инициализируем обработчики
	authHandler := handler.NewAuthHandler(authService, emailService)
	userHandler := handler.NewUserHandler(userService)
	familyHandler := handler.NewFamilyHandler(familyService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	savingsHandler := handler.NewSavingsHandler(savingsService)
	gamificationHandler := handler.NewGamificationHandler(gamificationService)
	syncHandler := handler.NewSyncHandler(syncService)

	// Настраиваем роутер Gin
	router := gin.Default()

	// Добавляем endpoint для проверки доступности сервера
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

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
		// User routes
		userGroup := apiRoutes.Group("/user")
		{
			userGroup.GET("/get/all", userHandler.GetAllUsers)
			userGroup.GET("/get/by/email", userHandler.GetUserByEmail)
			userGroup.DELETE("/delete/by/email", userHandler.DeleteUser)
		}

		// Family routes
		familyRoutes := router.Group("/family")
		{
			familyRoutes.POST("/create", familyHandler.CreateFamily)
			familyRoutes.GET("/get", familyHandler.GetFamilies)
			familyRoutes.POST("/join", familyHandler.JoinFamily)
		}

		// Transaction routes
		transactionRoutes := router.Group("/transactions")
		{
			transactionRoutes.POST("/create", transactionHandler.CreateTransaction)
			transactionRoutes.GET("/get/personal", transactionHandler.GetPersonalTransactions)
			transactionRoutes.GET("/get/group", transactionHandler.GetGroupTransactions)
		}

		// Savings routes
		savingsRoutes := router.Group("/savings")
		{
			savingsRoutes.POST("/create", savingsHandler.CreateSavingsGoal)
			savingsRoutes.GET("/get", savingsHandler.GetSavingsGoals)
			savingsRoutes.POST("/calculate", savingsHandler.CalculateSavingPlan)
		}

		// Gamification routes
		gamificationRoutes := router.Group("/gamification")
		{
			gamificationRoutes.GET("/get/score", gamificationHandler.GetUserScore)
		}

		// Sync routes
		syncRoutes := router.Group("/sync")
		{
			syncRoutes.POST("/transactions", syncHandler.SyncTransactions)
		}
	}

	// Запускаем сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Сервер запущен на порту %s", port)
	router.Run(fmt.Sprintf(":%s", port))
}
