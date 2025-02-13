package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"family_finance_back/internal/service"
)

type AuthHandler struct {
	authService  service.AuthService
	emailService service.EmailService
}

func NewAuthHandler(authService service.AuthService, emailService service.EmailService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		emailService: emailService,
	}
}

type registerRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code, err := h.authService.Register(req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Отправляем код на email пользователя
	if err := h.emailService.SendCode(req.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось отправить email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Регистрация успешна. Код отправлен на email."})
}

type loginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code, err := h.authService.Login(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Отправляем код на email пользователя
	if err := h.emailService.SendCode(req.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось отправить email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Код для входа отправлен на email."})
}

type verifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

func (h *AuthHandler) VerifyCode(c *gin.Context) {
	var req verifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.authService.VerifyCode(req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
