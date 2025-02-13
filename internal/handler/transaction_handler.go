package handler

import (
	"net/http"
	"strings"
	"time"

	"family_finance_back/internal/service"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

type createTransactionRequest struct {
	Type          string  `json:"type" binding:"required"`
	Category      string  `json:"category" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Date          string  `json:"date" binding:"required"`
	SavingsGoalID *string `json:"savings_goal_id"`
	Description   *string `json:"description"`
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	t, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты"})
		return
	}
	tx, err := h.transactionService.CreateTransaction(userID.(string), req.Type, req.Category, req.Amount, t, req.SavingsGoalID, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tx)
}

func (h *TransactionHandler) GetPersonalTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	txs, err := h.transactionService.GetPersonalTransactions(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, txs)
}

func (h *TransactionHandler) GetGroupTransactions(c *gin.Context) {
	usersParam := c.Query("user_ids")
	if usersParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указаны идентификаторы пользователей"})
		return
	}
	userIDs := splitAndTrim(usersParam)
	txs, err := h.transactionService.GetGroupTransactions(userIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, txs)
}

func splitAndTrim(s string) []string {
	var res []string
	parts := strings.Split(s, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}
