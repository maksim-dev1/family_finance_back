package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"family_finance_back/config"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	gomail "gopkg.in/gomail.v2"
)

// AuthService предоставляет методы для работы с авторизацией.
type AuthService struct {
	db    *sql.DB
	redis *redis.Client
	cfg   *config.Config
}

// NewAuthService создаёт новый экземпляр AuthService.
func NewAuthService(db *sql.DB, redis *redis.Client, cfg *config.Config) *AuthService {
	return &AuthService{
		db:    db,
		redis: redis,
		cfg:   cfg,
	}
}

// GenerateCode генерирует безопасный 6-значный числовой код для верификации.
func (a *AuthService) GenerateCode() (string, error) {
	max := big.NewInt(900000) // диапазон: 0 ... 899999
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	code := int(n.Int64()) + 100000 // смещаем диапазон до 100000 ... 999999
	return fmt.Sprintf("%06d", code), nil
}

// SendCodeEmail отправляет письмо с верификационным кодом на указанный email.
func (a *AuthService) SendCodeEmail(recipient, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", a.cfg.SMTPUsername)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "Ваш верификационный код")
	m.SetBody("text/plain", fmt.Sprintf("Ваш верификационный код: %s", code))

	port, err := strconv.Atoi(a.cfg.SMTPPort)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(a.cfg.SMTPHost, port, a.cfg.SMTPUsername, a.cfg.SMTPPassword)
	d.SSL = true // для порта 465 требуется SSL

	return d.DialAndSend(m)
}

// InitiateRegistration начинает процесс регистрации: генерирует код, сохраняет его в Redis (TTL 120 сек) и отправляет письмо.
func (a *AuthService) InitiateRegistration(name, email string) error {
	// Проверяем, существует ли пользователь
	var exists bool
	err := a.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("пользователь с данным email уже зарегистрирован")
	}

	code, err := a.GenerateCode()
	if err != nil {
		return err
	}
	// Сохраняем код и имя в Redis в формате "код:имя" с временем жизни 120 секунд
	data := fmt.Sprintf("%s:%s", code, name)
	ctx := context.Background()
	if err := a.redis.Set(ctx, "register:"+email, data, 120*time.Second).Err(); err != nil {
		return err
	}
	return a.SendCodeEmail(email, code)
}

// InitiateLogin начинает процесс входа: генерирует код, сохраняет его в Redis (TTL 120 сек) и отправляет письмо.
func (a *AuthService) InitiateLogin(email string) error {
	// Проверяем, существует ли пользователь
	var exists bool
	err := a.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("пользователь не найден, пройдите регистрацию")
	}
	code, err := a.GenerateCode()
	if err != nil {
		return err
	}
	ctx := context.Background()
	if err := a.redis.Set(ctx, "login:"+email, code, 120*time.Second).Err(); err != nil {
		return err
	}
	return a.SendCodeEmail(email, code)
}

// VerifyCode проверяет код подтверждения и генерирует JWT токены.
func (a *AuthService) VerifyCode(email, code string) (string, error) {
	ctx := context.Background()
	regKey := "register:" + email
	loginKey := "login:" + email

	var redisKey string
	val, err := a.redis.Get(ctx, regKey).Result()
	if err == redis.Nil {
		val, err = a.redis.Get(ctx, loginKey).Result()
		if err == redis.Nil {
			return "", errors.New("код подтверждения истёк или неверный")
		} else if err != nil {
			return "", err
		}
		redisKey = loginKey
	} else if err != nil {
		return "", err
	} else {
		redisKey = regKey
	}

	var storedCode, name string
	if redisKey == regKey {
		parts := strings.SplitN(val, ":", 2)
		if len(parts) != 2 {
			return "", errors.New("неверный формат данных в redis")
		}
		storedCode = parts[0]
		name = parts[1]
	} else {
		storedCode = val
	}

	if storedCode != code {
		return "", errors.New("введён неверный код подтверждения")
	}

	// Удаляем ключ из Redis после успешной проверки
	a.redis.Del(ctx, redisKey)

	// Если это регистрация, создаём пользователя, если его ещё нет
	if redisKey == regKey {
		var exists bool
		err := a.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
		if err != nil {
			return "", err
		}
		if !exists {
			_, err = a.db.Exec("INSERT INTO users (id, name, email, created_at, updated_at) VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())", name, email)
			if err != nil {
				return "", err
			}
		}
	}

	// Генерируем access и refresh токены с новыми сроками
	accessToken, err := a.generateJWT(email, time.Hour, "access")
	if err != nil {
		return "", err
	}
	refreshToken, err := a.generateJWT(email, time.Hour*24*7, "refresh")
	if err != nil {
		return "", err
	}

	tokenData := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	jsonBytes, err := json.Marshal(tokenData)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// generateJWT генерирует JWT токен с заданным временем жизни для access-токена
// или для refresh-токена, который теперь живёт 30 дней.
func (a *AuthService) generateJWT(email string, duration time.Duration, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"email":      email,
		"token_type": tokenType,
	}
	var secret []byte
	if tokenType == "refresh" {
		// Для refresh-токена задаём срок действия 30 дней.
		claims["exp"] = time.Now().Add(30 * 24 * time.Hour).Unix()
		secret = []byte(a.cfg.JWTSecret + "_refresh")
	} else {
		// Для access-токена задаём срок действия, переданный в параметре (например, 1 час).
		claims["exp"] = time.Now().Add(duration).Unix()
		secret = []byte(a.cfg.JWTSecret)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// RefreshToken проверяет refresh-токен и генерирует новый access-токен и новый refresh-токен.
func (a *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма подписи.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(a.cfg.JWTSecret + "_refresh"), nil
	})
	if err != nil {
		return "", "", errors.New("refresh токен недействителен, требуется повторная авторизация")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return "", "", errors.New("неверные данные токена")
		}
		// Генерируем новый access-токен (1 час).
		newAccessToken, err := a.generateJWT(email, time.Hour, "access")
		if err != nil {
			return "", "", err
		}
		// Генерируем новый refresh-токен (30 дней).
		newRefreshToken, err := a.generateJWT(email, 0, "refresh")
		if err != nil {
			return "", "", err
		}
		return newAccessToken, newRefreshToken, nil
	}
	return "", "", errors.New("неверный refresh токен")
}


// Logout "выходит" из системы, добавляя токен в blacklist в Redis.
func (a *AuthService) Logout(tokenString string) error {
	// Определяем тип токена (access) для blacklist'а
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(a.cfg.JWTSecret), nil
	})
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("неверный токен")
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("неверное время истечения токена")
	}
	duration := time.Until(time.Unix(int64(exp), 0))
	ctx := context.Background()
	return a.redis.Set(ctx, "blacklist:"+tokenString, "true", duration).Err()
}
