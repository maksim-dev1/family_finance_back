package config

import (
    "log"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

// Config хранит все настройки приложения.
type Config struct {
    DBHost        string
    DBPort        string
    DBUser        string
    DBPassword    string
    DBName        string
    RedisAddr     string
    RedisPassword string
    SMTPHost      string
    SMTPPort      int
    SMTPUsername  string
    SMTPPassword  string
    JWTSecret     string
}

// LoadConfig загружает настройки из .env файла и переменных окружения.
func LoadConfig() (*Config, error) {
    // Загружаем .env файл
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, reading environment variables")
    }

    smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
    if err != nil {
        return nil, err
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
        SMTPPort:      smtpPort,
        SMTPUsername:  os.Getenv("SMTP_USERNAME"),
        SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
        JWTSecret:     os.Getenv("JWT_SECRET"),
    }
    return config, nil
}
