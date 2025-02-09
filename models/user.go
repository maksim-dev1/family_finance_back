// models/user.go
package models

import "time"

// User описывает пользователя
type User struct {
	ID               int
	Email            string
	Name             string
	VerificationCode string
	IsVerified       bool
	CreatedAt        time.Time
}
