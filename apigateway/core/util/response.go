package util

import "github.com/gin-gonic/gin"

func SendSuccess(context *gin.Context, status int, data interface{}) {
	context.JSON(status, gin.H{
		"success": true,
		"data":    data,
	})
}

func SendError(context *gin.Context, status int, message string, err error) {
	context.JSON(status, gin.H{
		"success": false,
		"error": gin.H{
			"message": message,
			"details": err.Error(),
		},
	})
}
