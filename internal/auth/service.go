package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"family_finance_back/internal/models"
	redisstore "family_finance_back/internal/redisstore"
	"family_finance_back/internal/smtp"

	redis "github.com/go-redis/redis/v8"

	"github.com/golang-jwt/jwt/v4"
)

// Service описывает бизнес-логику аутентификации.
type Service interface {
	SendVerificationCode(ctx context.Context, email string) error
	VerifyCode(ctx context.Context, email, code, name string) (string, error)
}

// AuthService реализует Service.
type AuthService struct {
	repo        Repository
	redisClient *redis.Client
	mailer      *smtp.Mailer
	jwtSecret   string
	codeTTL     time.Duration
}

// NewAuthService создаёт новый экземпляр AuthService.
func NewAuthService(repo Repository, redisClient *redis.Client, mailer *smtp.Mailer, jwtSecret string) *AuthService {
	return &AuthService{
		repo:        repo,
		redisClient: redisClient,
		mailer:      mailer,
		jwtSecret:   jwtSecret,
		codeTTL:     10 * time.Minute, // Код действителен 10 минут
	}
}

// generateCode генерирует случайный 6-значный код.
func generateCode() (string, error) {
	max := big.NewInt(899999) // диапазон: 0..899999, затем +100000
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	code := 100000 + n.Int64()
	return fmt.Sprintf("%06d", code), nil
}

// SendVerificationCode отправляет на email 6-значный код подтверждения.
func (s *AuthService) SendVerificationCode(ctx context.Context, email string) error {
	code, err := generateCode()
	if err != nil {
		return err
	}
	// Сохраняем код в Redis с временем жизни codeTTL
	err = redisstore.SetCode(ctx, s.redisClient, email, code, s.codeTTL)
	if err != nil {
		return err
	}
	// Формируем сообщение письма
	subject := "Ваш код подтверждения для Family Finance"
	body := fmt.Sprintf("<p>Ваш код подтверждения: <b>%s</b></p><p>Код действителен в течение 10 минут.</p>", code)
	// Отправляем email
	err = s.mailer.SendEmail(email, subject, body)
	if err != nil {
		return err
	}
	return nil
}

// VerifyCode проверяет код. Если код корректный, то
// – создаёт пользователя (если отсутствует) и
// – возвращает JWT-токен.
func (s *AuthService) VerifyCode(ctx context.Context, email, code, name string) (string, error) {
	storedCode, err := redisstore.GetCode(ctx, s.redisClient, email)
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("verification code expired or not found")
		}
		return "", err
	}

	if storedCode != code {
		return "", errors.New("invalid verification code")
	}

	// Удаляем код из Redis
	err = redisstore.DeleteCode(ctx, s.redisClient, email)
	if err != nil {
		return "", err
	}

	// Проверяем наличие пользователя
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		// Если пользователь отсутствует – создаём нового.
		user = &models.User{
			Name:  name,
			Email: email,
		}
		err = s.repo.CreateUser(user)
		if err != nil {
			return "", err
		}
	}

	// Генерируем JWT-токен
	tokenString, err := s.generateJWTToken(user)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// generateJWTToken генерирует JWT-токен для пользователя.
func (s *AuthService) generateJWTToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // токен действителен 24 часа
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
