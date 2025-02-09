// config/config.go
package config

// Config хранит настройки приложения
type Config struct {
	DBPath       string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

// GetConfig возвращает настройки по умолчанию
func GetConfig() Config {
	return Config{
		DBPath:       "users.db",            // для SQLite (файл будет создан в корне проекта)
		SMTPHost:     "smtp.gmail.com",      // пример для Gmail
		SMTPPort:     587,
		SMTPUsername: "your_email@gmail.com",  // замените на ваш email
		SMTPPassword: "your_email_password",   // замените на ваш пароль или используйте App Password (рекомендуется)
	}
}
