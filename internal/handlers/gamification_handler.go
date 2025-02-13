package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"myapp/internal/service"
)

type GamificationHandler struct {
	gamificationService service.GamificationService
}

func NewGamificationHandler(gamificationService service.GamificationService) *GamificationHandler {
	return &GamificationHandler{
		gamificationService: gamificationService,
	}
}

func (h *GamificationHandler) GetUserScore(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	score, err := h.gamificationService.GetUserScore(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"score": score})
}
