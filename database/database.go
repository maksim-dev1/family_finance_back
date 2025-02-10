package database

import (
	"database/sql"
	"family_finance_back/config"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// InitDB открывает подключение к PostgreSQL и выполняет миграцию схемы.
func InitDB(cfg config.Config) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка при открытии БД: %v", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		verification_code TEXT,
		is_verified BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы users: %v", err)
	}

	return db
}
