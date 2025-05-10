package service

import (
	"errors"

	"family_finance_back/internal/models"
	"family_finance_back/internal/repository"
	"family_finance_back/internal/util"
)

// UserService определяет интерфейс для работы с данными пользователей
type UserService interface {
	// GetUserByToken получает данные пользователя по JWT токену
	// Проверяет валидность токена и возвращает данные пользователя
	GetUserByToken(token string) (*models.User, error)

	// GetUserByEmail получает данные пользователя по email
	// Возвращает nil, если пользователь не найден
	GetUserByEmail(email string) (*models.User, error)

	// UpdateUser обновляет данные пользователя
	// Обновляет только предоставленные поля
	UpdateUser(user *models.User) error
}

// userService реализует интерфейс UserService
type userService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewUserService создает новый экземпляр UserService
func NewUserService(userRepo repository.UserRepository, jwtSecret string) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *userService) GetUserByToken(token string) (*models.User, error) {
	// Декодируем JWT и получаем email
	claims, err := util.ValidateJWT(token, s.jwtSecret)
	if err != nil {
		return nil, errors.New("недействительный токен")
	}

	// Получаем пользователя по email
	user, err := s.userRepo.GetByEmail(claims.Email)
	if err != nil {
		return nil, errors.New("не удалось получить данные пользователя")
	}
	if user == nil {
		return nil, errors.New("пользователь не найден")
	}

	return user, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("не удалось получить данные пользователя")
	}
	if user == nil {
		return nil, errors.New("пользователь не найден")
	}
	return user, nil
}

func (s *userService) UpdateUser(user *models.User) error {
	if user == nil {
		return errors.New("данные пользователя не предоставлены")
	}
	if user.Email == "" {
		return errors.New("email пользователя обязателен")
	}
	return s.userRepo.Update(user)
}
