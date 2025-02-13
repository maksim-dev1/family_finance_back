package service

import (
	"family_finance_back/internal/repository"
)

type GamificationService interface {
	GetUserScore(userID string) (int, error)
}

type gamificationService struct {
	txRepo repository.TransactionRepository
}

func NewGamificationService(txRepo repository.TransactionRepository) GamificationService {
	return &gamificationService{
		txRepo: txRepo,
	}
}

func (s *gamificationService) GetUserScore(userID string) (int, error) {
	txs, err := s.txRepo.GetTransactionsByUser(userID)
	if err != nil {
		return 0, err
	}
	score := len(txs) * 10
	return score, nil
}
