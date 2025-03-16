package auth

import (
    "net/http"
    "context"

    "github.com/gin-gonic/gin"
)

// AuthHandler содержит обработчики для аутентификации.
type AuthHandler struct {
    service Service
}

// NewAuthHandler создаёт новый экземпляр AuthHandler.
func NewAuthHandler(service Service) *AuthHandler {
    return &AuthHandler{service: service}
}

// SendCodeRequest описывает входной JSON для отправки кода.
type SendCodeRequest struct {
    Email string `json:"email" binding:"required,email"`
}

// VerifyCodeRequest описывает входной JSON для верификации кода.
// Поле Name используется при регистрации нового пользователя.
type VerifyCodeRequest struct {
    Email string `json:"email" binding:"required,email"`
    Code  string `json:"code" binding:"required,len=6,numeric"`
    Name  string `json:"name" binding:"omitempty,min=1"`
}

// SendCode обрабатывает запрос на отправку кода.
func (h *AuthHandler) SendCode(c *gin.Context) {
    var req SendCodeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := h.service.SendVerificationCode(context.Background(), req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "verification code sent"})
}

// VerifyCode обрабатывает запрос на проверку кода и возвращает JWT-токен.
func (h *AuthHandler) VerifyCode(c *gin.Context) {
    var req VerifyCodeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    token, err := h.service.VerifyCode(context.Background(), req.Email, req.Code, req.Name)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"token": token})
}