package models

import "time"

// User описывает пользователя.
type User struct {
	ID               int       `json:"id"`
	Email            string    `json:"email"`
	Name             string    `json:"name,omitempty"`
	VerificationCode string    `json:"verification_code"`
	CodeExpiresAt    time.Time `json:"code_expires_at"` // время истечения кода
	IsVerified       bool      `json:"is_verified"`
	CreatedAt        time.Time `json:"created_at"`
}
