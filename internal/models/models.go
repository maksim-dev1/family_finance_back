package models

import "time"

// Пользователь
type User struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Семья/группа
type Family struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Член семьи
type FamilyMember struct {
	ID       string    `db:"id"`
	FamilyID string    `db:"family_id"`
	UserID   string    `db:"user_id"`
	Role     string    `db:"role"`
	JoinedAt time.Time `db:"joined_at"`
}

// Транзакция (доход, расход, накопление)
type Transaction struct {
	ID             string     `db:"id"`
	UserID         string     `db:"user_id"`
	Type           string     `db:"type"` // expense, income, saving
	Category       string     `db:"category"`
	Amount         float64    `db:"amount"`
	Date           time.Time  `db:"date"`
	SavingsGoalID  *string    `db:"savings_goal_id"`
	Description    *string    `db:"description"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

// Цель накоплений
type SavingsGoal struct {
	ID             string    `db:"id"`
	CreatedBy      string    `db:"created_by"`
	FamilyID       *string   `db:"family_id"`
	TargetAmount   float64   `db:"target_amount"`
	TargetDate     string    `db:"target_date"`
	StartDate      string    `db:"start_date"`
	PeriodicAmount *float64  `db:"periodic_amount"`
	Description    *string   `db:"description"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
