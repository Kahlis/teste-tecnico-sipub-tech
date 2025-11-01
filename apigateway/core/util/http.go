package util

import "github.com/gin-gonic/gin"

func SendJson(context *gin.Context, status int, payload interface{}) {
	context.JSON(status, payload)
}
