package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"family_finance_back/internal/models"
	"family_finance_back/internal/repository"
	"family_finance_back/internal/util"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type AuthService interface {
	RequestLoginCode(email string) (string, error)
	VerifyLoginCode(tempID, code string) (string, error)
	RequestRegistrationCode(email string) (string, error)
	VerifyRegistrationCode(tempID, code, name, surname, nickname string) (string, error)
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
func (s *authService) RequestLoginCode(email string) (string, error) {
	// Проверяем, существует ли пользователь
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("не удалось получить данные пользователя, попробуйте позже")
	}
	if user == nil {
		return "", errors.New("пользователь с указанным email не найден")
	}

	code := util.GenerateCode()
	tempID := uuid.New().String()

	data := map[string]string{
		"email": email,
		"code":  code,
	}
	serialized, err := json.Marshal(data)
	if err != nil {
		return "", errors.New("не удалось сформировать данные для отправки кода")
	}

	err = s.redisClient.Set(s.ctx, "login:"+tempID, serialized, 90*time.Second).Err()
	if err != nil {
		return "", errors.New("не удалось сохранить данные авторизации, повторите попытку позже")
	}

	if err = s.emailSvc.SendCode(email, code); err != nil {
		return "", errors.New("не удалось отправить письмо с кодом, пожалуйста, проверьте адрес электронной почты")
	}

	return tempID, nil
}

// VerifyLoginCode проверяет код и генерирует JWT
func (s *authService) VerifyLoginCode(tempID, code string) (string, error) {
	val, err := s.redisClient.Get(s.ctx, "login:"+tempID).Result()
	if err != nil {
		return "", errors.New("код не найден или срок действия кода истёк, повторите запрос")
	}
	var data map[string]string
	if err = json.Unmarshal([]byte(val), &data); err != nil {
		return "", errors.New("не удалось обработать данные авторизации, повторите попытку")
	}
	if data["code"] != code {
		return "", errors.New("введён неверный код, пожалуйста, проверьте и повторите попытку")
	}

	token, err := util.GenerateJWT(data["email"], s.jwtSecret)
	if err != nil {
		return "", errors.New("не удалось сгенерировать токен авторизации, повторите попытку позже")
	}
	s.redisClient.Del(s.ctx, "login:"+tempID)
	return token, nil
}

// RequestRegistrationCode отправляет код для регистрации и сохраняет связь uuid -> email
func (s *authService) RequestRegistrationCode(email string) (string, error) {
	// Проверка существования пользователя
	existingUser, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("не удалось проверить данные. Попробуйте позже")
	}
	if existingUser != nil {
		return "", errors.New("пользователь с указанным email уже существует")
	}

	code := util.GenerateCode()
	tempID := uuid.New().String()

	data := map[string]string{
		"email": email,
		"code":  code,
	}
	serialized, err := json.Marshal(data)
	if err != nil {
		return "", errors.New("не удалось сформировать данные для отправки кода")
	}

	// Сохраняем данные в Redis с TTL 90 секунд
	err = s.redisClient.Set(s.ctx, "register:"+tempID, serialized, 90*time.Second).Err()
	if err != nil {
		return "", errors.New("не удалось сохранить данные регистрации, повторите попытку позже")
	}

	// Отправляем код на указанный email
	if err = s.emailSvc.SendCode(email, code); err != nil {
		return "", errors.New("не удалось отправить письмо с кодом, пожалуйста, проверьте адрес электронной почты")
	}

	return tempID, nil
}

// VerifyRegistrationCode проверяет код и регистрирует нового пользователя, используя UUID
func (s *authService) VerifyRegistrationCode(tempID, code, name, surname, nickname string) (string, error) {
	val, err := s.redisClient.Get(s.ctx, "register:"+tempID).Result()
	if err != nil {
		return "", errors.New("код не найден или срок действия кода истёк, повторите запрос")
	}
	var data map[string]string
	if err = json.Unmarshal([]byte(val), &data); err != nil {
		return "", errors.New("не удалось обработать данные регистрации, повторите попытку")
	}
	if data["code"] != code {
		return "", errors.New("введён неверный код, пожалуйста, проверьте и повторите попытку")
	}

	if nickname == "" {
		nickname = name
	}
	newUser := &models.User{
		Name:     name,
		Surname:  surname,
		Nickname: nickname,
		Email:    data["email"],
	}
	if err = s.userRepo.Create(newUser); err != nil {
		return "", errors.New("не удалось создать пользователя, попробуйте позже")
	}
	token, err := util.GenerateJWT(data["email"], s.jwtSecret)
	if err != nil {
		return "", errors.New("не удалось сгенерировать токен авторизации, повторите попытку позже")
	}
	s.redisClient.Del(s.ctx, "register:"+tempID)
	return token, nil
}


// Logout добавляет токен в blacklist в Redis, чтобы его нельзя было использовать далее
func (s *authService) Logout(token string) error {
	// Сохраняем токен в blacklist с большим TTL (например, 100 лет)
	return s.redisClient.Set(s.ctx, "blacklist:"+token, "true", 100*365*24*time.Hour).Err()
}
