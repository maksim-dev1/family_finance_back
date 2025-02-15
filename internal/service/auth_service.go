package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"family_finance_back/config"
	"family_finance_back/internal/models"
	"family_finance_back/internal/repository"
)

type AuthService interface {
	Register(name, email string) (string, error)
	Login(email string) (string, error)
	VerifyCode(email, code string) (string, string, error)
	RefreshToken(refreshToken string) (string, string, error)
	Logout(userID string) error
}

type authService struct {
	userRepo repository.UserRepository
	redis    *redis.Client
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, redis *redis.Client, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		redis:    redis,
		cfg:      cfg,
	}
}

// Генерация 6-значного кода
func generateCode() (string, error) {
	max := big.NewInt(899999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

// Генерация access и refresh токенов
func (s *authService) generateTokens(user *models.User) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	accessTokenStr, err := accessToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken := uuid.New().String()
	ctx := context.Background()
	err = s.redis.Set(ctx, "refresh:"+user.ID, refreshToken, 7*24*time.Hour).Err()
	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshToken, nil
}

func (s *authService) Register(name, email string) (string, error) {
	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		return "", fmt.Errorf("пользователь с таким email уже существует")
	}
	code, err := generateCode()
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	if err := s.redis.Set(ctx, email, code, time.Minute).Err(); err != nil {
		return "", err
	}
	user := &models.User{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
	}
	if err := s.userRepo.CreateUser(user); err != nil {
		return "", err
	}
	return code, nil
}

func (s *authService) Login(email string) (string, error) {
	_, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("пользователь не существует")
	}
	code, err := generateCode()
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	if err := s.redis.Set(ctx, email, code, time.Minute).Err(); err != nil {
		return "", err
	}
	return code, nil
}

func (s *authService) VerifyCode(email, code string) (string, string, error) {
	ctx := context.Background()
	storedCode, err := s.redis.Get(ctx, email).Result()
	if err != nil {
		return "", "", fmt.Errorf("код истёк или неверный")
	}
	if storedCode != code {
		return "", "", fmt.Errorf("неверный код")
	}
	s.redis.Del(ctx, email)
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}
	return s.generateTokens(user)
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	ctx := context.Background()
	keys, err := s.redis.Keys(ctx, "refresh:*").Result()
	if err != nil {
		return "", "", fmt.Errorf("ошибка проверки refresh-токена")
	}

	var userID string
	for _, key := range keys {
		storedToken, _ := s.redis.Get(ctx, key).Result()
		if storedToken == refreshToken {
			userID = key[len("refresh:"):]
			break
		}
	}

	if userID == "" {
		return "", "", fmt.Errorf("недействительный refresh-токен")
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return "", "", fmt.Errorf("пользователь не найден")
	}

	return s.generateTokens(user)
}

func (s *authService) Logout(userID string) error {
	ctx := context.Background()
	return s.redis.Del(ctx, "refresh:"+userID).Err()
}
