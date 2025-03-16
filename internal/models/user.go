package models

import "time"

// User представляет структуру пользователя в системе Family Finance.
type User struct {
	ID        string    `json:"id"`         // UUID
	Name      string    `json:"name"`       // Имя пользователя
	Email     string    `json:"email"`      // Email (уникальное, используется для авторизации)
	CreatedAt time.Time `json:"created_at"` // Дата регистрации
	UpdatedAt time.Time `json:"updated_at"` // Дата последнего обновления профиля
}
