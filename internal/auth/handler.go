package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthHandler содержит HTTP‑обработчики для авторизации.
type AuthHandler struct {
	authService *AuthService
	// Добавим доступ к базе для проверки регистрации
}

// NewAuthHandler создаёт новый экземпляр AuthHandler.
func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRequest описывает структуру запроса для регистрации.
type RegisterRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// LoginRequest описывает структуру запроса для входа.
type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyRequest описывает структуру запроса для проверки кода.
type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// RefreshRequest описывает структуру запроса для обновления токена.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest описывает структуру запроса для выхода.
type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

// checkAlreadyAuthorized проверяет, передан ли валидный access токен в заголовке.
func checkAlreadyAuthorized(c *gin.Context, secret string) bool {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return false
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return false
	}
	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	return err == nil && token.Valid
}

// Register обрабатывает запрос на регистрацию.
func (h *AuthHandler) Register(c *gin.Context) {
	// Если уже авторизован, возвращаем сообщение
	if checkAlreadyAuthorized(c, h.authService.cfg.JWTSecret) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "пользователь уже авторизован"})
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Внутри InitiateRegistration проверяется наличие пользователя с данным email.
	if err := h.authService.InitiateRegistration(req.Name, req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "код подтверждения отправлен на email"})
}

// Login обрабатывает запрос на вход.
func (h *AuthHandler) Login(c *gin.Context) {
	// Если уже авторизован, возвращаем сообщение
	if checkAlreadyAuthorized(c, h.authService.cfg.JWTSecret) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "пользователь уже авторизован"})
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.authService.InitiateLogin(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "код подтверждения отправлен на email"})
}

// Verify обрабатывает проверку кода и возвращает JWT токены.
func (h *AuthHandler) Verify(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokensJSON, err := h.authService.VerifyCode(req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var tokenData map[string]string
	if err := json.Unmarshal([]byte(tokensJSON), &tokenData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось обработать токены"})
		return
	}
	c.JSON(http.StatusOK, tokenData)
}

// Refresh обрабатывает обновление access токена по refresh токену.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newAccessToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}

// Logout обрабатывает выход пользователя.
func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.authService.Logout(req.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "вы успешно вышли"})
}
