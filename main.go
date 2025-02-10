package main

import (
	"context"
	"family_finance_back/config"
	"family_finance_back/database"
	"family_finance_back/handlers"
	"family_finance_back/repository"
	"family_finance_back/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Нет файла .env, продолжаем использовать системные переменные окружения")
	}

	// Загружаем конфигурацию из переменных окружения
	cfg := config.GetConfig()

	// Инициализируем базу данных PostgreSQL
	db := database.InitDB(cfg)
	defer db.Close()

	// Создаем репозиторий, сервис и обработчики
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	// Настраиваем роутер с Gorilla Mux
	router := mux.NewRouter()
	router.HandleFunc("/send-code", authHandler.SendCode).Methods("POST")
	router.HandleFunc("/verify-code", authHandler.VerifyCode).Methods("POST")
	router.HandleFunc("/ping", handlers.PingHandler).Methods("GET")

	// Создаем HTTP-сервер с таймаутами
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Println("Сервер запущен на порту 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Сервер останавливается...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}
	log.Println("Сервер завершил работу")
}
