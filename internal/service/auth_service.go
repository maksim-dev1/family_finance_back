package service

import (
	"context"
	"errors"
	"time"

	"family_finance_back/internal/models"
	"family_finance_back/internal/repository"
	"family_finance_back/internal/util"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type AuthService interface {
	RequestLoginCode(email string) error
	VerifyLoginCode(email, code string) (string, error)
	RequestRegistrationCode(email string) (string, error)
	VerifyRegistrationCode(email, code, name, surname, nickname string) error
	Logout(token string) error
}

type authService struct {
	userRepo    repository.UserRepository
	emailSvc    EmailService
	redisClient *redis.Client
	jwtSecret   string
	ctx         context.Context
}

func NewAuthService(userRepo repository.UserRepository, emailSvc EmailService, redisClient *redis.Client, jwtSecret string) AuthService {
	return &authService{
		userRepo:    userRepo,
		emailSvc:    emailSvc,
		redisClient: redisClient,
		jwtSecret:   jwtSecret,
		ctx:         context.Background(),
	}
}

// RequestLoginCode отправляет код, если почта существует
func (s *authService) RequestLoginCode(email string) error {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("почта не найдена")
	}

	code := util.GenerateCode()
	// Сохраняем код в Redis с TTL 90 секунд
	err = s.redisClient.Set(s.ctx, "login:"+email, code, 90*time.Second).Err()
	if err != nil {
		return err
	}

	// Отправляем код на почту
	return s.emailSvc.SendCode(email, code)
}

// VerifyLoginCode проверяет код и генерирует JWT
func (s *authService) VerifyLoginCode(email, code string) (string, error) {
	storedCode, err := s.redisClient.Get(s.ctx, "login:"+email).Result()
	if err != nil {
		return "", errors.New("код не найден или истёк")
	}
	if storedCode != code {
		return "", errors.New("неверный код")
	}
	// Генерируем JWT (с очень большим сроком жизни)
	token, err := util.GenerateJWT(email, s.jwtSecret)
	if err != nil {
		return "", err
	}
	// Удаляем код из Redis
	s.redisClient.Del(s.ctx, "login:"+email)
	return token, nil
}

// RequestRegistrationCode отправляет код для регистрации и сохраняет связь uuid -> email
func (s *authService) RequestRegistrationCode(email string) (string, error) {
	// Проверяем, что пользователя с таким email ещё нет
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if user != nil {
		return "", errors.New("пользователь с такой почтой уже существует")
	}

	// Генерируем код и сохраняем его в Redis с TTL 90 секунд
	code := util.GenerateCode()
	err = s.redisClient.Set(s.ctx, "register:"+email, code, 90*time.Second).Err()
	if err != nil {
		return "", err
	}

	// Генерируем UUID для регистрации и сохраняем его вместе с email на 15 минут
	uuidKey := uuid.New().String()
	err = s.redisClient.Set(s.ctx, "registration_uuid:"+uuidKey, email, 15*time.Minute).Err()
	if err != nil {
		return "", err
	}

	// Отправляем код на почту
	if err := s.emailSvc.SendCode(email, code); err != nil {
		return "", err
	}

	// Возвращаем UUID клиенту
	return uuidKey, nil
}

// VerifyRegistrationCode проверяет код и регистрирует нового пользователя, используя UUID
func (s *authService) VerifyRegistrationCode(registrationUUID, code, name, surname, nickname string) error {
	// Получаем email по UUID
	email, err := s.redisClient.Get(s.ctx, "registration_uuid:"+registrationUUID).Result()
	if err != nil {
		return errors.New("регистрация не найдена или истёкла")
	}

	// Получаем сохранённый код по email и сверяем
	storedCode, err := s.redisClient.Get(s.ctx, "register:"+email).Result()
	if err != nil {
		return errors.New("код не найден или истёк")
	}
	if storedCode != code {
		return errors.New("неверный код")
	}

	// Если nickname пустой – используем name
	if nickname == "" {
		nickname = name
	}
	newUser := &models.User{
		Name:     name,
		Surname:  surname,
		Nickname: nickname,
		Email:    email,
	}
	err = s.userRepo.Create(newUser)
	if err != nil {
		return err
	}
	// Удаляем данные регистрации из Redis
	s.redisClient.Del(s.ctx, "register:"+email)
	s.redisClient.Del(s.ctx, "registration_uuid:"+registrationUUID)
	return nil
}

// Logout добавляет токен в blacklist в Redis, чтобы его нельзя было использовать далее
func (s *authService) Logout(token string) error {
	// Сохраняем токен в blacklist с большим TTL (например, 100 лет)
	return s.redisClient.Set(s.ctx, "blacklist:"+token, "true", 100*365*24*time.Hour).Err()
}
