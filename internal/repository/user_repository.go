package repository

import (
	"database/sql"
	"errors"
	"log"

	"family_finance_back/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	const operation = "userRepository.CreateUser"
	
	query := `INSERT INTO users (id, name, email, created_at, updated_at)
              VALUES ($1, $2, $3, NOW(), NOW())`

	log.Printf("%s: executing query\n%s\nWith params: [%s, %s, %s]",
		operation,
		query,
		user.ID,
		user.Name,
		user.Email,
	)

	result, err := r.db.Exec(query, user.ID, user.Name, user.Email)
	if err != nil {
		log.Printf("%s: ERROR in Exec - %v", operation, err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("%s: query successful. Rows affected: %d", operation, rowsAffected)
	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	const operation = "userRepository.GetUserByEmail"
	
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE email = $1`

	log.Printf("%s: executing query\n%s\nWith param: [%s]", operation, query, email)

	row := r.db.QueryRow(query, email)
	var user models.User

	log.Printf("%s: starting row scan", operation)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("%s: user not found with email: %s", operation, email)
			return nil, errors.New("пользователь не найден")
		}
		
		log.Printf("%s: ERROR in Scan - %v", operation, err)
		log.Printf("%s: TIP: Check data types compatibility between DB and struct fields", operation)
		return nil, err
	}

	log.Printf("%s: successfully retrieved user: %+v", operation, user)
	return &user, nil
}