package user

import (
	"database/sql"
	"errors"

	"family_finance_back/internal/models"
)

// UserService предоставляет методы для работы с данными пользователей.
type UserService struct {
	db *sql.DB
}

// NewUserService создаёт новый экземпляр UserService.
func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// GetAllUsers возвращает список всех пользователей.
func (s *UserService) GetAllUsers() ([]models.User, error) {
	rows, err := s.db.Query("SELECT id, name, email, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUserByEmail возвращает данные пользователя по email.
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow("SELECT id, name, email, created_at, updated_at FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUserByEmail удаляет пользователя по email.
func (s *UserService) DeleteUserByEmail(email string) error {
	result, err := s.db.Exec("DELETE FROM users WHERE email=$1", email)
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
