package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"family_finance_back/config"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	gomail "gopkg.in/gomail.v2"
)

// AuthService предоставляет методы для работы с авторизацией:
// генерация и отправка кода, верификация, генерация JWT токенов, refresh и logout.
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

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// InitiateRegistration начинает процесс регистрации: генерирует код, сохраняет его в Redis вместе с именем и отправляет письмо.
func (a *AuthService) InitiateRegistration(name, email string) error {
	code, err := a.GenerateCode()
	if err != nil {
		return err
	}
	// Сохраняем код и имя в Redis в формате "код:имя" с временем жизни 10 минут
	data := fmt.Sprintf("%s:%s", code, name)
	ctx := context.Background()
	err = a.redis.Set(ctx, "register:"+email, data, time.Minute*10).Err()
	if err != nil {
		return err
	}
	// Отправляем письмо с кодом
	err = a.SendCodeEmail(email, code)
	if err != nil {
		return err
	}
	return nil
}

// InitiateLogin начинает процесс входа: генерирует код, сохраняет его в Redis и отправляет письмо.
// Проверяется, что пользователь с таким email уже существует.
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
	err = a.redis.Set(ctx, "login:"+email, code, time.Minute*10).Err()
	if err != nil {
		return err
	}
	err = a.SendCodeEmail(email, code)
	if err != nil {
		return err
	}
	return nil
}

// VerifyCode проверяет введённый пользователем код (для регистрации или входа).
// Если код корректный, для регистрации создаётся пользователь (если его нет)
// и генерируются JWT токены (access и refresh).
func (a *AuthService) VerifyCode(email, code string) (string, error) {
	ctx := context.Background()
	regKey := "register:" + email
	loginKey := "login:" + email

	var redisKey string
	val, err := a.redis.Get(ctx, regKey).Result()
	if err == redis.Nil {
		// Если регистрационный ключ не найден, пробуем ключ для входа
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
			// Здесь используется функция PostgreSQL для генерации UUID (например, gen_random_uuid())
			_, err = a.db.Exec("INSERT INTO users (id, name, email, created_at, updated_at) VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())", name, email)
			if err != nil {
				return "", err
			}
		}
	}

	// Генерируем access и refresh токены
	accessToken, err := a.generateJWT(email, time.Minute*15)
	if err != nil {
		return "", err
	}
	refreshToken, err := a.generateJWT(email, time.Hour*24*7)
	if err != nil {
		return "", err
	}

	// Формируем JSON-ответ с токенами
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

// generateJWT генерирует JWT токен для указанного email с заданным временем жизни.
func (a *AuthService) generateJWT(email string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(duration).Unix(),
	})
	tokenString, err := token.SignedString([]byte(a.cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// RefreshToken проверяет refresh токен и генерирует новый access токен.
func (a *AuthService) RefreshToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(a.cfg.JWTSecret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return "", errors.New("неверные данные токена")
		}
		newAccessToken, err := a.generateJWT(email, time.Minute*15)
		if err != nil {
			return "", err
		}
		return newAccessToken, nil
	}
	return "", errors.New("неверный refresh токен")
}

// Logout "выходит" из системы, добавляя токен в blacklist в Redis до момента его истечения.
func (a *AuthService) Logout(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
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
	err = a.redis.Set(ctx, "blacklist:"+tokenString, "true", duration).Err()
	if err != nil {
		return err
	}
	return nil
}
