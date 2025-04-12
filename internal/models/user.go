package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Surname   string    `gorm:"size:100;not null" json:"surname"`
	Nickname  string    `gorm:"size:100;not null" json:"nickname"`
	Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
