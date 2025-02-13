package repository

import (
	"database/sql"

	"family_finance_back/internal/models"
)

type SavingsGoalRepository interface {
	CreateSavingsGoal(goal *models.SavingsGoal) error
	GetSavingsGoalsByUser(userID string) ([]*models.SavingsGoal, error)
	GetSavingsGoalsByFamily(familyID string) ([]*models.SavingsGoal, error)
}

type savingsGoalRepository struct {
	db *sql.DB
}

func NewSavingsGoalRepository(db *sql.DB) SavingsGoalRepository {
	return &savingsGoalRepository{db: db}
}

func (r *savingsGoalRepository) CreateSavingsGoal(goal *models.SavingsGoal) error {
	query := `INSERT INTO savings_goals (id, created_by, family_id, target_amount, target_date, start_date, periodic_amount, description, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`
	_, err := r.db.Exec(query, goal.ID, goal.CreatedBy, goal.FamilyID, goal.TargetAmount, goal.TargetDate, goal.StartDate, goal.PeriodicAmount, goal.Description)
	return err
}

func (r *savingsGoalRepository) GetSavingsGoalsByUser(userID string) ([]*models.SavingsGoal, error) {
	query := `SELECT id, created_by, family_id, target_amount, target_date, start_date, periodic_amount, description, created_at, updated_at
              FROM savings_goals WHERE created_by = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var goals []*models.SavingsGoal
	for rows.Next() {
		var goal models.SavingsGoal
		if err := rows.Scan(&goal.ID, &goal.CreatedBy, &goal.FamilyID, &goal.TargetAmount, &goal.TargetDate, &goal.StartDate, &goal.PeriodicAmount, &goal.Description, &goal.CreatedAt, &goal.UpdatedAt); err != nil {
			return nil, err
		}
		goals = append(goals, &goal)
	}
	return goals, nil
}

func (r *savingsGoalRepository) GetSavingsGoalsByFamily(familyID string) ([]*models.SavingsGoal, error) {
	query := `SELECT id, created_by, family_id, target_amount, target_date, start_date, periodic_amount, description, created_at, updated_at
              FROM savings_goals WHERE family_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(query, familyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var goals []*models.SavingsGoal
	for rows.Next() {
		var goal models.SavingsGoal
		if err := rows.Scan(&goal.ID, &goal.CreatedBy, &goal.FamilyID, &goal.TargetAmount, &goal.TargetDate, &goal.StartDate, &goal.PeriodicAmount, &goal.Description, &goal.CreatedAt, &goal.UpdatedAt); err != nil {
			return nil, err
		}
		goals = append(goals, &goal)
	}
	return goals, nil
}
