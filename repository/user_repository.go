package repository

import (
	"database/sql"
	"errors"
	"family_finance_back/models"
)

type UserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUserVerification(email string, code string, isVerified bool) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	query := "SELECT id, email, verification_code, is_verified, created_at FROM users WHERE email = $1"
	row := r.db.QueryRow(query, email)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.VerificationCode, &user.IsVerified, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // пользователь не найден
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (email, verification_code, is_verified) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, user.Email, user.VerificationCode, user.IsVerified)
	return err
}

func (r *userRepository) UpdateUserVerification(email string, code string, isVerified bool) error {
	query := "UPDATE users SET verification_code = $1, is_verified = $2 WHERE email = $3"
	result, err := r.db.Exec(query, code, isVerified, email)
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
