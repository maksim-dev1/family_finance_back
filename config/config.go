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
		SMTPUsername: getEnv("SMTP_USERNAME", "maks-vasilev-2017@inbox.ru"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "mzgmqflxabxsqqkl"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
