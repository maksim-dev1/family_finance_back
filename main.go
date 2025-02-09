// main.go
package main

import (
	"family_finance_back/config"
	"family_finance_back/database"
	"family_finance_back/handlers"
	"family_finance_back/repository"
	"family_finance_back/service"
	"log"
	"net/http"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.GetConfig()

	// Инициализируем базу данных
	db := database.InitDB(cfg.DBPath)
	defer db.Close()

	// Создаем репозиторий пользователей
	userRepo := repository.NewUserRepository(db)

	// Создаем сервис аутентификации
	authService := service.NewAuthService(userRepo, cfg)

	// Создаем обработчики
	authHandler := handlers.NewAuthHandler(authService)

	// Настраиваем маршруты
	http.HandleFunc("/send-code", authHandler.SendCode)
	http.HandleFunc("/verify-code", authHandler.VerifyCode)

	log.Println("Сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
