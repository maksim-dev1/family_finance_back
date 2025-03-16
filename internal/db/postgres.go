package db

import (
	"database/sql"
	"fmt"

	"family_finance_back/config"

	_ "github.com/lib/pq"
)

// NewPostgresDB инициализирует подключение к PostgreSQL.
func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Автоматически запускаем миграцию для создания таблицы пользователей, если её нет.
	if err = Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate выполняет миграции базы данных, создавая таблицу пользователей если её нет.
func Migrate(db *sql.DB) error {
	// Создаём расширение для генерации UUID, если оно не существует.
	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		return fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	// Запрос на создание таблицы пользователей, если она отсутствует.
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR NOT NULL,
		email VARCHAR UNIQUE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL
	);`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil
}
