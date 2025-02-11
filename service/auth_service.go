package service

import (
	"errors"
	"fmt"
	"log"

	"family_finance_back/config"
	"family_finance_back/models"
	"family_finance_back/repository"
	"family_finance_back/utils"
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

// SendVerificationCode генерирует код, отправляет его на email и сохраняет/обновляет пользователя в БД.
func (s *authService) SendVerificationCode(emailAddr string) error {
	code, err := utils.GenerateVerificationCode()
	if err != nil {
		log.Printf("[ERROR] Ошибка генерации кода: %v", err)
		return fmt.Errorf("ошибка генерации кода")
	}

	subject := "Ваш проверочный код"
	body := fmt.Sprintf("Ваш проверочный код: %s", code)
	if err := utils.SendEmail(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUsername, s.cfg.SMTPPassword, emailAddr, subject, body); err != nil {
		log.Printf("[ERROR] Ошибка при отправке email на %s: %v", emailAddr, err)
		return fmt.Errorf("ошибка отправки email")
	}

	user, err := s.userRepo.GetUserByEmail(emailAddr)
	if err != nil {
		log.Printf("[ERROR] Ошибка поиска пользователя %s: %v", emailAddr, err)
		return fmt.Errorf("ошибка при поиске пользователя")
	}

	if user == nil {
		newUser := &models.User{
			Email:            emailAddr,
			VerificationCode: code,
			IsVerified:       false,
		}
		if err := s.userRepo.CreateUser(newUser); err != nil {
			log.Printf("[ERROR] Ошибка создания пользователя %s: %v", emailAddr, err)
			return fmt.Errorf("ошибка создания пользователя")
		}
	} else {
		if err := s.userRepo.UpdateUserVerification(emailAddr, code, false); err != nil {
			log.Printf("[ERROR] Ошибка обновления кода для пользователя %s: %v", emailAddr, err)
			return fmt.Errorf("ошибка обновления кода пользователя")
		}
	}

	log.Printf("[INFO] Код отправлен на %s", emailAddr)
	return nil
}

// VerifyCode проверяет корректность введенного кода и, если верно, помечает пользователя как подтвержденного.
func (s *authService) VerifyCode(emailAddr, code string) error {
	user, err := s.userRepo.GetUserByEmail(emailAddr)
	if err != nil {
		log.Printf("[ERROR] Ошибка поиска пользователя %s: %v", emailAddr, err)
		return fmt.Errorf("ошибка при поиске пользователя")
	}

	if user == nil {
		log.Printf("[ERROR] Пользователь %s не найден", emailAddr)
		return errors.New("пользователь не найден")
	}

	if user.VerificationCode != code {
		log.Printf("[ERROR] Неверный проверочный код для %s", emailAddr)
		return errors.New("неверный проверочный код")
	}

	if err := s.userRepo.UpdateUserVerification(emailAddr, code, true); err != nil {
		log.Printf("[ERROR] Ошибка подтверждения пользователя %s: %v", emailAddr, err)
		return fmt.Errorf("ошибка подтверждения пользователя")
	}

	log.Printf("[INFO] Пользователь %s успешно подтвержден", emailAddr)
	return nil
}
