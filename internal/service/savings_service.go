package service

import (
	"errors"
	"time"

	"family_finance_back/internal/models"
	"family_finance_back/internal/repository"

	"github.com/google/uuid"
)

type SavingsService interface {
	CreateSavingsGoal(createdBy string, familyID *string, targetAmount float64, targetDate, startDate string, description *string) (*models.SavingsGoal, error)
	GetUserSavingsGoals(userID string) ([]*models.SavingsGoal, error)
	GetFamilySavingsGoals(familyID string) ([]*models.SavingsGoal, error)
	CalculateSavingPlan(targetAmount float64, startDate, targetDate time.Time) (daily, weekly, monthly float64, err error)
}

type savingsService struct {
	goalRepo repository.SavingsGoalRepository
}

func NewSavingsService(goalRepo repository.SavingsGoalRepository) SavingsService {
	return &savingsService{
		goalRepo: goalRepo,
	}
}

func (s *savingsService) CreateSavingsGoal(createdBy string, familyID *string, targetAmount float64, targetDate, startDate string, description *string) (*models.SavingsGoal, error) {
	goal := &models.SavingsGoal{
		ID:             uuid.New().String(),
		CreatedBy:      createdBy,
		FamilyID:       familyID,
		TargetAmount:   targetAmount,
		TargetDate:     targetDate,
		StartDate:      startDate,
		PeriodicAmount: nil,
		Description:    description,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	err := s.goalRepo.CreateSavingsGoal(goal)
	if err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *savingsService) GetUserSavingsGoals(userID string) ([]*models.SavingsGoal, error) {
	return s.goalRepo.GetSavingsGoalsByUser(userID)
}

func (s *savingsService) GetFamilySavingsGoals(familyID string) ([]*models.SavingsGoal, error) {
	return s.goalRepo.GetSavingsGoalsByFamily(familyID)
}

func (s *savingsService) CalculateSavingPlan(targetAmount float64, startDate, targetDate time.Time) (daily, weekly, monthly float64, err error) {
	if targetDate.Before(startDate) {
		return 0, 0, 0, errors.New("targetDate must be after startDate")
	}
	days := targetDate.Sub(startDate).Hours() / 24
	daily = targetAmount / days
	weekly = daily * 7
	monthly = daily * 30
	return
}
