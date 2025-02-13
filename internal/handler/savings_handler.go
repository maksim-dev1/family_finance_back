package handler

import (
	"net/http"
	"time"

	"family_finance_back/internal/service"

	"github.com/gin-gonic/gin"
)

type SavingsHandler struct {
	savingsService service.SavingsService
}

func NewSavingsHandler(savingsService service.SavingsService) *SavingsHandler {
	return &SavingsHandler{
		savingsService: savingsService,
	}
}

type createSavingsRequest struct {
	TargetAmount float64 `json:"target_amount" binding:"required"`
	TargetDate   string  `json:"target_date" binding:"required"`
	StartDate    string  `json:"start_date" binding:"required"`
	Description  *string `json:"description"`
	FamilyID     *string `json:"family_id"`
}

func (h *SavingsHandler) CreateSavingsGoal(c *gin.Context) {
	var req createSavingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	goal, err := h.savingsService.CreateSavingsGoal(userID.(string), req.FamilyID, req.TargetAmount, req.TargetDate, req.StartDate, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, goal)
}

func (h *SavingsHandler) CalculateSavingPlan(c *gin.Context) {
	var req struct {
		TargetAmount float64 `json:"target_amount" binding:"required"`
		StartDate    string  `json:"start_date" binding:"required"`
		TargetDate   string  `json:"target_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты начала. Используйте YYYY-MM-DD."})
		return
	}
	target, err := time.Parse("2006-01-02", req.TargetDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат целевой даты. Используйте YYYY-MM-DD."})
		return
	}
	daily, weekly, monthly, err := h.savingsService.CalculateSavingPlan(req.TargetAmount, start, target)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"daily":   daily,
		"weekly":  weekly,
		"monthly": monthly,
	})
}

func (h *SavingsHandler) GetSavingsGoals(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	familyID := c.Query("family_id")
	if familyID != "" {
		goals, err := h.savingsService.GetFamilySavingsGoals(familyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, goals)
	} else {
		goals, err := h.savingsService.GetUserSavingsGoals(userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, goals)
	}
}
