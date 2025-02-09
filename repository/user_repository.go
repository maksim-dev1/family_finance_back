// repository/user_repository.go
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
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	query := "SELECT id, email, verification_code, is_verified, created_at FROM users WHERE email = ?"
	row := r.db.QueryRow(query, email)
	var user models.User
	var isVerifiedInt int
	err := row.Scan(&user.ID, &user.Email, &user.VerificationCode, &isVerifiedInt, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // пользователь не найден
		}
		return nil, err
	}
	user.IsVerified = isVerifiedInt == 1
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (email, verification_code, is_verified) VALUES (?, ?, ?)"
	isVerifiedInt := 0
	if user.IsVerified {
		isVerifiedInt = 1
	}
	_, err := r.db.Exec(query, user.Email, user.VerificationCode, isVerifiedInt)
	return err
}

func (r *userRepository) UpdateUserVerification(email string, code string, isVerified bool) error {
	query := "UPDATE users SET verification_code = ?, is_verified = ? WHERE email = ?"
	isVerifiedInt := 0
	if isVerified {
		isVerifiedInt = 1
	}
	result, err := r.db.Exec(query, code, isVerifiedInt, email)
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
