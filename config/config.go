package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
<<<<<<< HEAD
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBURL         string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	SMTPHost      string
	SMTPPort      string
	SMTPUsername  string
	SMTPPassword  string
	JWTSecret     string
=======
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

	JWTSecret        string // секрет для подписи JWT
	JWTExpiryMinutes int    // время жизни токена в минутах
>>>>>>> a87482f9ae3a0c3f31f94620ce2de9b4ff6244d3
}

func LoadConfig() (*Config, error) {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0 // по умолчанию используем 0
	}
<<<<<<< HEAD

	cfg := &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SMTPUsername:  os.Getenv("SMTP_USERNAME"),
		SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
=======
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		smtpPort = 587
	}
	jwtExpiry, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_MINUTES"))
	if err != nil {
		jwtExpiry = 60 // по умолчанию 60 минут
	}

	return Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           dbPort,
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "password"),
		DBName:           getEnv("DB_NAME", "familyfinanceDB"),
		SMTPHost:         getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:         smtpPort,
		SMTPUsername:     getEnv("SMTP_USERNAME", "your_email@example.com"),
		SMTPPassword:     getEnv("SMTP_PASSWORD", "your_email_password"),
		JWTSecret:        getEnv("JWT_SECRET", "your_secret_key"),
		JWTExpiryMinutes: jwtExpiry,
>>>>>>> a87482f9ae3a0c3f31f94620ce2de9b4ff6244d3
	}
	cfg.DBURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	return cfg, nil
}
