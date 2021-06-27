package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
	heartbeat connection
 */
func IsAlive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}