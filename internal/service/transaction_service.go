package service

import (
	"time"

	"family_finance_back/internal/models"
	"family_finance_back/internal/repository"

	"github.com/google/uuid"
)

type TransactionService interface {
	CreateTransaction(userID, txType, category string, amount float64, date time.Time, savingsGoalID *string, description *string) (*models.Transaction, error)
	GetPersonalTransactions(userID string) ([]*models.Transaction, error)
	GetGroupTransactions(userIDs []string) ([]*models.Transaction, error)
}

type transactionService struct {
	txRepo repository.TransactionRepository
}

func NewTransactionService(txRepo repository.TransactionRepository) TransactionService {
	return &transactionService{
		txRepo: txRepo,
	}
}

func (s *transactionService) CreateTransaction(userID, txType, category string, amount float64, date time.Time, savingsGoalID *string, description *string) (*models.Transaction, error) {
	tx := &models.Transaction{
		ID:            uuid.New().String(),
		UserID:        userID,
		Type:          txType,
		Category:      category,
		Amount:        amount,
		Date:          date,
		SavingsGoalID: savingsGoalID,
		Description:   description,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err := s.txRepo.CreateTransaction(tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *transactionService) GetPersonalTransactions(userID string) ([]*models.Transaction, error) {
	return s.txRepo.GetTransactionsByUser(userID)
}

func (s *transactionService) GetGroupTransactions(userIDs []string) ([]*models.Transaction, error) {
	return s.txRepo.GetTransactionsByUsers(userIDs)
}
