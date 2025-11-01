package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoute(rg *gin.RouterGroup) {
	rg.GET("/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
}
