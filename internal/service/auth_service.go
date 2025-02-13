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

	"myapp/config"
	"myapp/internal/models"
	"myapp/internal/repository"
)

type AuthService interface {
	Register(name, email string) (string, error) // возвращает сгенерированный код
	Login(email string) (string, error)            // возвращает сгенерированный код
	VerifyCode(email, code string) (string, error)   // возвращает JWT токен
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

func generateCode() (string, error) {
	max := big.NewInt(899999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	code := fmt.Sprintf("%06d", n.Int64()+100000)
	return code, nil
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

func (s *authService) VerifyCode(email, code string) (string, error) {
	ctx := context.Background()
	storedCode, err := s.redis.Get(ctx, email).Result()
	if err != nil {
		return "", fmt.Errorf("код истёк или неверный")
	}
	if storedCode != code {
		return "", fmt.Errorf("неверный код")
	}
	s.redis.Del(ctx, email)
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
