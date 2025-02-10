package models

import "time"

// User описывает пользователя.
type User struct {
	ID               int       `json:"id"`
	Email            string    `json:"email"`
	Name             string    `json:"name,omitempty"`
	VerificationCode string    `json:"verification_code"`
	IsVerified       bool      `json:"is_verified"`
	CreatedAt        time.Time `json:"created_at"`
}
