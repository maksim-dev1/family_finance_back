package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler содержит HTTP‑обработчики для работы с пользователями.
type UserHandler struct {
	userService *UserService
}

// NewUserHandler создаёт новый экземпляр UserHandler.
func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetAllUsers возвращает список всех пользователей.
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetCurrentUser возвращает данные аутентифицированного пользователя.
// Ожидается, что middleware установил в контекст email пользователя.
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	email, exists := c.Get("user_email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}
	user, err := h.userService.GetUserByEmail(email.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser удаляет текущего аутентифицированного пользователя по email, извлечённому из токена.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Извлекаем email пользователя, установленный JWT middleware
	email, exists := c.Get("user_email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}
	err := h.userService.DeleteUserByEmail(email.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "пользователь успешно удалён"})
}
