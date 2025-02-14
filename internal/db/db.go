package database

import (
	"database/sql"
	"family_finance_back/config"
	"fmt"
	"log"
)

// InitDB открывает подключение к PostgreSQL и выполняет миграцию схемы.
func InitDB(cfg config.Config) *sql.DB {
	// Формируем строку подключения. Обратите внимание, что порт в конфиге — строка.
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка при открытии БД: %v", err)
	}

	// Создаем таблицу пользователей
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы users: %v", err)
	}

	// Создаем таблицу семей
	createFamiliesTable := `
	CREATE TABLE IF NOT EXISTS families (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createFamiliesTable)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы families: %v", err)
	}

	// Создаем таблицу участников семей
	createFamilyMembersTable := `
	CREATE TABLE IF NOT EXISTS family_members (
		id TEXT PRIMARY KEY,
		family_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		role TEXT NOT NULL,
		joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	_, err = db.Exec(createFamilyMembersTable)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы family_members: %v", err)
	}

	// Создаем таблицу целей накоплений
	createSavingsGoalsTable := `
	CREATE TABLE IF NOT EXISTS savings_goals (
		id TEXT PRIMARY KEY,
		created_by TEXT NOT NULL,
		family_id TEXT,
		target_amount NUMERIC NOT NULL,
		target_date DATE NOT NULL,
		start_date DATE NOT NULL,
		periodic_amount NUMERIC,
		description TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE
	);`
	_, err = db.Exec(createSavingsGoalsTable)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы savings_goals: %v", err)
	}

	// Создаем таблицу транзакций
	createTransactionsTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		type TEXT NOT NULL,
		category TEXT NOT NULL,
		amount NUMERIC NOT NULL,
		date TIMESTAMP NOT NULL,
		savings_goal_id TEXT,
		description TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (savings_goal_id) REFERENCES savings_goals(id) ON DELETE SET NULL
	);`
	_, err = db.Exec(createTransactionsTable)
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы transactions: %v", err)
	}

	log.Println("Все таблицы успешно созданы (если их ещё не было)")
	return db
}
