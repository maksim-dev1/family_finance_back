package db

import (
	"database/sql"
	"fmt"

	"family_finance_back/config"

	_ "github.com/lib/pq"
)

// InitDB устанавливает соединение с PostgreSQL и возвращает объект *sql.DB.
func InitDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
