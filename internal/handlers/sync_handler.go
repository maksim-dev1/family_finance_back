package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"myapp/internal/models"
	"myapp/internal/service"
)

type SyncHandler struct {
	syncService service.SyncService
}

func NewSyncHandler(syncService service.SyncService) *SyncHandler {
	return &SyncHandler{
		syncService: syncService,
	}
}

func (h *SyncHandler) SyncTransactions(c *gin.Context) {
	var transactions []*models.Transaction
	if err := c.ShouldBindJSON(&transactions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	err := h.syncService.SyncTransactions(userID.(string), transactions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Синхронизация выполнена успешно"})
}
