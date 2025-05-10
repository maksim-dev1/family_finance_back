package repository

import (
	"errors"

	"family_finance_back/internal/models"

	"gorm.io/gorm"
)

// UserRepository определяет интерфейс для работы с данными пользователей в базе данных
type UserRepository interface {
	// GetByEmail получает пользователя по email
	// Возвращает nil, если пользователь не найден
	GetByEmail(email string) (*models.User, error)

	// Create создает нового пользователя
	// Возвращает ошибку, если пользователь с таким email уже существует
	Create(user *models.User) error

	// Update обновляет данные пользователя
	// Обновляет все поля модели
	Update(user *models.User) error
}

// userRepository реализует интерфейс UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}
