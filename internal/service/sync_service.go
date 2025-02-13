package service

import (
	"time"

	"github.com/google/uuid"
	"myapp/internal/models"
	"myapp/internal/repository"
)

type SyncService interface {
	SyncTransactions(userID string, transactions []*models.Transaction) error
}

type syncService struct {
	txRepo repository.TransactionRepository
}

func NewSyncService(txRepo repository.TransactionRepository) SyncService {
	return &syncService{
		txRepo: txRepo,
	}
}

func (s *syncService) SyncTransactions(userID string, transactions []*models.Transaction) error {
	for _, tx := range transactions {
		if tx.ID == "" {
			tx.ID = uuid.New().String()
		}
		tx.UserID = userID
		tx.CreatedAt = time.Now()
		tx.UpdatedAt = time.Now()
		err := s.txRepo.CreateTransaction(tx)
		if err != nil {
			return err
		}
	}
	return nil
}
