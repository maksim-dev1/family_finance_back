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
		SMTPUsername: "maks-vasilev-2017@inbox.ru",  // замените на ваш email
		SMTPPassword: "mzgm qflx abxs qqkl",   // замените на ваш пароль или используйте App Password (рекомендуется)
	}
}
