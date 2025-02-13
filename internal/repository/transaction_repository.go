package repository

import (
	"database/sql"
	"errors"
	"time"

	"myapp/internal/models"

	"github.com/lib/pq"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionsByUser(userID string) ([]*models.Transaction, error)
	GetTransactionsByUsers(userIDs []string) ([]*models.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(transaction *models.Transaction) error {
	query := `INSERT INTO transactions (id, user_id, type, category, amount, date, savings_goal_id, description, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`
	_, err := r.db.Exec(query, transaction.ID, transaction.UserID, transaction.Type, transaction.Category, transaction.Amount, transaction.Date, transaction.SavingsGoalID, transaction.Description)
	return err
}

func (r *transactionRepository) GetTransactionsByUser(userID string) ([]*models.Transaction, error) {
	query := `SELECT id, user_id, type, category, amount, date, savings_goal_id, description, created_at, updated_at
              FROM transactions WHERE user_id = $1 ORDER BY date DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions []*models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.UserID, &tx.Type, &tx.Category, &tx.Amount, &tx.Date, &tx.SavingsGoalID, &tx.Description, &tx.CreatedAt, &tx.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, &tx)
	}
	return transactions, nil
}

func (r *transactionRepository) GetTransactionsByUsers(userIDs []string) ([]*models.Transaction, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("нет пользователей")
	}
	query := `SELECT id, user_id, type, category, amount, date, savings_goal_id, description, created_at, updated_at
              FROM transactions WHERE user_id = ANY($1) ORDER BY date DESC`
	rows, err := r.db.Query(query, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions []*models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.UserID, &tx.Type, &tx.Category, &tx.Amount, &tx.Date, &tx.SavingsGoalID, &tx.Description, &tx.CreatedAt, &tx.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, &tx)
	}
	return transactions, nil
}
