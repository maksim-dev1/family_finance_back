package service

import (
	"fmt"
	"family_finance_back/config"
	"family_finance_back/models"
	"family_finance_back/repository"
	"family_finance_back/utils"
	"time"
)

type AuthService interface {
	SendVerificationCode(email string) error
	VerifyCode(email, code string) (string, error) // возвращает JWT-токен
}

type authService struct {
	userRepo repository.UserRepository
	cfg      config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// SendVerificationCode генерирует код, отправляет его на email и сохраняет/обновляет пользователя в БД.
func (s *authService) SendVerificationCode(emailAddr string) error {
	// Валидация email
	if !utils.IsValidEmail(emailAddr) {
		return fmt.Errorf("Некорректный формат почты")
	}

	code := utils.GenerateVerificationCode()
	// Код действителен 1 минуту
	expiryTime := time.Now().Add(1 * time.Minute)

	subject := "Ваш проверочный код"
	body := fmt.Sprintf("Ваш проверочный код: %s. Код действителен в течение 1 минуты.", code)
	if err := utils.SendEmail(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUsername, s.cfg.SMTPPassword, emailAddr, subject, body); err != nil {
		return err
	}

	user, err := s.userRepo.GetUserByEmail(emailAddr)
	if err != nil {
		return err
	}
	if user == nil {
		newUser := &models.User{
			Email:            emailAddr,
			VerificationCode: code,
			CodeExpiresAt:    expiryTime,
			IsVerified:       false,
		}
		return s.userRepo.CreateUser(newUser)
	} else {
		return s.userRepo.UpdateUserVerification(emailAddr, code, expiryTime, false)
	}
}

// VerifyCode проверяет корректность введенного кода, сравнивает время истечения и,
// если код действителен, помечает пользователя как подтвержденного и генерирует JWT-токен.
func (s *authService) VerifyCode(emailAddr, code string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(emailAddr)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("пользователь не найден")
	}
	// Проверка, не истёк ли код
	if time.Now().After(user.CodeExpiresAt) {
		return "", fmt.Errorf("код истёк. Пожалуйста, запросите новый код")
	}
	if user.VerificationCode != code {
		return "", fmt.Errorf("неверный проверочный код")
	}
	// Обновляем статус пользователя: подтверждён
	if err := s.userRepo.UpdateUserVerification(emailAddr, code, user.CodeExpiresAt, true); err != nil {
		return "", err
	}
	// Генерируем JWT-токен
	token, err := utils.GenerateJWTToken(user.Email, s.cfg.JWTSecret, s.cfg.JWTExpiryMinutes)
	if err != nil {
		return "", fmt.Errorf("ошибка при генерации токена: %v", err)
	}
	return token, nil
}
