package repository

import (
	"database/sql"
	"errors"

	"family_finance_back/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	DeleteUserByEmail(email string) error
	GetAllUsers() ([]*models.User, error)
}


type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (id, name, email, created_at, updated_at)
              VALUES ($1, $2, $3, NOW(), NOW())`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Email)
	return err
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)
	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAllUsers() ([]*models.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}


func (r *userRepository) DeleteUserByEmail(email string) error {
	query := `DELETE FROM users WHERE email = $1`
	result, err := r.db.Exec(query, email)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("пользователь не найден")
	}
	return nil
}