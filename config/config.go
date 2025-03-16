package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config хранит все переменные окружения, необходимые для работы приложения.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisAddr     string
	RedisPassword string

	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string

	JWTSecret string
}

// LoadConfig загружает переменные окружения из файла .env (если он есть) и возвращает Config.
func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, читаем переменные окружения")
	}

	config := &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SMTPUsername:  os.Getenv("SMTP_USERNAME"),
		SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
	}

	return config, nil
}
