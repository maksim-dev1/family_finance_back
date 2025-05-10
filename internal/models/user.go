package models

import "time"

// User представляет модель пользователя в системе
type User struct {
	// ID уникальный идентификатор пользователя
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// Name имя пользователя
	Name string `gorm:"size:100;not null" json:"name"`

	// Surname фамилия пользователя
	Surname string `gorm:"size:100;not null" json:"surname"`

	// Nickname псевдоним пользователя
	Nickname string `gorm:"size:100;not null" json:"nickname"`

	// Email электронная почта пользователя (уникальная)
	Email string `gorm:"size:100;uniqueIndex;not null" json:"email"`

	// CreatedAt время создания записи
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt время последнего обновления записи
	UpdatedAt time.Time `json:"updated_at"`
}
