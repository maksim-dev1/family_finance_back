// // config/config.go
// package config

// // Config хранит настройки приложения
// type Config struct {
// 	DBPath       string
// 	SMTPHost     string
// 	SMTPPort     int
// 	SMTPUsername string
// 	SMTPPassword string
// }

// // GetConfig возвращает настройки по умолчанию
// func GetConfig() Config {
// 	return Config{
// 		DBPath:       "users.db",            // для SQLite (файл будет создан в корне проекта)
// 		SMTPHost:     "smtp.gmail.com",      // пример для Gmail
// 		SMTPPort:     587,
// 		SMTPUsername: "maks-vasilev-2017@inbox.ru",  // замените на ваш email
// 		SMTPPassword: "mzgm qflx abxs qqkl",   // замените на ваш пароль или используйте App Password (рекомендуется)
// 	}
// }

package config

import (
	"os"
	"strconv"
)

// Config хранит настройки приложения.
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

// GetConfig читает настройки из переменных окружения или использует значения по умолчанию.
func GetConfig() Config {
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		dbPort = 5432
	}
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		smtpPort = 587
	}
	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "familyfinanceDB"),

		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     smtpPort,
		SMTPUsername: getEnv("SMTP_USERNAME", "your-email@gmail.com"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "your-email-password"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
