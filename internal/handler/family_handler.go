package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"family_finance_back/internal/service"
)

type FamilyHandler struct {
	familyService service.FamilyService
}

func NewFamilyHandler(familyService service.FamilyService) *FamilyHandler {
	return &FamilyHandler{
		familyService: familyService,
	}
}

type createFamilyRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *FamilyHandler) CreateFamily(c *gin.Context) {
	var req createFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	family, err := h.familyService.CreateFamily(req.Name, req.Description, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, family)
}

func (h *FamilyHandler) GetFamilies(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	families, err := h.familyService.GetUserFamilies(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, families)
}
