package auth

import (
    "database/sql"
    "time"

    "family_finance_back/internal/models"
    "github.com/google/uuid"
)

// Repository описывает методы работы с данными пользователей.
type Repository interface {
    GetUserByEmail(email string) (*models.User, error)
    CreateUser(user *models.User) error
}

// PostgresRepository реализует Repository с использованием PostgreSQL.
type PostgresRepository struct {
    db *sql.DB
}

// NewPostgresRepository создаёт новый экземпляр PostgresRepository.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
    return &PostgresRepository{db: db}
}

// GetUserByEmail возвращает пользователя по email.
func (r *PostgresRepository) GetUserByEmail(email string) (*models.User, error) {
    query := `SELECT id, name, email, created_at, updated_at FROM users WHERE email=$1`
    row := r.db.QueryRow(query, email)
    user := &models.User{}
    var id string
    err := row.Scan(&id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // пользователь не найден
        }
        return nil, err
    }
    user.ID, err = uuid.Parse(id)
    if err != nil {
        return nil, err
    }
    return user, nil
}

// CreateUser создаёт нового пользователя в базе данных.
func (r *PostgresRepository) CreateUser(user *models.User) error {
    if user.ID == uuid.Nil {
        user.ID = uuid.New()
    }
    now := time.Now().UTC()
    user.CreatedAt = now
    user.UpdatedAt = now

    query := `INSERT INTO users (id, name, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
    _, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
    if err != nil {
        return err
    }
    return nil
}
