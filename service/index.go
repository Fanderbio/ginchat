package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "welcome!",
	})
}
