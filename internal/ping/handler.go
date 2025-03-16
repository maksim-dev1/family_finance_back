package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler отвечает на запросы для проверки работоспособности сервера.
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
