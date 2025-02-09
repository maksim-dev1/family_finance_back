// service/auth_service.go
package service

import (
	"family_finance_back/config"
	"family_finance_back/models"
	"family_finance_back/repository"
    "family_finance_back/utils/code"
    "family_finance_back/utils/email"
	"fmt"
)

type AuthService interface {
	SendVerificationCode(email string) error
	VerifyCode(email, code string) error
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

// SendVerificationCode генерирует код, отправляет его на email и сохраняет/обновляет пользователя в БД
func (s *authService) SendVerificationCode(emailAddr string) error {
	code := code.GenerateVerificationCode()

	// Отправка письма с кодом
	subject := "Ваш проверочный код"
	body := fmt.Sprintf("Ваш проверочный код: %s", code)
	if err := email.SendEmail(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUsername, s.cfg.SMTPPassword, emailAddr, subject, body); err != nil {
		return err
	}

	// Проверка, существует ли уже пользователь
	user, err := s.userRepo.GetUserByEmail(emailAddr)
	if err != nil {
		return err
	}
	if user == nil {
		// Создаем нового пользователя
		newUser := &models.User{
			Email:            emailAddr,
			VerificationCode: code,
			IsVerified:       false,
		}
		return s.userRepo.CreateUser(newUser)
	} else {
		// Обновляем код в существующей записи и сбрасываем статус подтверждения
		return s.userRepo.UpdateUserVerification(emailAddr, code, false)
	}
}

// VerifyCode проверяет корректность введенного кода и, если все верно, помечает пользователя как подтвержденного
func (s *authService) VerifyCode(emailAddr, code string) error {
	user, err := s.userRepo.GetUserByEmail(emailAddr)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("пользователь не найден")
	}

	if user.VerificationCode != code {
		return fmt.Errorf("неверный проверочный код")
	}

	// Обновляем статус подтверждения
	return s.userRepo.UpdateUserVerification(emailAddr, code, true)
}
