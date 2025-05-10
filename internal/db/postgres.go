package db

import (
	"fmt"
	"log"

	"family_finance_back/config"
	"family_finance_back/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitPostgres инициализирует подключение к PostgreSQL
// Использует параметры подключения из конфигурации
// Возвращает экземпляр *gorm.DB для работы с базой данных
func InitPostgres(cfg config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	// Автоматическая миграция (создание таблиц, если их нет)
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("error during migration: %v", err)
	}

	return db
}
